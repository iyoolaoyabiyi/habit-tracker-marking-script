package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkCoverageConfig(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "vitest.config.ts"))
	if err != nil {
		return fail("coverage config", err.Error())
	}
	text := string(content)
	required := []string{
		"coverage:",
		"thresholds:",
		"lines: 80",
		"src/lib",
	}
	for _, marker := range required {
		if !strings.Contains(text, marker) {
			return fail("coverage config", "vitest.config.ts missing "+marker)
		}
	}

	return pass("coverage config", "Vitest coverage threshold for src/lib is configured")
}
