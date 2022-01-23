package main

import (
	"errors"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Xuanwo/go-locale"
	"github.com/hako/durafmt"
	"golang.org/x/text/language"
)

func removeComments(str string) (c string) {
	has := regexp.MustCompile(`^[^\\]*#.*$`)
	re := regexp.MustCompile(`^([^\\#]+)|(#.*)$`)
	if has.MatchString(str) {
		p := re.FindStringSubmatch(str)
		c = p[1]
	} else {
		c = str
	}
	return strings.TrimSpace(c)
}

func splitTabs(i string) []string {
	p := regexp.MustCompile(`\t+`)
	var r []string
	for _, e := range p.Split(i, -1) {
		if e != "" {
			r = append(r, strings.TrimSpace(e))
		}
	}
	return r
}

func pathExists(n string) bool {
	if _, err := os.Stat(n); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func deleteFile(f *os.File) (err error) {
	n := f.Name()
	if pathExists(n) {
		err = os.Remove(n)
	}
	return
}

func timeSince(t time.Time) string {
	return durafmt.Parse(time.Since(t)).LimitFirstN(1).String()
}

func setLocale() {
	tag, err := locale.Detect()
	if err != nil {
		tag = language.English
	}
	_locale = tag
}

func containsStr(arr []string, search string) bool {
	i := sort.SearchStrings(arr, search)
	return i < len(arr) && arr[i] == search
}
