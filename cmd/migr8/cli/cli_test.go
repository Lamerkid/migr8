package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_SimpleCommand(t *testing.T) {
	args := make([]string, 3)
	args[0] = "test"
	args[1] = "-t"
	args[2] = "testFlag"
	var executed bool

	app := NewApp()

	app.addCommand(&command{
		Name: "test",
		Action: func(_ context.Context, _ []string, _ map[string]string) error {
			executed = true
			return nil
		},
	})

	app.addFlag(&flag{
		Name:        "-t",
		Description: "test flag",
	})

	parsed, err := app.ParseArgs(args)
	if err != nil {
		t.Fatalf("failed to parse arguments: %v", err)
	}

	require.Equal(t, "test", parsed.Command.Name)
	require.Equal(t, "testFlag", parsed.Flags["-t"])

	err = app.ExecCommand(context.Background(), parsed)
	if err != nil {
		t.Fatalf("error executing command: %v", err)
	}

	if !executed {
		t.Error("command was not executed")
	}
}
