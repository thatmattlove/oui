package main

import (
	_ "embed"
)

//go:embed .version
var Version string

const TABLE_VERSION string = "v1"
