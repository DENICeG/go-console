package input

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// IsANSIEscape checks if the given rune slice contains an ANSI escape sequence on the given index.
func IsANSIEscape(input string, start int) bool {
	const expectedFirstRune = rune('\x1b')
	firstRune := rune(input[start])
	if firstRune != expectedFirstRune {
		return false
	}

	if utf8.RuneCountInString(input) < start+1 {
		return false
	}

	const expectedSecondRune = rune('[')
	secondRune := rune(input[start+1])
	return secondRune == expectedSecondRune
}

// ReadANSISequence reads an ANSI escape sequence from the given rune slice and returns the sequence and the end index.
func ReadANSISequence(input string, start int) (string, int) {
	var sb strings.Builder

	for i := start; i < utf8.RuneCountInString(input); i++ {
		if input[i] == 'm' {
			sb.WriteByte(input[i])
			return sb.String(), i + 1
		}

		sb.WriteByte(input[i])
	}

	return "", -1 // return the start index as the end index for invalid ANSI escape sequences
}

// ExtractRGB extracts the RGB values from an ANSI escape sequence.
func ExtractRGB(input string) RGB {
	regex := regexp.MustCompile("\x1b\\[(?P<FLAGS>\\d{1,3});(?P<MODE>\\d{1,3});(?P<R>\\d{1,3});(?P<G>\\d{1,3});(?P<B>\\d{1,3})m")
	matches := regex.FindStringSubmatch(input)

	if len(matches) != 6 {
		return RGB{}
	}

	rString := matches[3]
	gString := matches[4]
	bString := matches[5]

	r, _ := strconv.Atoi(rString)
	g, _ := strconv.Atoi(gString)
	b, _ := strconv.Atoi(bString)

	return RGB{
		R: uint8(r), //nolint
		G: uint8(g), //nolint
		B: uint8(b), //nolint
	}
}

// IsANSIReset checks if the given rune slice contains an ANSI reset sequence on the given index.
func IsANSIReset(input string, start int) bool {
	if utf8.RuneCountInString(input) < start+3 {
		return false
	}

	const expectedFirstRune = rune('\x1b')
	firstRune := rune(input[start])
	if firstRune != expectedFirstRune {
		return false
	}

	const expectedSecondRune = rune('[')
	secondRune := rune(input[start+1])
	match := secondRune == expectedSecondRune
	if !match {
		return false
	}

	return input[start+2] == '0' && input[start+3] == 'm'
}
