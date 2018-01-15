package controllers

import (
	"fmt"
	"image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/qor/qor"
	qorutils "github.com/qor/qor/utils"

	"github.com/cryptix/ssb-pubmon/config"
	sbmutils "github.com/cryptix/ssb-pubmon/config/utils"
	"github.com/cryptix/ssb-pubmon/hexagen"
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

func Overview(w http.ResponseWriter, req *http.Request) error {
	db := sbmutils.GetDB(req)

	var worked []models.Pub
	if err := db.Limit(5).Order("last_success desc").Where("state = ?", "worked").Find(&worked).Error; err != nil {
		return errors.Wrap(err, "overview: worked qry failed")
	}

	var failing []models.Address
	if err := db.Preload("Pub").Limit(5).Order("failures desc").Find(&failing).Error; err != nil {
		return errors.Wrap(err, "overview: failing qry failed")
	}

	err := config.View.Execute("overview", map[string]interface{}{
		"worked": worked,
		"since": func(when time.Time) string {
			return fmt.Sprintf("%v", time.Since(when))
		},
		"failing": failing,
	}, req, w)
	if err != nil {
		return errors.Wrap(err, "overview: tmpl execute failed")
	}
	return nil
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
	fmt.Fprintf(w, "%v", p)
	return nil
}

func Hexagen(w http.ResponseWriter, req *http.Request) error {
	width, err := strconv.ParseFloat(req.URL.Query().Get("width"), 64)
	if err != nil {
		return &BadReqError{req: req, field: "width", msg: " illegal width value"}
	}
	if width < 0 || width > 2048 {
		width = 512
	}

	key := req.URL.Query().Get("key")
	g, err := hexagen.Generate(key, width)
	if err != nil {
		return errors.Wrap(err, "hexagen: failed to generate image")
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, g); err != nil {
		return errors.Wrap(err, "hexagen: png encoding failed")
	}
	return nil
}

func SwitchLocale(w http.ResponseWriter, req *http.Request) {
	qorutils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
}
