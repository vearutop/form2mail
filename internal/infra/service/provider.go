package service

import (
	"github.com/vearutop/form2mail/internal/infra/email"
	"github.com/vearutop/form2mail/internal/infra/recaptcha"
)

type RecaptchaCheckerProvider interface {
	RecaptchaChecker() recaptcha.Checker
}

type EmailSenderProvider interface {
	EmailSender() email.Sender
}
