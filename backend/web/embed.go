package web

import (
	"embed"
)

// EmbeddedFiles contains the embedded web assets.
// The **/* pattern includes all files and subdirectories.
// This includes the fallback index.html for when the full web app is not available.
//
//go:embed **/* all:public
var EmbeddedFiles embed.FS
