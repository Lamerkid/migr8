package cli

import (
	"fmt"
	"strings"
)

// ParseResult is parsed commands and flags.
type ParseResult struct {
	Command     *command
	CommandArgs []string
	Flags       map[string]string
}

// ParseArgs parses all arguments passed to the CLI app.
func (a *App) ParseArgs(args []string) (*ParseResult, error) {
	result := &ParseResult{
		Flags: make(map[string]string),
	}

	var currentFlag *flag
	var currentCommand *command
	var commandArgs []string

	for _, arg := range args {
		// Handle flags
		if strings.HasPrefix(arg, "-") {
			if currentFlag != nil {
				result.Flags[currentFlag.Name] = ""
				currentFlag = nil
			}

			flag, flagValue, err := a.parseFlag(arg)
			if err != nil {
				return nil, err
			}

			if flagValue != "" {
				result.Flags[flag.Name] = flagValue
			} else {
				currentFlag = flag
			}
			continue
		}

		if currentFlag != nil {
			result.Flags[currentFlag.Name] = arg
			currentFlag = nil
			continue
		}

		// Handle commands
		if currentCommand == nil {
			cmd, exists := a.commands[arg]
			if !exists {
				return nil, fmt.Errorf("unknown command: %s", arg)
			}

			currentCommand = cmd
			result.Command = currentCommand
			result.CommandArgs = commandArgs
			continue
		}

		commandArgs = append(commandArgs, arg)
		result.Command = currentCommand
		result.CommandArgs = commandArgs
	}

	// Handle last flag with no value
	if currentFlag != nil {
		result.Flags[currentFlag.Name] = ""
	}

	return result, nil
}

func (a *App) parseFlag(arg string) (*flag, string, error) {
	if strings.Contains(arg, "=") {
		parts := strings.SplitN(arg, "=", 2)
		flagName := parts[0]
		flagValue := parts[1]

		flag, exists := a.flags[flagName]
		if !exists {
			return nil, "", fmt.Errorf("unknown option: %s", flagName)
		}

		return flag, flagValue, nil
	}

	flag, exists := a.flags[arg]
	if !exists {
		return nil, "", fmt.Errorf("unknown option: %s", arg)
	}

	return flag, "", nil
}
