package service

import (
	"github.com/bool64/brick"
	"github.com/vearutop/form2mail/internal/infra/email"
	"github.com/vearutop/form2mail/internal/infra/recaptcha"
)

// Name is the name of this application or service.
const Name = "form2mail"

// Config defines application configuration.
type Config struct {
	brick.BaseConfig

	SMTP      email.SMTPConfig
	Recaptcha recaptcha.Config
	StaticDir string `split_words:"true"`
}
