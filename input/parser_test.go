package input_test

import (
	"testing"

	"github.com/DENICeG/go-console/v2/input"
	"github.com/stretchr/testify/require"
)

const testData = `\x1b[38;2;117;113;94m<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\x1b[0m\x1b[38;2;248;248;242m
\x1b[0m\x1b[38;2;249;38;114m<registry-request\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mxmlns=\x1b[0m\x1b[38;2;230;219;116m\"http://registry.denic.de/global/5.0\"\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mxmlns:domain=\x1b[0m\x1b[38;2;230;219;116m\"http://registry.denic.de/domain/5.0\"\x1b[0m\x1b[38;2;249;38;114m>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m<domain:info\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mrecursive=\x1b[0m\x1b[38;2;230;219;116m\"false\"\x1b[0m\x1b[38;2;249;38;114m>\x1b[0m\x1b[38;2;248;248;242m
      \x1b[0m\x1b[38;2;249;38;114m<domain:handle>\x1b[0m\x1b[38;2;248;248;242mdomain-example-1000022.de\x1b[0m\x1b[38;2;249;38;114m</domain:handle>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m</domain:info>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m</registry-request>\x1b[0m`

func Test_IsANSIIEscape(t *testing.T) {
	t.Run("ANSI Escape at start", func(t *testing.T) {
		result := input.IsANSIEscape(testData, 0)
		require.True(t, result, "Expected ANSII escape at start")
	})

	t.Run("ANSI Escape within text", func(t *testing.T) {
		const testData = `38;2;117;113;94m<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\x1b[0m\x1b[38;2;248;248;242m
\x1b[0m\x1b[38;2;249;38;114m<registry-request\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mxmlns=\x1b[0m\x1b[38;2;230;219;116m\"http://registry.denic.de/global/5.0\"\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mxmlns:domain=\x1b[0m\x1b[38;2;230;219;116m\"http://registry.denic.de/domain/5.0\"\x1b[0m\x1b[38;2;249;38;114m>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m<domain:info\x1b[0m\x1b[38;2;248;248;242m \x1b[0m\x1b[38;2;166;226;46mrecursive=\x1b[0m\x1b[38;2;230;219;116m\"false\"\x1b[0m\x1b[38;2;249;38;114m>\x1b[0m\x1b[38;2;248;248;242m
      \x1b[0m\x1b[38;2;249;38;114m<domain:handle>\x1b[0m\x1b[38;2;248;248;242mdomain-example-1000022.de\x1b[0m\x1b[38;2;249;38;114m</domain:handle>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m</domain:info>\x1b[0m\x1b[38;2;248;248;242m
   \x1b[0m\x1b[38;2;249;38;114m</registry-request>\x1b[0m`

		result := input.IsANSIEscape(testData, 77)
		require.True(t, result, "Expected ANSII escape at start")
	})
}

func Test_ReadANSISequence(t *testing.T) {
	t.Run("Read ANSI sequence", func(t *testing.T) {
		const expected = "\x1b[38;2;117;113;94m"
		const data = "\x1b[38;2;117;113;94m<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\x1b[0m\x1b[38;2;248;248;242m"
		ansiSequence, end := input.ReadANSISequence(data, 0)

		expectedPosition := len(expected)
		require.Equal(t, expectedPosition, end)
		require.Equal(t, "\x1b[38;2;117;113;94m", ansiSequence)
	})
}

func Test_IsANSIReset(t *testing.T) {
	t.Run("Is ANSI reset sequence", func(t *testing.T) {
		const testData = "[0m"
		result := input.IsANSIReset(testData, 0)
		require.True(t, result, "Expected ANSI reset sequence")
	})
}

func Test_ExtractRGB(t *testing.T) {
	t.Run("Extract RGB from ANSI escape", func(t *testing.T) {
		const testData = `\x1b[38;2;117;113;94m`
		rgb := input.ExtractRGB(testData)
		require.Equal(t, input.RGB{R: 117, G: 113, B: 94}, rgb)
	})

	t.Run("Extract RGB from invalid ANSI escape", func(t *testing.T) {
		const testData = `\x1b[38;2;117m`
		rgb := input.ExtractRGB(testData)
		require.Equal(t, input.RGB{}, rgb)
	})
}