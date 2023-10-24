package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/thatmattlove/oui/v2/cmd"
)

func isPiped() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

var splitPattern = regexp.MustCompile(`[\n\r\t\s]`)

func main() {
	args := os.Args
	if isPiped() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			parts := splitPattern.Split(line, -1)
			for _, part := range parts {
				clean := splitPattern.ReplaceAllString(part, "")
				args = append(args, clean)
			}
		}
	}
	if len(args) == 1 {
		args = append(args, "--help")
	}
	err := cmd.New(Version).Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
