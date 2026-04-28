package checks

import "strings"

func checkLocalPersistenceUsage(root string) Result {
	var issues []string
	storageText, storageOK := readFileIssue(root, "src/lib/storage.ts", &issues)
	constantsText, constantsOK := readFileIssue(root, "src/lib/constants.ts", &issues)

	storageRequired := []string{
		"window.localStorage.getItem",
		"window.localStorage.setItem",
	}
	if storageOK {
		for _, marker := range storageRequired {
			if !strings.Contains(storageText, marker) {
				issues = append(issues, "src/lib/storage.ts missing "+marker)
			}
		}
	}

	keyRequired := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	if constantsOK {
		for _, marker := range keyRequired {
			if !strings.Contains(constantsText, marker) {
				issues = append(issues, "src/lib/constants.ts missing "+marker)
			}
		}
	}

	if len(issues) > 0 {
		return failWithIssues("local persistence usage", issues)
	}

	return pass("local persistence usage", "localStorage persistence is clearly implemented")
}
