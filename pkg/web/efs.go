package web

import "embed"

//go:embed "static/js"
var JSFiles embed.FS

//go:embed "static/css"
var CSSFiles embed.FS
