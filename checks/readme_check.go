package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkReadme(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		return fail("readme requirements", err.Error())
	}

	text := strings.ToLower(string(content))
	requiredPhrases := []string{
		"project overview",
		"setup instructions",
		"run instructions",
		"test instructions",
		"local persistence",
		"pwa support",
		"trade-offs",
		"test file mapping",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(text, phrase) {
			return fail("readme requirements", "README.md missing section containing "+phrase)
		}
	}

	return pass("readme requirements", "required README sections found")
}
