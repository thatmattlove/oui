package main

import (
	"fmt"
	"os"

	"github.com/Xuanwo/go-locale"
	"golang.org/x/text/language"
)

func init() {
	tag, err := locale.Detect()
	if err != nil {
		tag = language.English
	}
	_locale = tag
}

func getArgs() []string {
	args := os.Args
	if len(args) == 1 {
		args = append(args, "--help")
	}
	return args
}

func main() {
	args := getArgs()
	err := CLI().Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
