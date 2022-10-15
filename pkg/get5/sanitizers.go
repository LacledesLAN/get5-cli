package get5

import (
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

// sanitizeListItem ensures the input string is in the allowedList, otherwise returns an empty string.
func sanitizeListItem(allowList []string, input string) string {
	if allowList == nil || len(allowList) < 1 {
		return ""
	}

	if slices.Contains(allowList, input) {
		return input
	}

	return ""
}

// sanitizePrintable ensures a string is safe for printing in the CSGO game client.
func sanitizePrintable(raw string) string {
	raw = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}

		return -1
	}, raw)

	return strings.Join(strings.Fields(strings.TrimSpace(raw)), "_")
}

// sanitizeAndTruncatePrintable ensures a string is safe for printing in the CSGO game client and doesn't exceed a maxLength in size.
func sanitizeAndTruncatePrintable(raw string, maxLength int) string {
	raw = sanitizePrintable(raw)

	if maxLength < 1 || maxLength >= len(raw) {
		return raw
	}

	return raw[:maxLength]
}
