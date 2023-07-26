package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/text"
	"golang.org/x/text/message"
)

type Logger struct {
	Accent *text.Colors
	Writer *message.Printer
}

func NewLogger() *Logger {
	logger := &Logger{
		Accent: &text.Colors{text.Bold, text.Underline},
		Writer: message.NewPrinter(_locale),
	}
	return logger
}

func (l *Logger) WithAccent(color text.Color) *text.Colors {
	a := &text.Colors{color}
	*a = append(*a, *l.Accent...)
	return a
}

func (l *Logger) Success(s string, f ...interface{}) {
	t := &text.Colors{text.FgHiGreen}
	var v string
	if len(f) == 0 {
		v = t.Sprint(s)
	} else {
		v = t.Sprintf(s, l.Accent.Sprint(l.Writer.Sprint(f...)))
	}
	l.Writer.Fprintln(os.Stdout, v)
}

func (l *Logger) Info(s string, f ...interface{}) {
	t := &text.Colors{text.FgHiCyan}
	var v string
	if len(f) == 0 {
		v = t.Sprint(s)
	} else {
		v = t.Sprintf(s, l.Accent.Sprint(l.Writer.Sprint(f...)))
	}
	l.Writer.Fprintln(os.Stdout, v)
}

func (l *Logger) Warn(s string, f ...interface{}) {
	t := &text.Colors{text.FgHiYellow}

	var v string
	if len(f) == 0 {
		v = t.Sprint(s)
	} else {
		v = t.Sprintf(s, l.Accent.Sprint(l.Writer.Sprint(f...)))
	}
	l.Writer.Fprintln(os.Stdout, v)
}

func (l *Logger) Error(s string, f ...interface{}) {
	t := &text.Colors{text.FgHiRed}
	var v string
	if len(f) == 0 {
		v = t.Sprint(s)
	} else {
		v = t.Sprintf(s, l.Accent.Sprint(l.Writer.Sprint(f...)))
	}
	l.Writer.Fprintln(os.Stderr, v)
}

func (l *Logger) Err(err error, strs ...string) {
	t := &text.Colors{text.FgHiRed}
	e := err.Error()
	var v string
	if len(strs) == 0 {
		v = t.Sprint(e)
	} else if len(strs) == 1 {
		v = t.Sprintf("%s: %s", strs[0], e)
	} else {
		s := strs[0]
		f := strs[1:]
		fi := []interface{}{}
		for _, i := range f {
			fi = append(fi, l.Writer.Sprint(i))
		}
		wa := l.WithAccent(text.FgHiRed)
		ss := wa.Sprintf(s, fi...)
		v = fmt.Sprintf("%s: %s", ss, t.Sprint(e))
	}
	l.Writer.Fprintln(os.Stderr, v)
}
