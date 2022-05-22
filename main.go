package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/opendefinition/tuoda/parsers"
)

type Context struct {
	Debug bool
}

type ParseCmd struct {
	Parser  string `help:"Path to parser definition"`
	LogFile string `help:"Path to log file or"`
}

func (pc *ParseCmd) Run(ctx *Context) error {
	fmt.Println(pc.Parser)
	fmt.Println(pc.LogFile)

	test := new(parsers.CsvDefinition)
	doc, _ := json.Marshal(test)
	fmt.Println(string(doc))

	fmt.Println("Testings")
	parserdefRaw := ReadParserDefinition(pc.Parser)
	fmt.Println(parserdefRaw)

	// Finding parser type
	regx := regexp.MustCompile("\"parser_type\":\\s\"(\\w+)\"")
	match := regx.FindStringSubmatch(parserdefRaw)

	if len(match) == 0 {
		fmt.Println("No match for parser, sorry")
	} else {
		// Pretend there's a switch here
		obj := new(parsers.CsvDefinition)
		json.Unmarshal([]byte(parserdefRaw), obj)

		obj.Parse()
	}

	return nil
}

func ReadParserDefinition(filepath string) string {
	definitionFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}

	buffer := new(strings.Builder)
	io.Copy(buffer, definitionFile)

	definitionFile.Close()

	return buffer.String()
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Parse ParseCmd `cmd: help:"Parse log"`
}

func main() {
	fmt.Println("Tuoda Log Importer - v.0.0.1")
	fmt.Println("By Roger Johnsen - Opendefinition 2022\n")

	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
