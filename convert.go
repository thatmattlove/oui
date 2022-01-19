package main

import (
	"strings"

	"github.com/thatmattlove/go-macaddr"
)

type Formats struct {
	Hex    string
	Dotted string
	Dashed string
	Int    int64
	Bytes  string
}

func Convert(i string) (fmts *Formats) {
	i = strings.Split(i, "/")[0]
	mac, err := macaddr.ParseMACAddress(i)
	MaybePanic(err)
	fmts = &Formats{
		Hex:    mac.String(),
		Dotted: mac.Dots(),
		Dashed: mac.Dashes(),
		Int:    mac.Int(),
		Bytes:  mac.ByteString(),
	}
	return
}
