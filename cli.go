package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/camelcase"
)

const (
	ExitCodeOk = iota
	ExitCodeNoArg
	ExitCodeFileOpenError
)

var helpText = `
Usage: ghead [file ...]
`

var replaceCharRegex = regexp.MustCompile("[^a-zA-Z]")
var replaceSpaceRegex = regexp.MustCompile(" +")

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	if len(args) == 1 {
		fmt.Print(helpText)
		return ExitCodeNoArg
	}

	for _, filename := range args[1:] {
		file, oErr := os.Open(filename)
		if oErr != nil {
			return ExitCodeFileOpenError
		}
		defer file.Close()
		tokens := tokenize(file)
		cErr := check(tokens)
		if cErr != nil {
			return ExitCodeFileOpenError
		}
	}

	return ExitCodeOk
}

func tokenize(file io.Reader) []string {
	tokens := make([]string, 200, 1000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		filtered := replaceCharRegex.ReplaceAllString(line, " ")
		filtered = replaceSpaceRegex.ReplaceAllString(filtered, " ")
		filtered = strings.Trim(filtered, " ")
		if len(filtered) > 2 {
			tokens = append(tokens, filtered)
		}
	}
	return tokens
}

func check(tokens []string) error {
	for _, token := range tokens {
		for _, word := range strings.Split(token, " ") {
			for _, w := range camelcase.Split(word) {
				if len(w) > 2 {
					fmt.Println(w)
				}
			}
		}
	}

	return nil
}
