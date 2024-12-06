package input

import "github.com/sbreitf1/go-console"

// RGB is a RGB color representation.
type RGB struct {
	R uint8
	G uint8
	B uint8
}

type screen interface {
	Clear()
	Size() (int, int)
	SetCell(x, y int, r rune)
	SetCellColored(x, y int, r rune, foreground, background RGB)
	GetDefaultColor() RGB
	Flush()
	SetCursor(x, y int)
	PollEvent() event
	Close()
}

type event any

type errorEvent struct {
	Error error
}

type keyEvent struct {
	Key  console.Key
	Rune rune
}

type resizeEvent struct{}

func printCells(screen screen, str string, x, y int) {
	//TODO support for multiline string
	runes := []rune(str)
	for i := range runes {
		screen.SetCell(x+i, y, runes[i])
	}
}
