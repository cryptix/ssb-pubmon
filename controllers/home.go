package controllers

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/qor/qor"
	qorutils "github.com/qor/qor/utils"

	"github.com/cryptix/ssb-pubmon/config"
	sbmutils "github.com/cryptix/ssb-pubmon/config/utils"
	"github.com/cryptix/ssb-pubmon/hexagen"
	"github.com/cryptix/ssb-pubmon/models"
)

func HomeIndex(w http.ResponseWriter, req *http.Request) {

	db := sbmutils.GetDB(req)

	var worked []models.Pub
	if err := db.Limit(5).Order("last_success desc").Where("state = ?", "worked").Find(&worked).Error; err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var failing []models.Address
	if err := db.Preload("Pub").Limit(5).Order("failures desc").Find(&failing).Error; err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := config.View.Execute("home_index", map[string]interface{}{
		"worked": worked,
		"since": func(when time.Time) string {
			return fmt.Sprintf("%v", time.Since(when))
		},
		"failing": failing,
	}, req, w)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Hexagen(w http.ResponseWriter, req *http.Request) {
	width, err := strconv.ParseFloat(req.URL.Query().Get("width"), 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "error: illegal width value")
		return
	}
	if width < 0 || width > 2048 {
		width = 512
	}

	key := req.URL.Query().Get("key")
	g, err := hexagen.Generate(key, width)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, g); err != nil {
		log.Println(err)
		return
	}
}

func SwitchLocale(w http.ResponseWriter, req *http.Request) {
	qorutils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
}
