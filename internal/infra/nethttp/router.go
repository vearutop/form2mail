// Package nethttp manages application http interface.
package nethttp

import (
	"net/http"

	"github.com/bool64/brick"
	"github.com/swaggest/rest/nethttp"
	"github.com/vearutop/form2mail/internal/infra/nethttp/ui"
	"github.com/vearutop/form2mail/internal/infra/service"
	"github.com/vearutop/form2mail/internal/usecase"
)

// NewRouter creates an instance of router filled with handlers and docs.
func NewRouter(deps *service.Locator) http.Handler {
	r := brick.NewBaseRouter(deps.BaseLocator)

	r.Method(http.MethodPost, "/receive", nethttp.NewHandler(usecase.Receive(deps),
		nethttp.SuccessStatus(http.StatusSeeOther)))

	r.Method(http.MethodGet, "/", ui.Index(deps.Config))
	r.Mount("/static/", http.StripPrefix("/static", ui.Static(deps.Config)))

	return r
}
