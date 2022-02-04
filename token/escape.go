package token

import (
	"strings"
)

const EscapeChar byte = '\\'

var unescapeReplacer = strings.NewReplacer(
	"\\ ", " ",
	"\\:", ":",
	"\\(", "(",
	"\\)", ")",
	"\\\\", "\\",
)

func Unescape(s string) string {
	if !ContainsSpecialChars(s) {
		return s
	}
	return unescapeReplacer.Replace(s)
}

var escapeReplacer = strings.NewReplacer(
	" ", "\\ ",
	":", "\\:",
	"(", "\\(",
	")", "\\)",
	"\\", "\\\\",
)

func Escape(s string) string {
	if !ContainsSpecialChars(s) {
		return s
	}
	return escapeReplacer.Replace(s)
}

var specialChars string = " :()\\"

func ContainsSpecialChars(s string) bool {
	return strings.ContainsAny(s, specialChars)
}
