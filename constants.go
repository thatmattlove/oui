package main

import (
	_ "embed"

	"golang.org/x/text/language"
)

const (
	_tableVersion string = "v1"
	_nilStr       string = "<nil>"
	_updateMsg    string = "table '%s' is missing from the database"
)

var _locale language.Tag

//go:embed .version
var Version string
