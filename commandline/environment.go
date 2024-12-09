package commandline

import (
	"errors"

	"github.com/DENICeG/go-console/v2"
)

// Environment represents a command line interface environment with history and auto-completion.
type Environment struct {
	history                  CommandHistory
	Prompt                   PromptHandler
	PrintOptions             PrintOptionsHandler
	ExecUnknownCommand       ExecUnknownCommandHandler
	CompleteUnknownCommand   CommandCompletionHandler
	ErrorHandler             CommandErrorHandler
	commands                 map[string]Command
	RecoverPanickedCommands  bool
	UseCommandNameCompletion bool
}

// NewEnvironment returns a new command line environment.
func NewEnvironment() *Environment {
	return &Environment{
		Prompt:       func() string { return "cle" },
		PrintOptions: DefaultOptionsPrinter(),
		ExecUnknownCommand: func(cmd string, _ []string) error {
			_, err := console.Printlnf("Unknown command %q", cmd)
			return err
		},
		CompleteUnknownCommand: nil,
		ErrorHandler: func(_ string, _ []string, err error) error {
			if errors.Is(err, ErrCommandPanicked{}) {
				console.Printlnf("PANIC: %s", err.Error()) //nolint
			} else {
				console.Printlnf("ERROR: %s", err.Error()) //nolint
			}
			return nil
		},
		RecoverPanickedCommands:  true,
		UseCommandNameCompletion: true,
		history:                  NewCommandHistory(100),
		commands:                 make(map[string]Command),
	}
}

// SetStaticPrompt sets a constant prompt to display for command input.
func (b *Environment) SetStaticPrompt(prompt string) {
	b.Prompt = func() string { return prompt }
}

func (b *Environment) prompt() string {
	if b.Prompt == nil {
		return ""
	}
	return b.Prompt()
}

// RegisterCommand adds a new command to the command line environment.
func (b *Environment) RegisterCommand(cmd Command) {
	b.commands[cmd.Name()] = cmd
}

// UnregisterCommand removes a command from the command line environment and returns true if it was existent before.
func (b *Environment) UnregisterCommand(commandName string) bool {
	_, exists := b.commands[commandName]
	if exists {
		delete(b.commands, commandName)
	}
	return exists
}

// ReadCommand reads a command for the configured environment.
func (b *Environment) ReadCommand() ([]string, error) {
	return b.readCommand(ReadCommand)
}

func (b *Environment) readCommand(handler func(prompt string, opts *ReadCommandOptions) ([]string, error)) ([]string, error) {
	opts := &ReadCommandOptions{
		GetHistoryEntry:      b.history.GetHistoryEntry,
		GetCompletionOptions: b.GetCompletionOptions,
		PrintOptionsHandler:  b.PrintOptions,
	}
	cmd, err := handler(b.prompt(), opts)
	if err != nil {
		return nil, err
	}

	if len(cmd) > 0 && len(cmd[0]) > 0 {
		b.history.Put(cmd)
	}
	return cmd, nil
}

// Run reads and processes commands until an error is returned. Use ErrExit to gracefully stop processing.
func (b *Environment) Run() error {
	for {
		cmd, err := b.readCommand(ReadCommand)
		if err != nil {
			return err
		}

		if len(cmd) > 0 {
			if err := b.ExecCommand(cmd[0], cmd[1:]); err != nil {
				if errors.Is(err, ErrExit) {
					return nil
				}
				if b.ErrorHandler == nil {
					return err
				}

				b.ErrorHandler(cmd[0], cmd[1:], err)
			}
		}
	}
}

// GetCompletionOptions returns completion options for the given command. This method can be used as callback for ReadCommand.
func (b *Environment) GetCompletionOptions(currentCommand []string, entryIndex int) []CompletionOption {
	if entryIndex == 0 {
		if b.UseCommandNameCompletion {
			// completion for command
			options := make([]CompletionOption, 0)
			for name := range b.commands {
				options = append(options, &completionOption{replacement: name})
			}
			return options
		}
		return nil
	}

	cmd, exists := b.commands[currentCommand[0]]
	if !exists {
		if b.CompleteUnknownCommand != nil {
			return b.CompleteUnknownCommand(currentCommand, entryIndex)
		}
		return nil
	}

	return cmd.GetCompletionOptions(currentCommand, entryIndex)
}

// ExecCommand executes a command as if it has been entered in terminal.
func (b *Environment) ExecCommand(cmd string, args []string) error {
	var recovered any

	err := func() error {
		defer func() {
			if b.RecoverPanickedCommands {
				// recover from panic and save reason for error handling
				recovered = recover()
			}
		}()

		// execute command
		if c, exists := b.commands[cmd]; exists {
			return c.Exec(args)
		}
		if b.ExecUnknownCommand == nil {
			return ErrUnknownCommand(cmd)
		}
		return b.ExecUnknownCommand(cmd, args)
	}()

	if recovered != nil {
		return NewErrCommandPanicked(recovered)
	}
	return err
}
