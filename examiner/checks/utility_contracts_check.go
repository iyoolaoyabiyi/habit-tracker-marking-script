package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkUtilityContracts(root string) Result {
	utilityChecks := []struct {
		file     string
		patterns []string
	}{
		{
			file: "src/lib/slug.ts",
			patterns: []string{
				"export function getHabitSlug(name: string): string",
			},
		},
		{
			file: "src/lib/validators.ts",
			patterns: []string{
				"export function validateHabitName(name: string): {",
				"Habit name is required",
				"Habit name must be 60 characters or fewer",
			},
		},
		{
			file: "src/lib/streaks.ts",
			patterns: []string{
				"export function calculateCurrentStreak(completions: string[], today",
			},
		},
		{
			file: "src/lib/habits.ts",
			patterns: []string{
				"export function toggleHabitCompletion(habit: Habit, date: string): Habit",
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
	}

	return pass("utility contracts", "required exported utility signatures found")
}
