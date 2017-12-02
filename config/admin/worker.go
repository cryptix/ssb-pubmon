package admin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cryptix/ssb-pubmon/models"
	"github.com/prasannavl/go-errors"
	"github.com/qor/media/oss"
	"github.com/qor/worker"
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

			//context := &qor.Context{DB: db.DB}

			loc := filepath.Join("public", argument.File.URL())

			f, err := os.Open(loc)
			if err != nil {
				return errors.NewWithCause("failed to open upload", err)
			}
			defer f.Close()
			var pubs []models.Pub
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

			for i, p := range pubs {
				qorJob.AddResultsRow([]worker.TableCell{
					{Value: fmt.Sprint(i)},
					{Value: p.Host},
					{Value: fmt.Sprint(p.Port)},
					{Value: fmt.Sprint(p.Key)},
				}...)
			}

			qorJob.AddLog(fmt.Sprintf("Parsed %d Pub announcments", len(pubs)))
			return nil
		},
		Resource: Admin.NewResource(&importGossipJSONArg{}),
	})

	return Worker
}
