package shake

import (
	"regexp"
	"strings"
)

var emptyLines *regexp.Regexp = regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)

// A sentence is a series of characters, starting with at lease one whitespace
// character, that ends in one of ., ! or ?
var sentence *regexp.Regexp = regexp.MustCompile(`\s+[^.!?]*[.!?]`)

func startAtSentence(s string, noSkip string) string {
	indices := sentence.FindAllStringIndex(s, -1)
	start := 0
	if len(indices) > 1 {
		start = indices[1][0] + 1 // 1 for space
	}

	if caseInsensitiveContains(s[:start], noSkip) {
		return s
	}

	return s[start:]
}

func endAtSentenceOrWord(s string) string {
	end := len(s)
	for i := end - 1; i >= 0; i-- {
		if s[i] == '?' || s[i] == '!' || s[i] == '.' {
			break
		}
		if s[i] == ' ' {
			end = i
			break
		}
	}

	return s[:end]
}

func caseInsensitiveContains(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func removeEmptyLines(s string) string {
	return emptyLines.ReplaceAllString(s, "\n")
}
