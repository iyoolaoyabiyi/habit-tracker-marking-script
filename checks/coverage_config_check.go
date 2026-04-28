package checks

import "strings"

func checkCoverageConfig(root string) Result {
	var issues []string
	content, ok := readFileIssue(root, "vitest.config.ts", &issues)
	if !ok {
		return failWithIssues("coverage config", issues)
	}
	text := content
	required := []string{
		"coverage:",
		"thresholds:",
		"lines: 80",
		"src/lib",
	}
	for _, marker := range required {
		if !strings.Contains(text, marker) {
			issues = append(issues, "vitest.config.ts missing "+marker)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("coverage config", issues)
	}

	return pass("coverage config", "Vitest coverage threshold for src/lib is configured")
}
