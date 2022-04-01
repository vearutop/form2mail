// Package ui provides application web user interface.
package ui

import (
	"bytes"
	"github.com/vearutop/form2mail/internal/infra/service"
	"net/http"
	"os"

	"github.com/vearutop/form2mail/resources/static"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

// Static serves static assets.
var Static http.Handler

// nolint:gochecknoinits
func init() {
	if _, err := os.Stat("./resources/static"); err == nil {
		// path/to/whatever exists
		Static = http.FileServer(http.Dir("./resources/static"))
	} else {
		Static = statigz.FileServer(static.Assets, brotli.AddEncoding, statigz.EncodeOnInit)
	}
}

// Index serves index page of the application.
func Index(cfg service.Config) http.Handler {
	var file string

	switch {
	case cfg.Recaptcha.V3 && cfg.Recaptcha.SiteKey != "":
		file = "recaptcha_v3.html"
	case cfg.Recaptcha.SiteKey != "":
		file = "recaptcha_v2.html"
	default:
		file = "no_recaptcha.html"
	}

	if cfg.Recaptcha.V3 {
		file = "recaptcha_v3.html"
	}

	index, err := static.Assets.ReadFile(file)
	if err != nil {
		panic(err)
	}

	index = bytes.Replace(index, []byte("RECAPTCHA_SITE_KEY"), []byte(cfg.Recaptcha.SiteKey), 1)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(index)
	})
}
