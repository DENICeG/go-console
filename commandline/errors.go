package commandline

import (
	"fmt"
)

var errExit = fmt.Errorf("exit application")
var errCtrlC = fmt.Errorf("Ctrl+C")

type errUnknownCommand struct {
	commandName string
}

func (e errUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %q", e.commandName)
}

// ErrUnknownCommand returns a new error that indicates an unknown command.
func ErrUnknownCommand(commandName string) error {
	return errUnknownCommand{commandName}
}

type errCommandPanicked struct {
	recovered any
}

func (e errCommandPanicked) Error() string {
	return fmt.Sprintf("%v", e.recovered)
}

// ErrCommandPanicked returns a new error that indicates a panicked command.
func ErrCommandPanicked(recovered any) error {
	return errCommandPanicked{recovered}
}
