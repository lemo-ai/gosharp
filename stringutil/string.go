// Package stringutil provides string helpers.
// Implementations live in package judge; this package re-exports them for compatibility.
package stringutil

import "github.com/lemo-ai/gosharp/judge"

// IsLetterUpper checks whether the given byte b is in upper case.
func IsLetterUpper(b byte) bool { return judge.IsLetterUpper(b) }

// IsLetterLower checks whether the given byte b is in lower case.
func IsLetterLower(b byte) bool { return judge.IsLetterLower(b) }

// IsLetter checks whether the given byte b is a letter.
func IsLetter(b byte) bool { return judge.IsLetter(b) }

// IsNumeric checks whether the given string s is numeric.
func IsNumeric(s string) bool { return judge.IsNumeric(s) }

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
func UcFirst(s string) string { return judge.UcFirst(s) }

// ReplaceByMap returns a copy of origin, replaced by a map in unordered way, case-sensitively.
func ReplaceByMap(origin string, replaces map[string]string) string {
	return judge.ReplaceByMap(origin, replaces)
}

// RemoveSymbols removes non-alphanumeric ASCII characters.
func RemoveSymbols(s string) string { return judge.RemoveSymbols(s) }

// EqualFoldWithoutChars checks whether s1 and s2 are equal case-insensitively
// after removing non-alphanumeric ASCII characters.
func EqualFoldWithoutChars(s1, s2 string) bool {
	return judge.EqualFoldWithoutChars(s1, s2)
}
