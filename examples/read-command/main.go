package main

import (
	"errors"

	"github.com/DENICeG/go-console/v2"
	"github.com/DENICeG/go-console/v2/commandline"
)

func main() {
	history := commandline.NewCommandHistory(3)

	opts := &commandline.ReadCommandOptions{
		GetHistoryEntry:     history.GetHistoryEntry,
		PrintOptionsHandler: commandline.DefaultOptionsPrinter(),
	}

	console.Println("type exit to leave")
	for {
		cmd, err := commandline.ReadCommand("command", opts)
		if err != nil {
			if errors.Is(err, commandline.ErrCtrlC) {
				console.Println()
				break
			}

			console.Fatallnf("ReadCommand failed: %s", err.Error())
		}

		if len(cmd) > 0 {
			console.Printlnf("# %q", cmd[0])
			for i := 1; i < len(cmd); i++ {
				console.Printlnf("-> %q", cmd[i])
			}

			history.Put(cmd)

			if cmd[0] == "exit" {
				break
			}
		} else {
			// empty command
		}
	}
}
