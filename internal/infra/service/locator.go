package service

import (
	"github.com/bool64/brick"
)

// Locator defines application resources.
type Locator struct {
	*brick.BaseLocator

	Config Config

	RecaptchaCheckerProvider
	EmailSenderProvider
}

// ServiceConfig returns service configuration.
func (l *Locator) ServiceConfig() Config {
	return l.Config
}
