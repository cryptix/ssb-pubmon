package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/agl/ed25519"
	"github.com/prasannavl/go-errors"
	"github.com/qor/media/oss"
	"github.com/qor/worker"
	goon "github.com/shurcooL/go-goon"

	"github.com/cryptix/ssb-pubmon/db"
	"github.com/cryptix/ssb-pubmon/models"
)

type SimpleQueue struct {
	Worker *worker.Worker
}

func (q *SimpleQueue) Add(j worker.QorJobInterface) error {
	return q.Worker.RunJob(j.GetJobID())
}

func (q *SimpleQueue) Run(j worker.QorJobInterface) error {
	job := j.GetJob()

	if job.Handler != nil {
		err := job.Handler(j.GetSerializableArgument(j), j)
		if err != nil {
			j.AddLog(err.Error())
			return err

		}
		return nil
	}

	return errors.New("SimpleQue:no handler found for job " + job.Name)
}

func (q *SimpleQueue) Kill(j worker.QorJobInterface) error {
	return errors.New("SimpleQue:kill not implemented")
}

func (q *SimpleQueue) Remove(j worker.QorJobInterface) error {
	return errors.New("SimpleQue:remove not implemented")
}

func getWorker() *worker.Worker {
	sq := SimpleQueue{}
	Worker := worker.New(&worker.Config{
		Queue: &sq,
	})
	sq.Worker = Worker

	type sendNewsletterArgument struct {
		Subject      string
		Content      string `sql:"size:65532"`
		SendPassword string
		worker.Schedule
	}

	Worker.RegisterJob(&worker.Job{
		Name: "Send Newsletter",
		Handler: func(argument interface{}, qorJob worker.QorJobInterface) error {
			qorJob.AddLog("Started sending newsletters...")
			qorJob.AddLog(fmt.Sprintf("Argument: %+v", argument.(*sendNewsletterArgument)))
			for i := 1; i <= 100; i++ {
				time.Sleep(100 * time.Millisecond)
				qorJob.AddLog(fmt.Sprintf("Sending newsletter %v...", i))
				qorJob.SetProgress(uint(i))
			}
			qorJob.AddLog("Finished send newsletters")
			return nil
		},
		Resource: Admin.NewResource(&sendNewsletterArgument{}),
	})

	type importGossipJSONArg struct {
		File oss.OSS
	}

	Worker.RegisterJob(&worker.Job{
		Name: "Import gossip.json",
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
			argument := arg.(*importGossipJSONArg)

			d := db.DB
			d.LogMode(true)
			defer d.LogMode(false)

			loc := filepath.Join("public", argument.File.URL())
			f, err := os.Open(loc)
			if err != nil {
				return errors.NewWithCause("failed to open upload", err)
			}
			defer f.Close()

			type gossipAddr struct {
				Host, Key string
				Port      int
			}
			var pubs []gossipAddr
			if err := json.NewDecoder(f).Decode(&pubs); err != nil {
				return errors.NewWithCause("failed to decode uploaded file", err)
			}

			var heads = []worker.TableCell{
				{Value: "i"},
				{Value: "Host"},
				{Value: "Port"},
				{Value: "Key"},
			}
			qorJob.AddResultsRow(heads...)

			for i, jsonPub := range pubs {
				qorJob.AddResultsRow([]worker.TableCell{
					{Value: fmt.Sprint(i)},
					{Value: jsonPub.Host},
					{Value: fmt.Sprint(jsonPub.Port)},
					{Value: fmt.Sprint(jsonPub.Key)},
				}...)

				key, err := base64.StdEncoding.DecodeString(
					strings.TrimSuffix(strings.TrimPrefix(jsonPub.Key, "@"), ".ed25519"))
				if err != nil {
					return errors.NewWithCause("base64 decode of public part failed", err)
				}

				if len(key) != ed25519.PublicKeySize {
					return errors.New("illegal pubkey size")
				}

				var p models.Pub
				p.Key = jsonPub.Key

				if err := d.FirstOrCreate(&p, p).Error; err != nil {
					return errors.NewWithCause("could not find/create Pub", err)
				}

				var a models.Address
				a.PubID = p.ID
				a.Host = jsonPub.Host
				a.Port = jsonPub.Port
				goon.Dump(a)

				if err := d.FirstOrCreate(&a, a).Error; err != nil {
					return errors.NewWithCause("could not find/create Addr", err)
				}

				qorJob.AddLog(fmt.Sprintf("Processing %d: pub(%s)", i, p.Key))

			}

			qorJob.AddLog(fmt.Sprintf("Parsed %d Pub announcments", len(pubs)))
			return nil
		},
		Resource: Admin.NewResource(&importGossipJSONArg{}),
	})

	return Worker
}
