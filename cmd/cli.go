package cmd

import (
	"fmt"

	"github.com/Xuanwo/go-locale"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
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
	tag, err := locale.Detect()
	if err != nil {
		tag = language.English
	}
	p = message.NewPrinter(tag)
	return
}

func versionPrinter(c *cli.Context) {
	fmt.Println(c.App.Version)
}

func New(version string) *cli.App {
	subs := []*cli.Command{UpdateCmd(), ConvertCmd(), CountCmd()}

	flags := []cli.Flag{debugFlag}

	cli.VersionPrinter = versionPrinter

	cmd := &cli.App{
		Name:        "oui",
		Usage:       "MAC Address CLI Toolkit",
		Action:      MainCmd,
		Commands:    subs,
		Flags:       flags,
		Version:     version,
		HideVersion: false,
	}

	return cmd
}
