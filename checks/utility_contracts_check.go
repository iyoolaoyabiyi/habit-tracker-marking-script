package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func checkUtilityContracts(root string) Result {
	utilityChecks := []struct {
		file     string
		patterns []string
		regexes  []string
	}{
		{
			file: "src/lib/slug.ts",
			regexes: []string{
				`export\s+function\s+getHabitSlug\s*\(\s*name\s*:\s*string\s*\)\s*:\s*string`,
			},
		},
		{
			file: "src/lib/validators.ts",
			patterns: []string{
				"Habit name is required",
				"Habit name must be 60 characters or fewer",
			},
			regexes: []string{
				`export\s+function\s+validateHabitName\s*\(\s*name\s*:\s*string\s*\)\s*:\s*\{`,
				`\.trim\s*\(`,
				`>\s*60|length\s*>\s*60`,
			},
		},
		{
			file: "src/lib/streaks.ts",
			regexes: []string{
				`export\s+function\s+calculateCurrentStreak\s*\(\s*completions\s*:\s*string\[\]\s*,\s*today`,
				`new\s+Set\s*\(\s*completions\s*\)`,
				`includes\s*\(\s*today\s*\)`,
			},
		},
		{
			file: "src/lib/habits.ts",
			regexes: []string{
				`export\s+function\s+toggleHabitCompletion\s*\(\s*habit\s*:\s*Habit\s*,\s*date\s*:\s*string\s*\)\s*:\s*Habit`,
				`new\s+Set\s*\(`,
				`\.\.\.\s*habit`,
			},
		},
	}

	for _, check := range utilityChecks {
		content, err := os.ReadFile(filepath.Join(root, check.file))
		if err != nil {
			return fail("utility contracts", err.Error())
		}
		text := string(content)
		for _, pattern := range check.patterns {
			if !strings.Contains(text, pattern) {
				return fail("utility contracts", fmt.Sprintf("%s missing %q", check.file, pattern))
			}
		}
		for _, pattern := range check.regexes {
			ok, err := regexp.MatchString(pattern, text)
			if err != nil {
				return fail("utility contracts", fmt.Sprintf("%s has invalid examiner regex %q: %v", check.file, pattern, err))
			}
			if !ok {
				return fail("utility contracts", fmt.Sprintf("%s missing pattern %q", check.file, pattern))
			}
		}
	}

	return pass("utility contracts", "required exported utility signatures and core implementation markers found")
}
