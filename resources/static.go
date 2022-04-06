// Package static provides embedded static assets.
package resources

import (
	"embed"
)

// Static provides embedded static assets for web application.
//go:embed static/*
var Static embed.FS
