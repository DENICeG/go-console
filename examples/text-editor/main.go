package main

import (
	"github.com/DENICeG/go-console/v2"
	"github.com/DENICeG/go-console/v2/input"
)

func main() {
	str, ok, err := input.Text("default string")
	if err != nil {
		console.Fatallnf("FATAL: %s", err.Error())
	}

	if ok {
		console.Printlnf("You entered %q", str)
	} else {
		console.Println("Abort")
	}
}
