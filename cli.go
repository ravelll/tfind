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
	ExitCodeTokenizeError
)

var helpText = `
Usage: ghead [file ...]
`

var replaceRegex = regexp.MustCompile("[^a-zA-Z]")
var dict = make([]string, 0, 40000)

type CLI struct {
	outStream, errStream io.Writer
}

func init() {
	f, err := os.Open("assets/dict.txt")
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		dict = append(dict, s.Text())
	}
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

		tokens, tErr := tokenize(file)
		if tErr != nil {
			return ExitCodeTokenizeError
		}

		typos := make([]string, 10, 10)
		for _, token := range tokens {
			typos = append(typos, check(token))
		}
	}

	return ExitCodeOk
}

func tokenize(file io.Reader) ([]string, error) {
	tokens := make([]string, 200, 1000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		filtered := replaceRegex.ReplaceAllString(line, " ")
		for _, symbol := range strings.Split(filtered, " ") {
			for _, word := range camelcase.Split(symbol) {
				if len(word) >= 3 {
					tokens = append(tokens, strings.ToLower(word))
				}
			}
		}
	}

	return tokens, nil
}

func check(token string) string {
	for _, w := range dict {
		if w == token {
			return ""
		}
	}

	fmt.Println(token)

	return ""
}
