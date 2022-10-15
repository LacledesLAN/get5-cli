package get5

import (
	"testing"
)

func Test_sanitizeListItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input     string
		allowList []string
		expected  string
	}{
		{"abc", nil, ""},
		{"abc", []string{}, ""},
		{"abc", []string{"abc"}, "abc"},
		{"abc", []string{"def"}, ""},
		{"abc", []string{"abc", "def"}, "abc"},
		{"abc", []string{"abc", "def", "def"}, "abc"},
		{"abc", []string{"abc", "def", ""}, "abc"},
	}

	for _, test := range tests {
		actual := sanitizeListItem(test.allowList, test.input)

		if actual != test.expected {
			t.Errorf("With input '%s' and allow list %s, expected %s, but got '%s'. ", test.input, test.allowList, test.expected, actual)
		}
	}
}

func Test_sanitizePrintable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello ", "hello"},
		{" hello ", "hello"},
		{" hello ", "hello"},
		{"\t\vhello\r\n", "hello"},
		{"hello    \r\n\t\v  there", "hello_there"},
		{" hello there ", "hello_there"},
		{"hello " + string('\a') + " there", "hello_there"},
	}

	for _, test := range tests {
		actual := sanitizePrintable(test.input)

		if actual != test.expected {
			t.Errorf("With input '%s' expected '%s', but got '%s'.",
				test.input, test.expected, actual)
		}
	}
}

func Test_sanitizeAndTruncatePrintable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"hello there", 0, "hello_there"},
		{"hello there", 1, "h"},
		{"hello there", 11, "hello_there"},
		{"hello there", 12, "hello_there"},
	}

	for _, test := range tests {
		actual := sanitizeAndTruncatePrintable(test.input, test.maxLength)

		if actual != test.expected {
			t.Errorf("With input '%s', and maxLength '%d', expected '%s', but got '%s'.",
				test.input, test.maxLength, test.expected, actual)
		}
	}
}
