package main

import (
	"github.com/DENICeG/go-console/v2"
	"github.com/DENICeG/go-console/v2/commandline"
)

func main() {
	console.Println("Enter an empty line to exit")

	history := commandline.NewLineHistory(5)

	for {
		console.Print("enter> ")
		line, err := commandline.ReadLineWithHistory(history)
		if err != nil {
			console.Fatallnf("ReadLineWithHistory failed: %s", err.Error())
		}

		if len(line) == 0 {
			break
		}

		history.Put(line)
		console.Printlnf("-> %q", line)
	}
}
