package commandline

import (
	"fmt"
)

// ErrExit is returned when the user wants to exit the application.
var ErrExit = fmt.Errorf("exit application")

// ErrCtrlC is returned when the user wants to stop the application.
var ErrCtrlC = fmt.Errorf("Ctrl+C")

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

type ErrCommandPanicked struct {
	recovered any
}

func (e ErrCommandPanicked) Error() string {
	return fmt.Sprintf("%v", e.recovered)
}

// NewErrCommandPanicked returns a new error that indicates a panicked command.
func NewErrCommandPanicked(recovered any) error {
	return ErrCommandPanicked{recovered}
}
