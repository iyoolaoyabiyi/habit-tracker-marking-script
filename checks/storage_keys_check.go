package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkStorageKeys(root string) Result {
	var issues []string
	var content strings.Builder
	readAny := false
	for _, rel := range []string{"src/lib/constants.ts", "src/lib/storage.ts"} {
		fileContent, err := os.ReadFile(filepath.Join(root, rel))
		if err == nil {
			content.Write(fileContent)
			content.WriteByte('\n')
			readAny = true
		}
	}
	if !readAny {
		issues = append(issues, "src/lib/constants.ts and src/lib/storage.ts could not be read")
		return failWithIssues("storage keys", issues)
	}

	required := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	for _, key := range required {
		if !strings.Contains(content.String(), key) {
			issues = append(issues, "src/lib/constants.ts or src/lib/storage.ts missing key "+key)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("storage keys", issues)
	}

	return pass("storage keys", "required localStorage keys found")
}
