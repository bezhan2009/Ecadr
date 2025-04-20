package aiSerivce

import "strings"

func addBrackets(text string) string {
	trimmedText := strings.TrimSpace(text)

	if !strings.HasPrefix(trimmedText, "[") {
		trimmedText = "[" + trimmedText
	}
	if !strings.HasSuffix(trimmedText, "]") {
		trimmedText = trimmedText + "]"
	}

	return trimmedText
}
