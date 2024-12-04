package commandline

// Command denotes a named command with completion and execution handler.
type Command interface {
	// Name returns the name of the command as used in the command line.
	Name() string
	// GetCompletionOptions denotes a custom completion handler as used for ReadCommand.
	GetCompletionOptions(currentCommand []string, entryIndex int) []CompletionOption
	// Exec is called to execute the command with a set of arguments.
	Exec(args []string) error
}

// ExecCommandHandler is called when processing a command. Return ErrExit to gracefully stop processing.
type ExecCommandHandler func(args []string) error

type customCommand struct {
	completionHandler CommandCompletionHandler
	execHandler       ExecCommandHandler
	name              string
}

func (c *customCommand) Name() string {
	return c.name
}
func (c *customCommand) GetCompletionOptions(currentCommand []string, entryIndex int) []CompletionOption {
	if c.completionHandler != nil {
		return c.completionHandler(currentCommand, entryIndex)
	}
	return nil
}
func (c *customCommand) Exec(args []string) error {
	if c.execHandler != nil {
		return c.execHandler(args)
	}
	return nil
}

// NewExitCommand returns a named command to stop command line processing.
func NewExitCommand(name string) Command {
	return &customCommand{
		name:              name,
		completionHandler: func([]string, int) []CompletionOption { return nil },
		execHandler:       func([]string) error { return errExit },
	}
}

// NewParameterlessCommand returns a named command that takes no parameters.
func NewParameterlessCommand(name string, handler ExecCommandHandler) Command {
	return &customCommand{
		name:              name,
		completionHandler: func([]string, int) []CompletionOption { return nil },
		execHandler:       handler,
	}
}

// NewCustomCommand returns a named command with completion and execution handler.
func NewCustomCommand(name string, completionHandler CommandCompletionHandler, execHandler ExecCommandHandler) Command {
	return &customCommand{
		name:              name,
		completionHandler: completionHandler,
		execHandler:       execHandler,
	}
}
