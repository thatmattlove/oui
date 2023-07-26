package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/message"
)

var tableStyle = &table.Style{
	Name:    "StyleRounded",
	Box:     table.StyleBoxRounded,
	Color:   table.ColorOptionsDefault,
	Options: table.OptionsDefault,
	Title:   table.TitleOptionsDefault,
}

var debugFlag *cli.BoolFlag = &cli.BoolFlag{Name: "debug", Usage: "Enable debugging", Value: false}

func withLocale() (p *message.Printer) {
	p = message.NewPrinter(_locale)
	return
}

func versionPrinter(c *cli.Context) {
	fmt.Println(c.App.Version)
}

func CLI() *cli.App {
	subs := []*cli.Command{UpdateCmd(), ConvertCmd(), CountCmd()}

	flags := []cli.Flag{debugFlag}

	cli.VersionPrinter = versionPrinter

	cmd := &cli.App{
		Name:        "oui",
		Usage:       "MAC Address CLI Toolkit",
		Action:      MainCmd,
		Commands:    subs,
		Flags:       flags,
		Version:     Version,
		HideVersion: false,
	}

	return cmd
}
