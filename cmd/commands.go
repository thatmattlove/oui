package cmd

import (
	"fmt"

	"github.com/gookit/gcli/v3/progress"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/thatmattlove/go-macaddr"
	"github.com/thatmattlove/oui/v2/internal/logger"
	"github.com/thatmattlove/oui/v2/internal/util"
	"github.com/thatmattlove/oui/v2/oui"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func resultsTitle(search string) string {
	s := &text.Colors{text.FgHiCyan, text.Bold, text.Underline}
	t := &text.Colors{text.Bold}
	return fmt.Sprintf("\n %s%s\n", s.Sprint(search), t.Sprint(" Results"))
}

func createTable(search string) (t table.Writer) {
	t = table.NewWriter()
	t.SetStyle(*tableStyle)
	h := &text.Colors{text.FgHiMagenta, text.Bold}
	t.AppendHeader(table.Row{h.Sprint("Prefix"), h.Sprint("Organization"), h.Sprint("Range"), h.Sprint("Registry")})
	return
}

func CountCmd() *cli.Command {
	return &cli.Command{
		Name:    "entires",
		Usage:   "Show the number of MAC addresses in the database",
		Aliases: []string{"e", "count"},
		Action: func(c *cli.Context) error {
			logger := logger.New()
			sqlite, err := oui.CreateSQLiteOption()
			if err != nil {
				return err
			}
			db, err := oui.New(oui.WithVersion(c.App.Version), oui.WithLogging(logger), sqlite)
			if err != nil {
				return err
			}
			count, err := db.Count()
			if err != nil {
				return err
			}
			logger.Info("MAC Address database has %s entries", count)
			return nil
		},
	}
}

func MainCmd(c *cli.Context) error {
	log := logger.New()
	search := c.Args().First()

	sqlite, err := oui.CreateSQLiteOption()
	if err != nil {
		return err
	}

	db, err := oui.New(oui.WithVersion(c.App.Version), oui.WithLogging(log), sqlite)
	if err != nil {
		return err
	}
	defer db.Close()
	count, err := db.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		log.Warn("MAC Address database has not been populated.")
		err = UpdateCmd().Run(c)
		if err != nil {
			return err
		}
	}

	results, err := db.Find(search)
	if err != nil {
		return err
	}

	t := createTable(search)
	fmt.Println(resultsTitle(search))

	if len(results) == 0 {
		log.Error("\n No results found\n\n")
	} else {
		for _, result := range results {
			_, mp, err := macaddr.ParseMACPrefix(result.PrefixString())
			if err != nil {
				return err
			}
			m := (&text.Colors{text.FgHiGreen, text.Bold}).Sprint(mp.MAC.String())
			p := text.FgHiCyan.Sprintf("/%d", mp.PrefixLen())
			rf := text.FgHiCyan.Sprint(mp.First())
			rl := text.FgHiRed.Sprint(mp.Last())
			r := fmt.Sprintf("%s-%s", rf, rl)
			t.AppendRow(table.Row{m + p, result.Org, r, result.Registry})
		}
		fmt.Println(t.Render())
	}
	return nil
}

func UpdateCmd() *cli.Command {
	cmd := &cli.Command{
		Name:    "update",
		Usage:   "Refresh the MAC address database",
		Aliases: []string{"u", "up"},
	}
	cmd.Action = func(c *cli.Context) error {
		statuses := map[int]string{
			5:   "Downloading vendor data...",
			10:  "Processing vendor data...",
			99:  "Populating database...",
			100: "Completed",
		}

		style := progress.BarChars{Completed: '█', Processing: '▌', Remaining: '░'}
		b := (&text.Colors{text.Color(808080)}).Sprint("{@bar}")
		title := (&text.Colors{text.FgHiCyan, text.Bold}).Sprint("\nUpdating MAC Address Database")

		fmt.Println(title)

		p := progress.New(100).
			Config(func(p *progress.Progress) {
				p.Format = b + " {@percent:4s}% {@message}"
			}).
			AddWidget("bar", progress.BarWidget(50, style)).
			AddWidget("message", progress.DynamicTextWidget(statuses))
		p.Start()
		sqlite, err := oui.CreateSQLiteOption()
		if err != nil {
			return err
		}
		db, err := oui.New(oui.WithVersion(c.App.Version), oui.WithLogging(logger.New()), oui.WithProgress(p), sqlite)
		if err != nil {
			return err
		}
		defer db.Close()
		p.AdvanceTo(3)
		num, err := db.Populate()
		p.Finish()
		if err != nil {
			return err
		}
		message.NewPrinter(language.English)

		dur := util.TimeSince(p.StartedAt())
		v := (&text.Colors{text.FgHiGreen, text.Bold}).Sprint(c.App.Version)
		n := (&text.Colors{text.FgHiBlue, text.Bold}).Sprint(withLocale().Sprint(num))
		d := (&text.Colors{text.FgHiRed, text.Bold}).Sprint(dur)
		fmt.Printf("Updated MAC Address database (%s) with %s records in %s\n", v, n, d)
		return nil
	}

	return cmd
}

func ConvertCmd() *cli.Command {
	cmd := &cli.Command{
		Name:      "convert",
		Usage:     "Convert a MAC Address to other formats",
		Aliases:   []string{"c", "con"},
		ArgsUsage: "MAC Address to Convert",
	}
	cmd.Action = func(c *cli.Context) error {
		i := c.Args().First()
		f, err := oui.Convert(i)
		if err != nil {
			return err
		}
		t := table.NewWriter()
		t.SetStyle(*tableStyle)
		h := (&text.Colors{text.FgHiMagenta, text.Bold})
		t.AppendRow(table.Row{h.Sprint("Hexadecimal"), f.Hex})
		t.AppendRow(table.Row{h.Sprint("Dotted"), f.Dotted})
		t.AppendRow(table.Row{h.Sprint("Dashed"), f.Dashed})
		t.AppendRow(table.Row{h.Sprint("Integer"), f.Int})
		t.AppendRow(table.Row{h.Sprint("Bytes"), f.Bytes})
		m := (&text.Colors{text.FgHiCyan, text.Bold, text.Underline}).Sprintf("%s\n", i)
		fmt.Println("\n " + m)
		fmt.Println(t.Render())
		return nil
	}

	return cmd
}
