package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkHabitBehaviorMarkers(root string) Result {
	habitsContent, err := os.ReadFile(filepath.Join(root, "src/lib/habits.ts"))
	if err != nil {
		return fail("habit behavior markers", err.Error())
	}
	cardContent, err := os.ReadFile(filepath.Join(root, "src/components/habits/HabitCard.tsx"))
	if err != nil {
		return fail("habit behavior markers", err.Error())
	}

	required := []struct {
		text string
		src  string
	}{
		{"toggleHabitCompletion", "src/lib/habits.ts"},
		{"confirm-delete-button", "src/components/habits/HabitCard.tsx"},
		{"calculateCurrentStreak", "src/components/habits/HabitCard.tsx"},
	}

	sources := map[string]string{
		"src/lib/habits.ts":                   string(habitsContent),
		"src/components/habits/HabitCard.tsx": string(cardContent),
	}
	for _, item := range required {
		if !strings.Contains(sources[item.src], item.text) {
			return fail("habit behavior markers", fmt.Sprintf("%s missing %q", item.src, item.text))
		}
	}

	return pass("habit behavior markers", "habit completion, deletion confirmation, and streak markers are present")
}
