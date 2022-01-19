package main

import (
	_ "embed"

	"golang.org/x/text/language"
)

const (
	_ouiListUrl   string = "https://gitlab.com/wireshark/wireshark/-/raw/master/manuf"
	_tableVersion string = "v0"
	_nilStr       string = "<nil>"
	_updateMsg    string = "table '%s' is missing from the database"
)

var MaybePanic func(err error)

var _locale language.Tag

//go:embed .version
var Version string
