package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
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

	return config.View.Execute("pubalive", map[string]interface{}{
		"Pub":   p,
		"Tries": tries,
		"humanizeTime": func(when time.Time) string {
			return humanize.Time(when)
		},
	}, req, w)
}

func SwitchLocale(w http.ResponseWriter, req *http.Request) {
	qorutils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
}
