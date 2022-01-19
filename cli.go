package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
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

func createPanicFunc() func(err error) {
	return func(err error) {
		isTest := flag.Lookup("test.v") != nil
		if err != nil {
			if isTest {
				panic(err)
			} else {
				e := &text.Colors{text.FgYellow, text.Bold}
				fmt.Println(e.Sprint(err))
				os.Exit(1)
			}
		}
	}
}

func withLocale() (p *message.Printer) {
	p = message.NewPrinter(_locale)
	return
}

func versionPrinter(c *cli.Context) {
	fmt.Println(c.App.Version)
}

func CLI() *cli.App {
	subs := []*cli.Command{UpdateCmd(), ConvertCmd()}

	cli.VersionPrinter = versionPrinter

	cmd := &cli.App{
		Name:        "oui",
		Usage:       "MAC Address CLI Toolkit",
		Action:      MainCmd,
		Commands:    subs,
		Version:     Version,
		HideVersion: false,
	}

	return cmd
}
