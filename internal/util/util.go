package util

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hako/durafmt"
)

func RemoveComments(str string) (c string) {
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

func SplitTabs(i string) []string {
	p := regexp.MustCompile(`\t+`)
	var r []string
	for _, e := range p.Split(i, -1) {
		if e != "" {
			r = append(r, strings.TrimSpace(e))
		}
	}
	return r
}

func PathExists(n string) bool {
	if _, err := os.Stat(n); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func TimeSince(t time.Time) string {
	return durafmt.Parse(time.Since(t)).LimitFirstN(1).String()
}

func SplitSlice[T any](slice []T, max int) [][]T {
	result := make([][]T, 0)

	for i := 0; i < len(slice); i += max {
		end := i + max
		if end > len(slice) {
			end = len(slice)
		}
		part := slice[i:end]
		if part != nil {
			result = append(result, part)
		}
	}
	return result
}
