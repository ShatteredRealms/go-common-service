package util

// Reference: https://github.com/gojaguar/jaguar/tree/main/strings

func PascalCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	return string(append(t, lookupAndReplacePascalCaseWords(s, 0)...))
}

func lookupAndReplacePascalCaseWords(s string, i int) []byte {
	t := make([]byte, 0, 32)
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		t, i = appendLowercaseSequence(s, i, t)
	}
	return t
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// appendLowercaseSequence appends the lowercase sequence from s that begins at i into t
// returns the new t that contains all the chain of characters that should be lowercase
// and the new index where to start counting from.
func appendLowercaseSequence(s string, i int, t []byte) ([]byte, int) {
	for i+1 < len(s) && isASCIILower(s[i+1]) {
		i++
		t = append(t, s[i])
	}
	return t, i
}
