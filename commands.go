package main

import (
	"fmt"

	"github.com/gookit/gcli/v3/progress"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/thatmattlove/go-macaddr"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func resultsTitle(search string) string {
	s := &text.Colors{text.FgCyan, text.Bold, text.Underline}
	t := &text.Colors{text.Bold}
	return fmt.Sprintf("\n %s%s\n", s.Sprint(search), t.Sprint(" Results"))
}

func createTable(search string) (t table.Writer) {
	t = table.NewWriter()
	t.SetStyle(*tableStyle)
	h := &text.Colors{text.FgMagenta, text.Bold}
	t.AppendHeader(table.Row{h.Sprint("Prefix"), h.Sprint("Organization"), h.Sprint("Range")})
	return
}

func MainCmd(c *cli.Context) error {
	MaybePanic = createPanicFunc()
	search := c.Args().First()

	t := createTable(search)
	fmt.Println(resultsTitle(search))

	for result := range Find(search) {
		_, mp, err := macaddr.ParseMACPrefix(result.PrefixString())
		MaybePanic(err)
		m := (&text.Colors{text.FgHiGreen, text.Bold}).Sprint(mp.MAC.String())
		p := text.FgCyan.Sprintf("/%d", mp.PrefixLen())
		rf := text.FgCyan.Sprint(mp.First())
		rl := text.FgBlue.Sprint(mp.Last())
		r := fmt.Sprintf("%s-%s", rf, rl)
		t.AppendRow(table.Row{m + p, result.Org, r})
	}
	fmt.Println(t.Render())
	return nil
}

func UpdateCmd() *cli.Command {
	cmd := &cli.Command{
		Name:    "update",
		Usage:   "Refresh the MAC address database",
		Aliases: []string{"u", "up"},
	}
	cmd.Action = func(c *cli.Context) error {
		MaybePanic = createPanicFunc()
		statuses := map[int]string{
			10:  "Deleted old database...",
			20:  "Created new database...",
			30:  "Downloading vendor data...",
			40:  "Processing vendor data...",
			99:  "Populating database...",
			100: "Completed",
		}

		style := progress.BarChars{Completed: '█', Processing: '▌', Remaining: '░'}
		b := (&text.Colors{text.Color(808080)}).Sprint("{@bar}")
		title := (&text.Colors{text.FgCyan, text.Bold}).Sprint("\nUpdating MAC Address Database")

		fmt.Println(title)

		p := progress.New(100).
			Config(func(p *progress.Progress) {
				p.Format = b + " {@percent:4s}% {@message}"
			}).
			AddWidget("bar", progress.BarWidget(50, style)).
			AddWidget("message", progress.DynamicTextWidget(statuses))
		num, err := UpdateTable(p)
		p.Finish()
		MaybePanic(err)
		message.NewPrinter(language.English)

		dur := timeSince(p.StartedAt())
		v := (&text.Colors{text.FgGreen, text.Bold}).Sprint(_tableVersion)
		n := (&text.Colors{text.FgBlue, text.Bold}).Sprint(withLocale().Sprint(num))
		d := (&text.Colors{text.FgRed, text.Bold}).Sprint(dur)
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
		MaybePanic = createPanicFunc()
		i := c.Args().First()
		f := Convert(i)
		t := table.NewWriter()
		t.SetStyle(*tableStyle)
		h := (&text.Colors{text.FgMagenta, text.Bold})
		t.AppendRow(table.Row{h.Sprint("Hexadecimal"), f.Hex})
		t.AppendRow(table.Row{h.Sprint("Dotted"), f.Dotted})
		t.AppendRow(table.Row{h.Sprint("Dashed"), f.Dashed})
		t.AppendRow(table.Row{h.Sprint("Integer"), f.Int})
		t.AppendRow(table.Row{h.Sprint("Bytes"), f.Bytes})
		m := (&text.Colors{text.FgCyan, text.Bold, text.Underline}).Sprintf("%s\n", i)
		fmt.Println("\n " + m)
		fmt.Println(t.Render())
		return nil
	}

	return cmd
}
