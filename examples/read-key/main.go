package main

import (
	"github.com/sbreitf1/go-console"
)

func main() {
	console.Println("Press Escape to exit")
	for {
		key, char, err := console.ReadKey()
		if err != nil {
			panic(err)
		}
		console.Printlnf("%s -> %q", key, string(char))

		if key == console.KeyEscape {
			break
		}
	}
}