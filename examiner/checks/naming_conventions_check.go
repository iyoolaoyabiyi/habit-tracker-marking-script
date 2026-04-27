package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func checkNamingConventions(root string) Result {
	componentFiles := []string{
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
		"src/components/habits/HabitList.tsx",
		"src/components/habits/HabitCard.tsx",
		"src/components/shared/SplashScreen.tsx",
		"src/components/shared/ProtectedRoute.tsx",
	}
	for _, rel := range componentFiles {
		base := filepath.Base(rel)
		if matched, _ := regexp.MatchString(`^[A-Z][A-Za-z0-9]*\.tsx$`, base); !matched {
			return fail("naming conventions", fmt.Sprintf("component file must use PascalCase: %s", rel))
		}
	}

	libFiles := []string{
		"src/lib/auth.ts",
		"src/lib/habits.ts",
		"src/lib/storage.ts",
		"src/lib/streaks.ts",
		"src/lib/slug.ts",
		"src/lib/validators.ts",
		"src/lib/constants.ts",
	}
	for _, rel := range libFiles {
		base := filepath.Base(rel)
		if matched, _ := regexp.MatchString(`^[a-z0-9]+\.ts$`, base); !matched {
			return fail("naming conventions", fmt.Sprintf("utility file must be lowercase: %s", rel))
		}
	}

	testFiles := []string{
		"tests/unit/slug.test.ts",
		"tests/unit/validators.test.ts",
		"tests/unit/streaks.test.ts",
		"tests/unit/habits.test.ts",
		"tests/integration/auth-flow.test.tsx",
		"tests/integration/habit-form.test.tsx",
		"tests/e2e/app.spec.ts",
	}
	for _, rel := range testFiles {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			return fail("naming conventions", "missing required test file "+rel)
		}
	}

	return pass("naming conventions", "component, utility, and required test file names match the contract")
}
