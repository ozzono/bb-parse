package utils

import "strings"

func TrimSpaces(input string) string {
	for strings.Contains(input, "  ") {
		input = strings.ReplaceAll(input, "  ", " ")
	}
	input = strings.TrimPrefix(input, " ")
	input = strings.TrimSuffix(input, " ")
	return input
}
