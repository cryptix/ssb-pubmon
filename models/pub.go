package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/agl/ed25519"
	"github.com/cryptix/go-muxrpc"
	"github.com/cryptix/go/logging"
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

type Address struct {
	/* Composite primary keys can't be created on SQLite tables
	    https://github.com/jinzhu/gorm/issues/1037

	   	Host  string `gorm:"primary_key"`
	   	Port  int    `gorm:"primary_key"`
	*/
	gorm.Model
	Pub      Pub
	PubID    uint
	Addr     string
	Took     time.Duration
	Failures int
	LastTry  string
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
		var pub = value.(*Pub)
		if time.Since(pub.LastSuccess) < 2*time.Minute {
			return errors.New("to soon")
		}
		return checkPub(value, tx)
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

	for i, a := range addrs {
		var checkAddr = func(a Address) error {
			log.Log("msg", "dialing", "addr", a.Addr)
			dialStart := time.Now()
			c, err := shsDialer("tcp", a.Addr)
			defer func() {
				if err != nil {
					a.LastTry = err.Error()
					a.Failures++
				}
				a.Took = time.Since(dialStart)
				if err = tx.Save(a).Error; err != nil {
					log.Log("event", "error", "err", err, "msg", "failed to update addresses", "aid", a.ID)
				}
			}()
			if err != nil {
				return errors.Wrapf(err, "dialer(%d) - %s:%s\n", i, pub.Key, a.Addr)
			}

			rpc := muxrpc.NewClient(log, c)
			go rpc.Handle()

			log.Log("msg", "new rpc client", "addr", a.Addr)

			wait := make(chan struct{})

			rpc.HandleCall("gossip.ping", func(msg json.RawMessage) interface{} {
				wait <- struct{}{}
				return nil
			})

			rpc.HandleSource("blobs.createWants", func(msg json.RawMessage) chan interface{} {
				wait <- struct{}{}
				return nil
			})

			rpc.HandleSource("createHistoryStream", func(msg json.RawMessage) chan interface{} {
				wait <- struct{}{}
				return nil
			})

			<-wait
			go func() {
				<-wait
				<-wait
				close(wait)
			}()

			if err = rpc.Close(); err != nil {
				return errors.Wrapf(err, "close(%d) - %s:%s", i, pub.Key, a.Addr)
			}

			a.LastTry = "success"
			a.Failures = 0
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
			switch {
			case e == nil:
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
