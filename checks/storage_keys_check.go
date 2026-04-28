package checks

import "strings"

func checkStorageKeys(root string) Result {
	var issues []string
	content, ok := readFileIssue(root, "src/lib/constants.ts", &issues)
	if !ok {
		return failWithIssues("storage keys", issues)
	}

	required := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	for _, key := range required {
		if !strings.Contains(content, key) {
			issues = append(issues, "src/lib/constants.ts missing key "+key)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("storage keys", issues)
	}

	return pass("storage keys", "required localStorage keys found")
}
