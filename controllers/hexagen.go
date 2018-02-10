package controllers

import (
	"image/png"
	"net/http"
	"strconv"

	"cryptoscope.co/go/hexagen"
	"github.com/pkg/errors"
)

func Hexagen(w http.ResponseWriter, req *http.Request) error {
	wStr := req.URL.Query().Get("width")
	width, err := strconv.ParseFloat(wStr, 64)
	if err != nil {
		return &BadReqError{req: req, field: "width", msg: " illegal width value"}
	}
	if width < 0 || width > 2048 {
		width = 512
	}

	key := req.URL.Query().Get("key")
	etag := wStr + key

	if match := req.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	g, err := hexagen.Generate(key, width)
	if err != nil {
		return errors.Wrap(err, "hexagen: failed to generate image")
	}

	w.Header().Set("Etag", etag)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=31557600")
	if err := png.Encode(w, g); err != nil {
		return errors.Wrap(err, "hexagen: png encoding failed")
	}
	return nil
}
