package checks

import (
	"os"
	"path/filepath"
	"regexp"
)

func checkLocalPersistenceUsage(root string) Result {
	var issues []string
	storageText, storageOK := readFileIssue(root, "src/lib/storage.ts", &issues)
	constantsContent, _ := os.ReadFile(filepath.Join(root, "src/lib/constants.ts"))
	keyText := storageText + "\n" + string(constantsContent)

	storageRequired := []struct {
		label   string
		pattern string
	}{
		{label: "localStorage.getItem", pattern: `(?:window\.)?localStorage\.getItem\s*\(`},
		{label: "localStorage.setItem", pattern: `(?:window\.)?localStorage\.setItem\s*\(`},
	}
	if storageOK {
		for _, marker := range storageRequired {
			if !regexp.MustCompile(marker.pattern).MatchString(storageText) {
				issues = append(issues, "src/lib/storage.ts missing "+marker.label)
			}
		}
	}

	keyRequired := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	for _, marker := range keyRequired {
		if !containsAny(keyText, marker) {
			issues = append(issues, "src/lib/storage.ts or src/lib/constants.ts missing "+marker)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("local persistence usage", issues)
	}

	return pass("local persistence usage", "localStorage persistence is clearly implemented")
}
