package main

import (
	"fmt"
	"os"
)

func init() {
	MaybePanic = func(err error) {
		if err != nil {
			panic(err)
		}
	}
	setLocale()
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
