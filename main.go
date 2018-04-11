package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/config/admin/bindatafs"
	"github.com/cryptix/ssb-pubmon/db"
	_ "github.com/cryptix/ssb-pubmon/db/migrations"
	"github.com/cryptix/ssb-pubmon/models"
	"github.com/cryptix/ssb-pubmon/sbmhttp"
	"github.com/cryptix/ssb-pubmon/ssb"
)

var (
	log   logging.Interface
	check = logging.CheckFatal

	Revision string = "undefined"

	flagCompileTemplates = flag.Bool("compile-templates", false, "Compile Templates")
	flagKeyfile          = flag.String("keyfile", "secret", "ssb keypair file (#json)")
)

func main() {
	flag.Parse()

	// create timestamped logfile
	os.Mkdir("logs", 0700)
	logFileName := fmt.Sprintf("logs/%s-%s.log",
		filepath.Base(os.Args[0]),
		time.Now().Format("2006-01-02_15-04"))
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic(err) // logging not ready yet...
	}
	logging.SetupLogging(io.MultiWriter(os.Stderr, logFile))
	log = logging.Logger("ssb-pubmon")

	err = ssb.InitClient(*flagKeyfile)
	logging.CheckFatal(err)

	var h = sbmhttp.InitServ(log, Revision)

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		db := db.GetBase()
		for t := range ticker.C {
			start := time.Now()
			var pubs []models.Pub
			rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
			const qry = `SELECT * FROM pubs WHERE id NOT IN (SELECT DISTINCT pub_id FROM checks WHERE created_at > ?) ORDER BY RANDOM()`
			err := db.Raw(qry, rounded).Scan(&pubs).Error
			logging.CheckFatal(err)
			const limit = 15 // limit fan-out
			if len(pubs) > limit {
				pubs = pubs[:limit]
			}
			var wg sync.WaitGroup
			wg.Add(len(pubs))
			for _, p := range pubs {
				go func(pub models.Pub, d *sync.WaitGroup) {
					err := models.CheckPub(&pub, db)
					if err != nil {
						log.Log("CheckPub", "error", "pub", pub.ID, "err", err)
					}
					d.Done()
				}(p, &wg)
			}
			wg.Wait()
			if len(pubs) > 0 {
				log.Log("before", rounded, "pub#", len(pubs), "checkTick", time.Since(start))
			}
		}
	}()

	if *flagCompileTemplates {
		bindatafs.AssetFS.Compile()
		return
	}
	log.Log("event", "init", "msg", "http listen", "addr", config.Config.HTTPHost, "version", Revision)
	check(http.ListenAndServe(config.Config.HTTPHost, h))
}
