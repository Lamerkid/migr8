// Package main is an entrypoint for the migration CLI app.
package main

import (
	"context"
	"log"
	"os"

	"github.com/Lamerkid/migr8/cmd/migr8/cli"
)

var version string

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := cli.NewApp()
	app.Version = version

	cli.RegisterFlags(app)
	cli.RegisterCommands(app)

	if len(os.Args) < 2 {
		_ = app.ShowHelp()
		return 0
	}

	parsedArgs, err := app.ParseArgs(os.Args[1:])
	if err != nil {
		log.Printf("error parsing arguments:\n  %v", err)
		return 1
	}

	if parsedArgs.Command == nil {
		_ = app.ShowHelp()
		return 2
	}

	if err := app.ExecCommand(ctx, parsedArgs); err != nil {
		log.Printf("error executing command:\n  %v", err)
		return 3
	}

	return 0
}
