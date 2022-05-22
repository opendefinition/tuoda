package main

import (
	"fmt"

	"github.com/alecthomas/kong"
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

	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Pc ParseCmd `cmd: help:"Parse log"`
}

func main() {
	fmt.Println("Tuoda Log Importer - v.0.0.1")
	fmt.Println("By Roger Johnsen - Opendefinition 2022")

	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)

}
