package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/gcli/v3/progress"
	"github.com/thatmattlove/go-macaddr"
)

func processBody(b string) (l []string) {
	for _, line := range strings.Split(b, "\n") {
		line = removeComments(line)
		if line != "" {
			l = append(l, line)
		}
	}
	return
}

func DownloadFile(dir string, p *progress.Progress) (fn string, count int, err error) {
	tf, err := os.CreateTemp(dir, "data-*.txt")
	fn = tf.Name()
	if err != nil {
		return
	}
	p.AdvanceTo(30)
	res, err := http.Get(_ouiListUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	defer tf.Close()
	p.AdvanceTo(40)
	b, err := io.ReadAll(res.Body)
	lines := processBody(string(b))
	count = len(lines)
	for _, l := range lines {
		l += "\n"
		io.WriteString(tf, l)
	}
	return
}

func Collect(fn string, chunkSize int, prg *progress.Progress) chan VendorDef {
	chnl := make(chan VendorDef)

	go func() {
		f, err := os.Open(fn)
		MaybePanic(err)
		defer f.Close()
		defer deleteFile(f)
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		lineNo := 1
		chunk := chunkSize
		for scanner.Scan() {
			lineNo++
			line := scanner.Text()
			lp := splitTabs(line)
			om := strings.TrimSpace(lp[0])
			if !strings.Contains(om, "/") {
				om += "/24"
			}
			var p string
			var l int
			var o string
			if len(lp) == 3 {
				_, mp, err := macaddr.ParseMACPrefix(om)
				bm := mp.MAC.String()
				MaybePanic(err)
				l = mp.PrefixLen()
				p = fmt.Sprintf("%s/%d", bm, l)
				o = strings.TrimSpace(lp[2])
				chnl <- VendorDef{Prefix: p, Length: l, Org: o}
			} else if len(lp) == 2 {
				bm, mp, err := macaddr.ParseMACPrefix(om)
				MaybePanic(err)
				l = mp.PrefixLen()
				p = fmt.Sprintf("%s/%d", bm, l)
				o = strings.TrimSpace(lp[1])
				chnl <- VendorDef{Prefix: p, Length: l, Org: o}
			} else if len(lp) > 0 {
				fmt.Fprintf(os.Stderr, "Unable to parse line %d: '%v'\n", lineNo, line)
			}
			if chunk%lineNo == chunk {
				prg.Advance()
				chunk += chunkSize
			}
		}
		close(chnl)
	}()

	return chnl
}
