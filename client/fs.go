package client

import "embed"

//go:embed all:build/*
var DistFS embed.FS
