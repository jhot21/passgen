package web

import "embed"

//go:embed all:app/dist
var staticFiles embed.FS
