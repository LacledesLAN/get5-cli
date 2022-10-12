package get5

import "testing"

func Test_sanitizePrintable(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"hello", 0, "hello"},
		{"hello", 10, "hello"},
		{"hello", 1, "h"},
		{"hello ", 0, "hello"},
		{" hello ", 0, "hello"},
		{" hello ", 0, "hello"},
		{"\t\vhello\r\n", 0, "hello"},
		{" hello there ", 0, "hello_there"},
	}

	for _, test := range tests {
		actual := sanitizePrintable(test.input, test.maxLength)

		if actual != test.expected {
			t.Errorf("With input '%s' and max length %d, expected '%s', but got '%s'.",
				test.input, test.maxLength, test.expected, actual)
		}
	}
}
