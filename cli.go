package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var helpText = `
Usage: ghead [file ...]
`

var replacer = strings.NewReplacer(
	"@", " ",
	"_", " ",
	",", " ",
	".", " ",
	"-", " ",
	"+", " ",
	"=", " ",
	"&", " ",
	"$", " ",
	"*", " ",
	"^", " ",
	"#", " ",
	"!", " ",
	"?", " ",
	"/", " ",
	"\\", " ",
	"|", " ",
	"(", " ",
	")", " ",
	"[", " ",
	"]", " ",
	"{", " ",
	"}", " ",
	"<", " ",
	">", " ",
	"~", " ",
	"`", " ",
	";", " ",
	":", " ",
	"\"", " ",
	"'", " ",
)

const (
	ExitCodeOk = iota
	ExitCodeNoArg
	ExitCodeFileOpenError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	if len(args) == 0 {
		fmt.Print(helpText)
		return ExitCodeNoArg
	}

	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			return ExitCodeFileOpenError
		}

		tokens := tokenize(file)
		tokens[0] = "a"
	}

	return ExitCodeOk
}

func tokenize(file io.Reader) []string {
	tokens := make([]string, 200, 1000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		filtered := replacer.Replace(line)
		tokens = append(tokens, filtered)
		fmt.Println(filtered)
	}

	return tokens
}
