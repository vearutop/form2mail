// Package ui provides application web user interface.
package ui

import (
	"bytes"
	"github.com/vearutop/form2mail/internal/infra/service"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/vearutop/form2mail/resources/static"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

// Static serves static assets.
func Static(cfg service.Config) http.Handler {
	if cfg.StaticDir != "" {
		return http.FileServer(http.Dir(cfg.StaticDir))
	}

	if _, err := os.Stat("./resources/static"); err == nil {
		// path/to/whatever exists
		return http.FileServer(http.Dir("./resources/static"))
	}

	return statigz.FileServer(static.Assets, brotli.AddEncoding, statigz.EncodeOnInit)
}

// Index serves index page of the application.
func Index(cfg service.Config) http.Handler {
	var (
		file  string
		index []byte
		err   error
	)

	switch {
	case cfg.Recaptcha.V3 && cfg.Recaptcha.SiteKey != "":
		file = "recaptcha_v3.html"
	case cfg.Recaptcha.SiteKey != "":
		file = "recaptcha_v2.html"
	default:
		file = "no_recaptcha.html"
	}

	if cfg.StaticDir != "" {
		index, err = ioutil.ReadFile(cfg.StaticDir + "/index.html")

	} else {
		index, err = static.Assets.ReadFile(file)
	}
	if err != nil {
		panic(err)
	}

	index = bytes.Replace(index, []byte("RECAPTCHA_SITE_KEY"), []byte(cfg.Recaptcha.SiteKey), 1)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(index)
	})
}
