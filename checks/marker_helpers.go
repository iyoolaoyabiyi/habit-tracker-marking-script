package checks

import "strings"

func containsAny(text string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func normalizeJSEscapedText(text string) string {
	replacer := strings.NewReplacer(
		`\'`, `'`,
		`\"`, `"`,
		"\\`", "`",
	)
	return replacer.Replace(text)
}
