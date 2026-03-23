// Package cli provides command line usage of migration app.
package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
)

// App is the main CLI app.
type App struct {
	Name        string
	Version     string
	Description string
	commands    map[string]*command
	flags       map[string]*flag
}

// NewApp returns new instance of the CLI app.
func NewApp() *App {
	return &App{
		Name:        "migr8",
		Description: "Database migration tool for PostgreSQL",
		commands:    make(map[string]*command),
		flags:       make(map[string]*flag),
	}
}

// ExecCommand executes command of the CLI app.
func (a *App) ExecCommand(ctx context.Context, parsedArgs *ParseResult) error {
	if err := parsedArgs.Command.Action(ctx, parsedArgs.CommandArgs, parsedArgs.Flags); err != nil {
		return err
	}

	return nil
}

// ShowHelp shows info and usage of the CLI app.
func (a *App) ShowHelp() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintf(w, "  %s [options] <command>\n\n", a.Name)

	fmt.Fprintln(w, "OPTIONS:")

	seen := make(map[string]bool)

	for _, flag := range a.flags {
		if !seen[flag.Name] {
			fmt.Fprintf(w, "  %s\t%s\n", flag.Name, flag.Description)
			seen[flag.Name] = true
		}
	}
	fmt.Fprint(w, "\n")

	fmt.Fprintln(w, "COMMANDS:")
	for _, command := range a.commands {
		fmt.Fprintf(w, "  %s\t%s\n", command.Name, command.Description)
	}

	return w.Flush()
}
