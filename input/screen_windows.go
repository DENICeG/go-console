//go:build windows
// +build windows

package input

import (
	"github.com/DENICeG/go-console/v2"

	"github.com/gdamore/tcell"
)

type windowsScreen struct {
	screen tcell.Screen
}

func newScreen() (screen, error) {
	screen, err := tcell.NewConsoleScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	return &windowsScreen{screen}, nil
}
func (s *windowsScreen) GetDefaultColor() RGB {
	fg, _, _ := tcell.StyleDefault.Decompose()
	r, g, b := fg.RGB()
	return RGB{R: uint8(r), G: uint8(g), B: uint8(b)}
}

func (s *windowsScreen) Clear() {
	s.screen.Clear()
}

func (s *windowsScreen) Size() (int, int) {
	return s.screen.Size()
}

func (s *windowsScreen) SetCell(x, y int, r rune) {
	s.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
}

func (s *windowsScreen) SetCellColored(x, y int, r rune, fg, bg RGB) {
	style := tcell.Style(0)

	style = style.Background(tcell.NewRGBColor(int32(bg.R), int32(bg.G), int32(bg.B)))
	style = style.Foreground(tcell.NewRGBColor(int32(fg.R), int32(fg.G), int32(fg.B)))

	s.screen.SetContent(x, y, r, nil, style)
}

func (s *windowsScreen) Flush() {
	s.screen.Sync()
}

func (s *windowsScreen) SetCursor(x, y int) {
	s.screen.ShowCursor(x, y)
}

func (s *windowsScreen) PollEvent() event {
	// wait for supported event
	for {
		// translate received event
		switch e := s.screen.PollEvent().(type) {
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape:
				return keyEvent{console.KeyEscape, '\000'}

			case tcell.KeyCtrlW:
				return keyEvent{console.KeyCtrlW, '\000'}
			case tcell.KeyCtrlS:
				return keyEvent{console.KeyCtrlS, '\000'}

			case tcell.KeyUp:
				return keyEvent{console.KeyUp, '\000'}
			case tcell.KeyDown:
				return keyEvent{console.KeyDown, '\000'}
			case tcell.KeyLeft:
				return keyEvent{console.KeyLeft, '\000'}
			case tcell.KeyRight:
				return keyEvent{console.KeyRight, '\000'}
			case tcell.KeyHome:
				return keyEvent{console.KeyHome, '\000'}
			case tcell.KeyEnd:
				return keyEvent{console.KeyEnd, '\000'}
			case tcell.KeyPgUp:
				return keyEvent{console.KeyPageUp, '\000'}
			case tcell.KeyPgDn:
				return keyEvent{console.KeyPageDown, '\000'}

			case tcell.KeyBackspace:
				fallthrough
			case tcell.KeyBackspace2:
				return keyEvent{console.KeyBackspace, '\r'}
			case tcell.KeyDelete:
				return keyEvent{console.KeyDelete, '\000'}
			case tcell.KeyEnter:
				return keyEvent{console.KeyEnter, '\n'}
			case tcell.KeyTab:
				return keyEvent{console.KeyTab, '\t'}

			default:
				if e.Rune() == ' ' {
					return keyEvent{console.KeySpace, ' '}
				}
				return keyEvent{0, e.Rune()}
			}

		case *tcell.EventResize:
			return resizeEvent{}

		case *tcell.EventError:
			return errorEvent{e}
		}
	}
}
func (s *windowsScreen) Close() {
	s.screen.Fini() //nolint
}
