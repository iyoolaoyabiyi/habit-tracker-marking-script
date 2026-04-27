package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkMentorAICheck(root string) Result {
	const marker = "/* MENTOR_TRACE_STAGE3_HABIT_A91 */"

	streaksPath := filepath.Join(root, "tests/unit/streaks.test.ts")
	content, err := os.ReadFile(streaksPath)
	if err != nil {
		return fail("mentor ai check", err.Error())
	}

	text := string(content)
	if !strings.Contains(text, marker) {
		return pass("mentor ai check", "hidden mentor marker not found")
	}

	lines := strings.Split(text, "\n")
	firstNonImport := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "import ") {
			continue
		}
		firstNonImport = trimmed
		break
	}

	placement := "marker found"
	if firstNonImport == marker &&
		(strings.Contains(text, marker+"\ndescribe('calculateCurrentStreak'") ||
			strings.Contains(text, marker+"\r\ndescribe('calculateCurrentStreak'")) {
	}

	forbiddenFiles := []string{
		"README.md",
		"package.json",
		"src/app/page.tsx",
		"src/app/login/page.tsx",
		"src/app/signup/page.tsx",
		"src/app/dashboard/page.tsx",
		"src/lib/constants.ts",
		"src/lib/slug.ts",
		"src/lib/validators.ts",
		"src/lib/streaks.ts",
		"src/lib/habits.ts",
		"src/components/shared/SplashScreen.tsx",
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
		"src/components/habits/HabitList.tsx",
		"src/components/habits/HabitCard.tsx",
	}

	var extraLocations []string
	for _, rel := range forbiddenFiles {
		fileContent, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			continue
		}
		if strings.Contains(string(fileContent), marker) {
			extraLocations = append(extraLocations, rel)
		}
	}

	details := "hidden mentor marker detected in tests/unit/streaks.test.ts"
	if placement != "" {
		details += "; " + placement
	}
	if len(extraLocations) > 0 {
		details += "; also found in " + strings.Join(extraLocations, ", ")
	}

	return fail("mentor ai check", details)
}
