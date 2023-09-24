package main

import (
	"fmt"
	"os"

	"github.com/thatmattlove/oui/v2/cmd"
)

func getArgs() []string {
	args := os.Args
	if len(args) == 1 {
		args = append(args, "--help")
	}
	return args
}

func main() {
	args := getArgs()
	err := cmd.New(TABLE_VERSION).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
