package input

import (
	"strings"

	"github.com/sbreitf1/go-console"
)

type textEditor struct {
	lines      [][]rune
	caretLine  int
	caretPos   int
	InsertMode bool
}

func newTextEditor(str string) *textEditor {
	strLines := strings.Split(strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n"), "\n")
	lines := make([][]rune, len(strLines))
	for i := 0; i < len(strLines); i++ {
		lines[i] = []rune(strLines[i])
	}

	return &textEditor{InsertMode: false, lines: lines, caretLine: 0, caretPos: 0}
}

func (e *textEditor) Caret() (int, int) {
	caretLine := boundBy(e.caretLine, 0, len(e.lines)-1)
	caretPos := boundBy(e.caretPos, 0, len(e.lines[caretLine]))
	return caretLine, caretPos
}

func (e *textEditor) String() string {
	var sb strings.Builder
	for i := 0; i < len(e.lines); i++ {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(string(e.lines[i]))
	}
	return sb.String()
}

func (e *textEditor) Lines() []string {
	strLines := make([]string, len(e.lines))
	for i := 0; i < len(e.lines); i++ {
		strLines[i] = string(e.lines[i])
	}
	return strLines
}

func (e *textEditor) LineRange(start, count int) []string {
	if start >= len(e.lines) {
		return []string{}
	}

	count = min(len(e.lines)-start, count)
	strLines := make([]string, count)
	for i := start; i < (start + count); i++ {
		strLines[i-start] = string(e.lines[i])
	}
	return strLines
}

func (e *textEditor) MoveCaretLeft() bool {
	// check if caret in bounds. first check here to keep caret pos on short lines during vertical navigation
	if e.caretLine < len(e.lines) && e.caretPos >= len(e.lines[e.caretLine]) {
		e.caretPos = len(e.lines[e.caretLine])
	}

	e.caretPos--
	if e.caretPos < 0 {
		if e.caretLine <= 0 {
			e.caretPos = 0
		} else {
			e.caretLine--
			e.caretPos = len(e.lines[e.caretLine])
		}
	}
	return true
}

func (e *textEditor) MoveCaretRight() bool {
	e.caretPos++
	if e.caretPos >= len(e.lines[e.caretLine]) {
		if e.caretLine >= (len(e.lines) - 1) {
			e.caretPos = len(e.lines[e.caretLine])
		} else {
			e.caretLine++
			e.caretPos = 0
		}
	}
	return true
}

func (e *textEditor) MoveCaretUp(delta int) bool {
	e.caretLine -= delta
	if e.caretLine < 0 {
		e.caretLine = 0
	}

	return true
}

func (e *textEditor) MoveCaretDown(delta int) bool {
	e.caretLine += delta
	if e.caretLine >= len(e.lines) {
		e.caretLine = len(e.lines) - 1
	}

	return true
}

func (e *textEditor) MoveCaretToLineBegin() bool {
	e.caretPos = 0
	return true
}

func (e *textEditor) MoveCaretToLineEnd() bool {
	e.caretPos = len(e.lines[e.caretLine])
	return true
}

func (e *textEditor) InsertAtCaret(str string) {
	caretLine, caretPos := e.Caret()

	str = strings.ReplaceAll(strings.ReplaceAll(str, "\r\n", "\n"), "\r", "\n")
	insertLines := strings.Split(str, "\n")
	if len(insertLines) == 1 {
		// no new line inserted, simple case:
		prefix := string(e.lines[caretLine][:caretPos])
		suffix := string(e.lines[caretLine][caretPos:])
		e.lines[caretLine] = []rune(prefix + insertLines[0] + suffix)

		// move caret to end of inserted string
		e.caretPos = caretPos + len([]rune(insertLines[0]))
		return
	}

	targetCaretPos := len([]rune(insertLines[len(insertLines)-1]))

	// first line first part of old line (to caret) and first line of inserted string
	insertLines[0] = string(e.lines[caretLine][:caretPos]) + insertLines[0]
	// last line is last line of inserted string and last part of old line (behind caret)
	insertLines[len(insertLines)-1] = insertLines[len(insertLines)-1] + string(e.lines[caretLine][caretPos:])
	// insert new lines to slice
	newLines := make([][]rune, len(e.lines)+len(insertLines)-1)
	for i := 0; i < caretLine; i++ {
		newLines[i] = e.lines[i]
	}

	for i := 0; i < len(insertLines); i++ {
		newLines[caretLine+i] = []rune(insertLines[i])
	}

	for i := (caretLine + 1); i < len(e.lines); i++ {
		newLines[len(insertLines)+i-1] = e.lines[i]
	}

	e.lines = newLines

	// move caret to end of inserted string
	e.caretLine = caretLine + len(insertLines) - 1
	e.caretPos = targetCaretPos
}

func (e *textEditor) NewLineAtCaret() {
	e.InsertAtCaret("\n")
}

func (e *textEditor) RemoveLeftOfCaret() bool {
	caretLine, caretPos := e.Caret()
	// fix caret position to currently visible position
	e.caretPos = caretPos

	if caretPos > 0 {
		// removing only in current line
		prefix := string(e.lines[caretLine][:caretPos-1])
		suffix := string(e.lines[caretLine][caretPos:])
		e.lines[caretLine] = []rune(prefix + suffix)
		e.caretPos = caretPos - 1
		return true
	}

	if caretPos == 0 && caretLine > 0 {
		targetCaretPos := len(e.lines[caretLine-1])

		newLines := make([][]rune, len(e.lines)-1)
		for i := 0; i < caretLine; i++ {
			newLines[i] = e.lines[i]
		}

		newLines[caretLine-1] = append(e.lines[caretLine-1], e.lines[caretLine]...) //nolint
		for i := (caretLine + 1); i < len(e.lines); i++ {
			newLines[i-1] = e.lines[i]
		}
		e.lines = newLines

		// move caret to old end of previous line
		e.caretLine = caretLine - 1
		e.caretPos = targetCaretPos
		return true
	}

	return false
}

func (e *textEditor) RemoveRightOfCaret() bool {
	caretLine, caretPos := e.Caret()
	// fix caret position to currently visible position
	e.caretPos = caretPos

	if caretPos < len(e.lines[caretLine]) {
		// removing only in current line
		prefix := string(e.lines[caretLine][:caretPos])
		suffix := string(e.lines[caretLine][caretPos+1:])
		e.lines[caretLine] = []rune(prefix + suffix)
		return true
	}

	if caretPos >= len(e.lines[caretLine]) && caretLine < len(e.lines)-1 {
		newLines := make([][]rune, len(e.lines)-1)
		for i := 0; i <= caretLine; i++ {
			newLines[i] = e.lines[i]
		}

		newLines[caretLine] = append(e.lines[caretLine], e.lines[caretLine+1]...) //nolint
		for i := (caretLine + 2); i < len(e.lines); i++ {
			newLines[i-1] = e.lines[i]
		}
		e.lines = newLines

		return true
	}

	return false
}

func boundBy(val int, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

// Text opens a full screen text editor in console mode and returns the entered string.
//
// Returns true if saved.
func Text(str string) (string, bool, error) {
	screen, err := newScreen()
	if err != nil {
		return "", false, err
	}
	defer screen.Close()

	editor := newTextEditor(str)

	// currently visible rectangle
	firstLine := 0
	firstPos := 0

	for {
		// render current editor view
		screen.Clear()
		viewportWidth, viewportHeight := screen.Size()
		editorOffsetX := 1
		editorOffsetY := 1
		editorWidth := viewportWidth - 2
		editorHeight := viewportHeight - 3

		drawBox(editorOffsetX, editorWidth, screen, editorOffsetY, editorHeight)

		firstLine, firstPos = setCursor(editor, firstLine, editorHeight, firstPos, editorWidth, screen, editorOffsetX, editorOffsetY)

		printText(editor, firstLine, editorHeight, firstPos, editorWidth, screen, editorOffsetX, editorOffsetY)

		printCells(screen, "Esc to exit", 1, editorOffsetY+editorHeight+1)
		printCells(screen, "Strg+S to save", editorOffsetX+editorWidth-14, editorOffsetY+editorHeight+1)

		// display
		screen.Flush()

		switch e := screen.PollEvent().(type) {
		case keyEvent:
			switch e.Key {
			case console.KeyEscape:
				return str, false, nil

			case console.KeyCtrlW:
				// for all nano fans :)
				fallthrough
			case console.KeyCtrlS:
				return editor.String(), true, nil

			case console.KeyLeft:
				editor.MoveCaretLeft()
			case console.KeyRight:
				editor.MoveCaretRight()
			case console.KeyUp:
				editor.MoveCaretUp(1)
			case console.KeyDown:
				editor.MoveCaretDown(1)

			case console.KeyPageUp:
				editor.MoveCaretUp(editorHeight)
			case console.KeyPageDown:
				editor.MoveCaretDown(editorHeight)

			case console.KeyHome:
				editor.MoveCaretToLineBegin()
			case console.KeyEnd:
				editor.MoveCaretToLineEnd()

			case console.KeyBackspace:
				editor.RemoveLeftOfCaret()
			case console.KeyDelete:
				editor.RemoveRightOfCaret()

			case console.KeyEnter:
				editor.NewLineAtCaret()
			case console.KeySpace:
				editor.InsertAtCaret(" ")
			case console.KeyTab:
				editor.InsertAtCaret("    ")
			default:
				if e.Rune != '\000' {
					editor.InsertAtCaret(string(e.Rune))
				}
			}

		case resizeEvent:
			// do nothing, just redraw in next iteration

		case errorEvent:
			return "", false, e.Error
		}
	}
}

func setCursor(editor *textEditor, firstLine int, editorHeight int, firstPos int, editorWidth int, screen screen, editorOffsetX int, editorOffsetY int) (int, int) {
	caretLine, caretPos := editor.Caret()
	// ensure caret is visible
	if caretLine < firstLine {
		firstLine = caretLine
	}

	if caretLine >= (firstLine + editorHeight - 1) {
		firstLine = caretLine - editorHeight + 1
	}

	if caretPos < firstPos {
		firstPos = caretPos
	}

	if caretPos >= (firstPos + editorWidth - 1) {
		firstPos = caretPos - editorWidth + 1
	}

	// set relative caret location
	screen.SetCursor(editorOffsetX+caretPos-firstPos, editorOffsetY+caretLine-firstLine)
	return firstLine, firstPos
}

func printText(editor *textEditor, firstLine int, editorHeight int, firstPos int, editorWidth int, screen screen, editorOffsetX int, editorOffsetY int) {
	currentColor := screen.GetDefaultColor()

	for i, line := range editor.LineRange(firstLine, editorHeight) {
		runes := []rune(line)
		for j := firstPos; j < min(len(runes), firstPos+editorWidth); j++ {
			if IsANSIEscape(line, j) {
				if IsANSIReset(line, j) {
					currentColor = screen.GetDefaultColor()
					j += 3
				} else {
					sequence, end := ReadANSISequence(line, j)
					currentColor = ExtractRGB(sequence)
					j = end
					continue
				}
			}

			screen.SetCellColored(editorOffsetX+j-firstPos, editorOffsetY+i, runes[j], currentColor, screen.GetDefaultColor())
		}
	}
}

func drawBox(editorOffsetX int, editorWidth int, screen screen, editorOffsetY int, editorHeight int) {
	for x := editorOffsetX; x < (editorOffsetX + editorWidth); x++ {
		screen.SetCell(x, editorOffsetY-1, '─')
		screen.SetCell(x, editorOffsetY+editorHeight, '─')
	}

	for y := editorOffsetY; y < (editorOffsetY + editorHeight); y++ {
		screen.SetCell(editorOffsetX-1, y, '│')
		screen.SetCell(editorOffsetX+editorWidth, y, '│')
	}

	screen.SetCell(editorOffsetX-1, editorOffsetY-1, '┌')
	screen.SetCell(editorOffsetX+editorWidth, editorOffsetY-1, '┐')
	screen.SetCell(editorOffsetX-1, editorOffsetY+editorHeight, '└')
	screen.SetCell(editorOffsetX+editorWidth, editorOffsetY+editorHeight, '┘')
}
