package main // import "github.com/cryptix/ssb-pubmon"

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cryptix/go/logging"
	"github.com/cryptix/ssb-pubmon/config"
	"github.com/cryptix/ssb-pubmon/config/admin/bindatafs"
	_ "github.com/cryptix/ssb-pubmon/db/migrations"
	"github.com/cryptix/ssb-pubmon/sbmhttp"
)

var (
	log   logging.Interface
	check = logging.CheckFatal

	Revision string = "undefined"
)

func main() {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	cmdLine.Parse(os.Args[1:])

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

	var h = sbmhttp.InitServ(log, Revision)

	if *compileTemplate {
		bindatafs.AssetFS.Compile()
		return
	}
	log.Log("event", "init", "msg", "http listen", "addr", config.Config.HTTPHost, "version", Revision)
	check(http.ListenAndServe(config.Config.HTTPHost, h))
}
