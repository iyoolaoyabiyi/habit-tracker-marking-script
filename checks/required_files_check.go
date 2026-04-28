package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func checkRequiredFiles(root string) Result {
	required := []string{
		"package.json",
		"README.md",
		"src/app/globals.css",
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/app/login/page.tsx",
		"src/app/signup/page.tsx",
		"src/app/dashboard/page.tsx",
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
		"src/components/habits/HabitList.tsx",
		"src/components/habits/HabitCard.tsx",
		"src/components/shared/SplashScreen.tsx",
		"src/components/shared/ProtectedRoute.tsx",
		"src/lib/auth.ts",
		"src/lib/habits.ts",
		"src/lib/storage.ts",
		"src/lib/streaks.ts",
		"src/lib/slug.ts",
		"src/lib/validators.ts",
		"src/lib/constants.ts",
		"src/types/auth.ts",
		"src/types/habit.ts",
		"public/manifest.json",
		"public/sw.js",
		"public/icons/icon-192.png",
		"public/icons/icon-512.png",
		"tests/unit/slug.test.ts",
		"tests/unit/validators.test.ts",
		"tests/unit/streaks.test.ts",
		"tests/unit/habits.test.ts",
		"tests/integration/auth-flow.test.tsx",
		"tests/integration/habit-form.test.tsx",
		"tests/e2e/app.spec.ts",
	}

	var missing []string
	for _, rel := range required {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			missing = append(missing, rel)
		}
	}

	if len(missing) > 0 {
		return fail("required files", "missing: "+strings.Join(missing, ", "))
	}

	return pass("required files", fmt.Sprintf("%d required paths present", len(required)))
}
