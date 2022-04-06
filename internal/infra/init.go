package infra

import (
	"context"
	"net/http"

	"github.com/bool64/brick"
	"github.com/swaggest/rest/response/gzip"
	"github.com/vearutop/form2mail/internal/infra/email"
	"github.com/vearutop/form2mail/internal/infra/recaptcha"
	"github.com/vearutop/form2mail/internal/infra/schema"
	"github.com/vearutop/form2mail/internal/infra/service"
	"go.opencensus.io/plugin/ochttp"
)

// NewServiceLocator creates application service locator.
func NewServiceLocator(cfg service.Config) (loc *service.Locator, err error) {
	l := &service.Locator{}
	l.Config = cfg

	defer func() {
		if err != nil && l != nil && l.LoggerProvider != nil {
			l.CtxdLogger().Error(context.Background(), err.Error())
		}
	}()

	l.BaseLocator, err = brick.NewBaseLocator(cfg.BaseConfig)
	if err != nil {
		return nil, err
	}

	schema.SetupOpenapiCollector(l.OpenAPI)

	l.HTTPServerMiddlewares = append(l.HTTPServerMiddlewares, gzip.Middleware)

	l.RecaptchaCheckerProvider = &recaptcha.V2V3Checker{
		Config: cfg.Recaptcha,
		Transport: &ochttp.Transport{
			Base: http.DefaultTransport,
			FormatSpanName: func(request *http.Request) string {
				return "recaptcha"
			},
		},
	}

	l.EmailSenderProvider, err = email.NewSMTPSender(cfg.SMTP)
	if err != nil {
		return nil, err
	}

	return l, nil
}
