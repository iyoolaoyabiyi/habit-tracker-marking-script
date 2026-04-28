package checks

import (
	"fmt"
	"strings"
)

func checkHabitBehaviorMarkers(root string) Result {
	var issues []string
	habitsContent, _ := readFileIssue(root, "src/lib/habits.ts", &issues)
	cardContent, _ := readFileIssue(root, "src/components/habits/HabitCard.tsx", &issues)

	required := []struct {
		text string
		src  string
	}{
		{"toggleHabitCompletion", "src/lib/habits.ts"},
		{"confirm-delete-button", "src/components/habits/HabitCard.tsx"},
		{"calculateCurrentStreak", "src/components/habits/HabitCard.tsx"},
	}

	sources := map[string]string{
		"src/lib/habits.ts":                   habitsContent,
		"src/components/habits/HabitCard.tsx": cardContent,
	}
	for _, item := range required {
		if !strings.Contains(sources[item.src], item.text) {
			issues = append(issues, fmt.Sprintf("%s missing %q", item.src, item.text))
		}
	}

	if len(issues) > 0 {
		return failWithIssues("habit behavior markers", issues)
	}

	return pass("habit behavior markers", "habit completion, deletion confirmation, and streak markers are present")
}
