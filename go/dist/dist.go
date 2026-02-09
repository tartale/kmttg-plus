package dist

import "embed"

//go:embed webui/*
var Filesystem embed.FS
