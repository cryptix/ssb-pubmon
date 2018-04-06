package models

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"cryptoscope.co/go/muxrpc"
	"github.com/agl/ed25519"
	"github.com/cryptix/go/debug"
	"github.com/cryptix/go/logging"
	humanize "github.com/dustin/go-humanize"
	kitlog "github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/qor/notification"
	"github.com/qor/qor"
	"github.com/qor/transition"

	"github.com/cryptix/ssb-pubmon/config/notify"
	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/ssb"
)

type Pub struct {
	gorm.Model
	Key         string `gorm:"unique_index"`
	Addresses   []Address
	LastSuccess time.Time
	transition.Transition
}

/* Composite primary keys can't be created on SQLite tables
    https://github.com/jinzhu/gorm/issues/1037

   	Host  string `gorm:"primary_key"`
   	Port  int    `gorm:"primary_key"`
*/
type Address struct {
	gorm.Model
	Pub   Pub
	PubID uint
	Addr  string
}

//go:generate stringer -type=State

type State uint

const (
	Unavailable  State = iota // means something like dns failure or connection refused
	KeyExchanged              // dialed, past the SHS kex
	Muxed                     // did a muxrpc exchange

)

type Check struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Pub       Pub
	PubID     uint
	Addr      Address
	AddrID    uint
	Took      time.Duration
	State     State
	Error     string
}

var (
	PubHealth = transition.New(&Pub{})
)

func init() {
	PubHealth.Initial("unchecked")
	tryEvent := PubHealth.Event("try")
	tryEvent.To("trying").From("failed").Before(checkPub)
	tryEvent.To("trying").From("worked").Before(checkPub)
	tryEvent.To("trying").From("unchecked").Before(checkPub)
	tryEvent.To("trying").From("trying").Before(func(value interface{}, tx *gorm.DB) (err error) {
		return errors.New("TODO")
		/* TODO: deduce from db
		var pub = value.(*Pub)
		if time.Since(pub.LastSuccess) < 2*time.Minute {
			return errors.New("to soon")
		}
		return checkPub(value, tx)
		*/
	})
}

func checkPub(value interface{}, tx *gorm.DB) (err error) {

	var (
		pub    = value.(*Pub)
		pubKey [ed25519.PublicKeySize]byte
		addrs  []Address
		log    = logging.Logger(pub.Key)
	)

	start := time.Now()
	pk, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(strings.TrimPrefix(pub.Key, "@"), ".ed25519"))
	if err != nil {
		return errors.Wrap(err, "secrethandshake: base64 decode of public part failed.")
	}
	copy(pubKey[:], pk)

	shsDialer, err := ssb.SHSClient.NewDialer(pubKey)
	if err != nil {
		return errors.Wrapf(err, "creating dialer for pubkey failed: %q", pub.Key)
	}

	if err := tx.Model(value).Association("Addresses").Find(&addrs).Error; err != nil {
		return errors.Wrap(err, "failed to assoc pub with addresses")
	}

	if len(addrs) < 1 {
		return errors.Errorf("no addresses for %s", pub.Key)
	}

	var (
		errc = make(chan error)
	)

	ctx := context.TODO()

	for i, a := range addrs {
		var checkAddr = func(a Address) error {
			var check Check
			check.PubID = a.PubID
			check.AddrID = a.ID
			log.Log("msg", "dialing", "addr", a.Addr)
			dialStart := time.Now()
			c, err := shsDialer("tcp", a.Addr)
			defer func() {
				if err != nil {
					check.Error = err.Error()
				}
				check.Took = time.Since(dialStart)
				if err = tx.Save(&check).Error; err != nil {
					log.Log("event", "error", "err", err, "msg", "failed to update addresses", "aid", a.ID)
				}
			}()
			if err != nil {
				return errors.Wrapf(err, "dialer(%d) - %s:%s\n", i, pub.Key, a.Addr)
			}
			check.State = KeyExchanged

			counter := debug.WrapCounter(c)
			p := muxrpc.NewPacker(counter)
			handler := ssb.NewTryHandler(pub.Key, kitlog.With(log, "id", pub.Key))
			rpc := muxrpc.Handle(p, handler)

			log.Log("msg", "new rpc client", "addr", a.Addr)

			go func() {
				err := rpc.(muxrpc.Server).Serve(ctx)
				log.Log("event", "connection done",
					"id", pub.Key,
					"err", err,
					"took", time.Since(start),
					"sent", humanize.Bytes(counter.Cw.Count()),
					"rcvd", humanize.Bytes(counter.Cr.Count()),
				)
			}()

			<-handler.Wait
			go func() {
				<-handler.Wait
				<-handler.Wait
			}()

			if err = rpc.Terminate(); err != nil {
				return errors.Wrapf(err, "terminate(%d) - %s:%s", i, pub.Key, a.Addr)
			}

			check.State = Muxed
			return nil
		}

		go func(a Address) {
			errc <- checkAddr(a)
		}(a)
	}

	go func() {
		var success bool
		var errs []error
		for e := range errc {
			if e == nil {
				pub.LastSuccess = time.Now()
				pub.SetState("worked")
				success = true
			}
			errs = append(errs, e)
			if len(errs) == len(addrs) {
				close(errc)
			}
		}

		if !success {
			pub.SetState("failed")
		}
		took := time.Since(start)
		if took > 5*time.Second && success {
			err := notify.Sender.Send(&notification.Message{
				To:          "1",
				Title:       "Long Run!",
				Body:        fmt.Sprintf("Pub#%v attempt worked! Took:%v", pub.ID, took),
				MessageType: "connection_try",
			}, &qor.Context{DB: db.GetBase()})
			if err != nil {
				log.Log("err", errors.Wrap(err, "failed to send notify"))
			}
		}

		if err := tx.Save(pub).Error; err != nil {
			log.Log("err", errors.Wrap(err, "failed to assoc pub with addresses"))
		}
		log.Log("msg", "all checked", "cnt", len(addrs), "success", success, "took", fmt.Sprintf("%v", took))
	}()
	return nil
}
