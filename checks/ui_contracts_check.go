package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkUIContracts(root string) Result {
	fileChecks := []struct {
		file     string
		patterns []string
	}{
		{
			file:     "src/components/shared/SplashScreen.tsx",
			patterns: []string{`data-testid="splash-screen"`},
		},
		{
			file:     "src/components/auth/LoginForm.tsx",
			patterns: []string{`auth-login-email`, `auth-login-password`, `auth-login-submit`},
		},
		{
			file:     "src/components/auth/SignupForm.tsx",
			patterns: []string{`auth-signup-email`, `auth-signup-password`, `auth-signup-submit`},
		},
		{
			file:     "src/components/habits/HabitForm.tsx",
			patterns: []string{`habit-form`, `habit-name-input`, `habit-description-input`, `habit-frequency-select`, `habit-save-button`},
		},
		{
			file:     "src/components/habits/HabitCard.tsx",
			patterns: []string{`habit-card-`, `habit-streak-`, `habit-complete-`, `habit-edit-`, `habit-delete-`, `confirm-delete-button`},
		},
	}

	var issues []string
	for _, check := range fileChecks {
		content, err := os.ReadFile(filepath.Join(root, check.file))
		if err != nil {
			issues = append(issues, check.file+" could not be read: "+err.Error())
			continue
		}
		text := string(content)
		for _, pattern := range check.patterns {
			if !strings.Contains(text, pattern) {
				issues = append(issues, fmt.Sprintf("%s missing %q", check.file, pattern))
			}
		}
	}

	splashBytes, splashErr := os.ReadFile(filepath.Join(root, "src/components/shared/SplashScreen.tsx"))
	constantsContent, _ := os.ReadFile(filepath.Join(root, "src/lib/constants.ts"))
	if splashErr == nil && !strings.Contains(string(splashBytes)+"\n"+string(constantsContent), "Habit Tracker") {
		issues = append(issues, `src/components/shared/SplashScreen.tsx or src/lib/constants.ts missing "Habit Tracker"`)
	}

	dashboardText, err := readExistingFiles(root, []string{
		"src/app/dashboard/page.tsx",
		"src/components/habits/HabitList.tsx",
		"src/components/habits/HabitForm.tsx",
		"src/components/habits/HabitCard.tsx",
	})
	if err != nil {
		issues = append(issues, "dashboard UI files could not be read: "+err.Error())
	} else {
		for _, pattern := range []string{`dashboard-page`, `empty-state`, `create-habit-button`, `auth-logout-button`} {
			if !strings.Contains(dashboardText, pattern) {
				issues = append(issues, fmt.Sprintf("dashboard UI missing %q", pattern))
			}
		}
	}

	if len(issues) > 0 {
		return failWithIssues("ui contract markers", issues)
	}

	return pass("ui contract markers", "required test ids and UI markers found")
}
