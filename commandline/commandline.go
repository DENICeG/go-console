package commandline

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sbreitf1/go-console"
)

const (
	maxAutoPrintListLen = 100
)

var (
	// remember the last time Tab was pressed to detect double-tab.
	lastTabPress = time.Unix(0, 0)
)

const doubleTabSpan = 250 * time.Millisecond

// CommandHistoryHandler describes a function that returns a command from history at the given index.
//
// Index 0 denotes the latest command. nil is returned when the number of entries in history is exceeded. The index will never be negative.
type CommandHistoryHandler func(index int) ([]string, bool)

// ReadCommandOptions configures options and callbacks for ReadComman.
type ReadCommandOptions struct {
	// GetHistoryEntry denotes the handler for reading command history.
	GetHistoryEntry CommandHistoryHandler
	// GetCompletionOptions denotes the handler for auto completion.
	GetCompletionOptions CommandCompletionHandler
	// PrintOptionsHandler denotes the handler to print options on double-tab.
	PrintOptionsHandler PrintOptionsHandler
}

// ReadCommand reads a command from console input and offers history, aswell as completion functionality.
func ReadCommand(prompt string, opts *ReadCommandOptions) ([]string, error) {
	if opts == nil {
		opts = &ReadCommandOptions{
			PrintOptionsHandler: DefaultOptionsPrinter(),
		}
	}

	var cmd []string
	err := console.WithReadKeyContext(func() error {
		var err error
		cmd, err = readCommand(prompt, opts)
		return err
	})
	return cmd, err
}

func readCommand(prompt string, opts *ReadCommandOptions) ([]string, error) {
	var sb strings.Builder

	for {
		line, err := readCommandLine(&prompt, sb.String(), true, opts)
		if err != nil {
			return nil, err
		}

		sb.WriteString(line)

		if cmd, isComplete := ParseCommand(sb.String()); isComplete {
			return cmd, nil
		}

		// line break is part of command -> append to command because it has been omitted by the line reader
		sb.WriteRune('\n')
		// show empty prompt for further lines of same command
		prompt = ""
	}
}

func readCommandLine(prompt *string, currentCommand string, escapeHistory bool, opts *ReadCommandOptions) (string, error) {
	if prompt != nil {
		console.Printf("%s> ", *prompt) //nolint
	}

	var cmdToString func([]string) string
	if escapeHistory {
		cmdToString = GetCommandString
	} else {
		cmdToString = func(cmd []string) string { return strings.Join(cmd, " ") }
	}

	var sb strings.Builder

	lineLen := 0

	putRune := func(r rune) {
		sb.WriteRune(r)
		console.Print(string(r)) //nolint
		lineLen++
	}

	putString := func(str string) {
		sb.WriteString(str)
		console.Print(str) //nolint
		lineLen += len(str)
	}

	clearLine := func() {
		sb.Reset()
		str1 := strings.Repeat("\b", lineLen)
		str2 := strings.Repeat(" ", lineLen)
		console.Printf("%s%s%s", str1, str2, str1) //nolint
		lineLen = 0
	}

	replaceLine := func(newLine string) {
		clearLine()
		putString(newLine)
	}

	reprintLine := func() {
		if prompt == nil {
			console.Printf("%s", sb.String()) //nolint
		} else {
			console.Printf("%s> %s", *prompt, sb.String()) //nolint
		}
	}

	removeLastChar := func() {
		if lineLen > 0 {
			str := []rune(sb.String())
			sb.Reset()
			if len(str) > 0 {
				sb.WriteString(string(str[:len(str)-1]))
			}

			console.Print("\b \b") //nolint
			lineLen--
		}
	}

	historyIndex := -1

	for {
		key, r, err := console.ReadKey()
		if err != nil {
			return "", err
		}

		switch key {
		case console.KeyCtrlC:
			return "", errCtrlC

		case console.KeyEscape:
			clearLine()

		case console.KeyUp:
			if opts.GetHistoryEntry != nil {
				if newCmd, ok := opts.GetHistoryEntry(historyIndex + 1); ok {
					historyIndex++
					replaceLine(cmdToString(newCmd))
				}
			}
		case console.KeyDown:
			if opts.GetHistoryEntry != nil {
				if historyIndex >= 0 {
					historyIndex--

					if historyIndex >= 0 {
						if newCmd, ok := opts.GetHistoryEntry(historyIndex); ok {
							replaceLine(cmdToString(newCmd))
						} else {
							// something seems to have changed -> return to initial state
							historyIndex = -1
							clearLine()
						}
					} else {
						clearLine()
					}
				}
			}

		case console.KeyTab:
			if opts.GetCompletionOptions != nil {
				str := sb.String()
				cmd, _ := ParseCommand(fmt.Sprintf("%s%s", currentCommand, str))

				if len(cmd) == 0 {
					// append virtual entry to complete commands
					cmd = []string{""}
				} else if str[len(str)-1] == ' ' {
					// new command part already started by whitespace, but not recognized as part of command
					// -> append empty command part for processing
					cmd = append(cmd, "")
				}

				prefix := cmd[len(cmd)-1]
				options := filterOptions(opts.GetCompletionOptions(cmd, len(cmd)-1), prefix)
				if len(options) > 0 {
					if time.Since(lastTabPress) < doubleTabSpan {
						if opts.PrintOptionsHandler != nil {
							// double-tab detected -> print options
							console.Println()

							sort.Slice(options, func(i, j int) bool {
								return options[i].String() < options[j].String()
							})
							opts.PrintOptionsHandler(options)
							reprintLine()
						}
						// process next tab as single-press
						lastTabPress = time.Unix(0, 0)

					} else {
						if len(options) == 1 {
							if len(options[0].Replacement()) > 0 {
								suffix := Escape(options[0].Replacement()[len(prefix):])
								putString(suffix)

								if !options[0].IsPartial() {
									putRune(' ')
								}
							} else {
								// nothing changed? start double-tab combo
								lastTabPress = time.Now()
							}

						} else {
							longestCommonPrefix := findLongestCommonPrefix(options)
							suffix := Escape(longestCommonPrefix[len(prefix):])
							if len(suffix) > 0 {
								putString(suffix)
							} else {
								// nothing changed? start double-tab combo
								lastTabPress = time.Now()
							}
						}
					}
				} else {
					// nothing changed? start double-tab combo
					lastTabPress = time.Now()
				}
			}

		case console.KeyEnter:
			console.Println()
			return sb.String(), nil

		case console.KeyBackspace:
			removeLastChar()

		case console.KeySpace:
			putRune(' ')

		case 0:
			putRune(r)

		default:
			// ignore unknown special keys
		}
	}
}

func filterOptions(options []CompletionOption, prefix string) []CompletionOption {
	if options == nil {
		return nil
	}

	filtered := make([]CompletionOption, 0)
	for _, c := range options {
		if strings.HasPrefix(c.Replacement(), prefix) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func findLongestCommonPrefix(options []CompletionOption) string {
	if len(options) == 0 {
		return ""
	}

	longestCommonPrefix := ""
	for i := 1; ; i++ {
		if len(options[0].Replacement()) < i {
			// prefix cannot be any longer
			return longestCommonPrefix
		}

		prefix := options[0].Replacement()[:i]
		for _, c := range options {
			if !strings.HasPrefix(c.Replacement(), prefix) {
				// the next prefix would not be valid for all options
				return longestCommonPrefix
			}
		}

		longestCommonPrefix = prefix
	}
}

// ParseCommand parses a command input with escape sequences, single quotes and double quotes. The return parameter isComplete is false when a quote or escape sequence is not closed.
func ParseCommand(str string) (parts []string, isComplete bool) {
	cmd := make([]string, 0)

	var sb strings.Builder
	escape := false
	doubleQuote := false
	singleQuote := false

	for _, r := range str {
		switch {
		case singleQuote:
			if r == '\'' {
				singleQuote = false
			} else {
				sb.WriteRune(r)
			}
		case doubleQuote:
			if escape {
				if r != '\\' && r != '$' && r != '"' {
					// consume escape character only for actual escape sequences
					sb.WriteRune('\\')
				}
				sb.WriteRune(r)
				escape = false
			} else {
				if r == '"' {
					doubleQuote = false
				} else if r == '\\' {
					escape = true
				} else {
					sb.WriteRune(r)
				}
			}
		case escape:
			sb.WriteRune(r)
			escape = false
		default:
			if r == '\\' {
				escape = true
			} else if r == '\'' {
				singleQuote = true
			} else if r == '"' {
				doubleQuote = true
			} else if r == ' ' {
				if sb.Len() > 0 {
					cmd = append(cmd, sb.String())
					sb.Reset()
				}
			} else {
				sb.WriteRune(r)
			}
		}
	}

	if len(sb.String()) > 0 {
		cmd = append(cmd, sb.String())
	}

	return cmd, (!escape && !singleQuote && !doubleQuote)
}

// GetCommandString is the inverse function of Parse() and outputs a single string equal to the given command.
func GetCommandString(cmd []string) string {
	var sb strings.Builder
	for i, str := range cmd {
		if i > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(Quote(str))
	}
	return sb.String()
}

// Quote returns a quoted string if it contains special chars.
func Quote(str string) string {
	if NeedQuote(str) {
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(strings.ReplaceAll(str, "\\", "\\\\"), "\"", "\\\""))
	}

	return str
}

// NeedQuote returns true when the string contains characters that need to be quoted or escaped.
func NeedQuote(str string) bool {
	return strings.Contains(str, " ")
}

// Escape returns a string that escapes all special chars.
func Escape(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\"", "\\\"")
	str = strings.ReplaceAll(str, "'", "\\'")
	str = strings.ReplaceAll(str, " ", "\\ ")
	str = strings.ReplaceAll(str, "\n", "\\\n")
	str = strings.ReplaceAll(str, "\r", "\\\r")

	return str
}

// PromptHandler defines a function that returns the current command line prompt.
type PromptHandler func() string

// ExecUnknownCommandHandler is called when processing an unknown command.
type ExecUnknownCommandHandler func(cmd string, args []string) error

// CommandErrorHandler is called when a command has returned an error.
//
// Should return nil when the error has been handled, otherwise the command handler will stop and return the error.
type CommandErrorHandler func(cmd string, args []string, err error) error

// PrintOptionsHandler specifies a method to print options on double-tab.
type PrintOptionsHandler func([]CompletionOption)

// DefaultOptionsPrinter returns a function that is used to print options on double-tab.
//
// This method will ask the user for large lists to confirm printin.
func DefaultOptionsPrinter() PrintOptionsHandler {
	return func(options []CompletionOption) {
		if len(options) > maxAutoPrintListLen {
			console.Printlnf("  print all %d options? (y/N)", len(options)) //nolint
			// assume is only called during command reading here (keyboard needs to be prepared)
			_, r, err := console.ReadKey()
			if err != nil {
				return
			}

			if strings.ToLower(string(r)) != "y" {
				return
			}
		}

		console.PrintList(options) //nolint
	}
}
