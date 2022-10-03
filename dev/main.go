// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

	"git.tech.kora.id/go/dev/core"
	"github.com/gemalto/flume"
)

func init() {
	cfg := flume.Config{
		DefaultLevel: flume.InfoLevel,
		Development:  true,
		Encoding:     "term-color",
	}

	flume.Configure(cfg)
	core.Log = flume.New("dev-cli")
}

// list of available command.
var cmd = []*core.Command{
	runCommand,
	makeCommand,
}

func main() {
	flag.Usage = showUsage
	flag.Parse()
	log.SetFlags(0)

	// if no argument passes, show usage,
	args := flag.Args()
	if len(args) < 1 {
		showUsage()
	}

	// show help when argument is help
	if args[0] == "help" {
		showHelp(args[1:])
		return
	}

	// check the argument is registered as command.
	for _, c := range cmd {
		if c.Name == args[0] && c.IsRunable() {
			c.Flag.Usage = func() {
				c.ShowUsage()
			}

			if c.CustomFlags {
				args = args[1:]
			} else {
				c.Flag.Parse(args[1:])
				args = c.Flag.Args()
			}

			os.Exit(c.Run(c, args))
			return
		}
	}

	core.Log.Error(fmt.Sprintf("Unknown command %q.", args[0]))
	core.Log.Error("Run 'dev help' to see list of command")
	os.Exit(2)
}

// showUsage print usage message into cli.
func showUsage() {
	var usageTexts = "dev (tools) is cli tools for go projects.\n\nUsage:\n\tdev [command] [arguments]\n\nThe command available are:\n {{range .}}\n\t{{.Name | printf \"%-11s\"}} {{.Info}}{{end}}\n\nUse 'dev help [command]' for more information about a command.\n"

	compileTemplate(os.Stdout, usageTexts, cmd)
	os.Exit(2)
}

// showHelp print help message into cli.
func showHelp(args []string) {
	if len(args) == 0 {
		showUsage()
		return
	}

	if len(args) != 1 {
		core.Log.Error("Too many arguments given.")
		core.Log.Error("Usage: \n\tdev help command.")
		os.Exit(2)
	}

	arg := args[0]
	for _, c := range cmd {
		if c.Name == arg {
			helpText := "Usage:\n\tdev {{.Name}}\n\n{{.Info}}\n\n{{.Usage | trim}}"
			compileTemplate(os.Stdout, helpText, c)
			return
		}
	}

	core.Log.Error(fmt.Sprintf("Unknown help topic %#q.", arg))
	core.Log.Error("Run 'dev help' to see list of topic")
	os.Exit(2)
}

// compileTemplate, compiling html/template into writer with data command.
func compileTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
