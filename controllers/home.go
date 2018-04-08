package controllers

import (
	"fmt"
	"net/http"
	"text/tabwriter"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/qor/qor"
	qorutils "github.com/qor/qor/utils"

	"github.com/cryptix/ssb-pubmon/config"
	sbmutils "github.com/cryptix/ssb-pubmon/config/utils"
	"github.com/cryptix/ssb-pubmon/models"
)

type BadReqError struct {
	req        *http.Request
	field, msg string
}

func (e *BadReqError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.msg)
}

func Index(w http.ResponseWriter, req *http.Request) error {
	err := config.View.Execute("index", nil, req, w)
	if err != nil {
		return errors.Wrap(err, "overview: tmpl execute failed")
	}
	return nil
}

func LastChecks(w http.ResponseWriter, req *http.Request) error {
	db := sbmutils.GetDB(req)

	var checks []models.Check
	if err := db.Preload("Pub").Preload("Addr").Limit(25).Order("created_at desc").Find(&checks).Error; err != nil {
		return errors.Wrap(err, "lastchecks: checks qry failed")
	}

	return config.View.Execute("lastchecks", map[string]interface{}{
		"checks": checks,
		"since": func(when time.Time) string {
			return fmt.Sprintf("%v", time.Since(when))
		},
	}, req, w)
}

func Alive(w http.ResponseWriter, req *http.Request) error {
	db := sbmutils.GetDB(req)
	key := req.URL.Query().Get("key")
	if len(key) < 50 || key[0] != '@' {
		return &BadReqError{req: req, field: "key", msg: "illegal key value"}
	}

	var p models.Pub
	if err := db.Find(&p, "key = ?", key).Error; err != nil {
		return errors.Wrap(err, "alive: worked qry failed")
	}

	var tries []models.Check
	if err := db.Preload("Pub").Preload("Addr").Limit(100).Order("created_at desc").Find(&tries, "pub_id = ?", p.ID).Error; err != nil {
		return errors.Wrap(err, "alive: tries qry failed")
	}

	tw := tabwriter.NewWriter(w, 0, 0, '\t', ' ', 0)
	fmt.Fprintln(tw, "#\tAddr\tState\tSaved\tTook\tError\t")
	for i, try := range tries {
		_, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", i,
			try.Addr.Addr,
			try.State,
			humanize.Time(try.CreatedAt),
			try.Took,
			try.Error)
		if err != nil {
			return errors.Wrap(err, "alive: printing tries failed")
		}
	}
	return tw.Flush()
}

func SwitchLocale(w http.ResponseWriter, req *http.Request) {
	qorutils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
}
