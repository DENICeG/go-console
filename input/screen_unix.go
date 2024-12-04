//go:build !windows
// +build !windows

package input

import (
	"github.com/sbreitf1/go-console"

	"github.com/nsf/termbox-go"
)

type unixScreen struct{}

func newScreen() (screen, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}

	return &unixScreen{}, nil
}

func (s *unixScreen) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault) //nolint
}
func (s *unixScreen) Size() (int, int) {
	return termbox.Size()
}
func (s *unixScreen) SetCell(x, y int, r rune) {
	termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}

func (s *unixScreen) SetCellColored(x, y int, r rune) {
	termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
}

func (s *unixScreen) Flush() {
	termbox.Flush()
}
func (s *unixScreen) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}
func (s *unixScreen) PollEvent() event {
	// wait for supported event
	for {
		// translate received event
		switch e := termbox.PollEvent(); e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				return keyEvent{console.KeyEscape, '\000'}

			case termbox.KeyCtrlW:
				return keyEvent{console.KeyCtrlW, '\000'}
			case termbox.KeyCtrlS:
				return keyEvent{console.KeyCtrlS, '\000'}

			case termbox.KeyArrowUp:
				return keyEvent{console.KeyUp, '\000'}
			case termbox.KeyArrowDown:
				return keyEvent{console.KeyDown, '\000'}
			case termbox.KeyArrowLeft:
				return keyEvent{console.KeyLeft, '\000'}
			case termbox.KeyArrowRight:
				return keyEvent{console.KeyRight, '\000'}
			case termbox.KeyHome:
				return keyEvent{console.KeyHome, '\000'}
			case termbox.KeyEnd:
				return keyEvent{console.KeyEnd, '\000'}
			case termbox.KeyPgup:
				return keyEvent{console.KeyPageUp, '\000'}
			case termbox.KeyPgdn:
				return keyEvent{console.KeyPageDown, '\000'}

			case termbox.KeyBackspace:
				fallthrough
			case termbox.KeyBackspace2:
				return keyEvent{console.KeyBackspace, '\r'}
			case termbox.KeyDelete:
				return keyEvent{console.KeyDelete, '\000'}
			case termbox.KeyEnter:
				return keyEvent{console.KeyEnter, '\n'}
			case termbox.KeySpace:
				return keyEvent{console.KeySpace, ' '}
			case termbox.KeyTab:
				return keyEvent{console.KeyTab, '\t'}

			default:
				return keyEvent{0, e.Ch}
			}

		case termbox.EventResize:
			return resizeEvent{}

		case termbox.EventError:
			return errorEvent{e.Err}
		}
	}
}
func (s *unixScreen) Close() {
	termbox.Close()
}
