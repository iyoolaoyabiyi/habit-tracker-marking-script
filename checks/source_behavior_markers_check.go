package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func checkSourceBehaviorMarkers(root string) Result {
	allSource, err := readSourceFiles(root)
	if err != nil {
		return fail("source behavior markers", err.Error())
	}

	files := map[string]string{}
	for _, rel := range []string{
		"src/lib/auth.ts",
		"src/lib/habits.ts",
		"src/lib/storage.ts",
		"src/lib/streaks.ts",
		"src/lib/validators.ts",
		"src/components/auth/LoginForm.tsx",
		"src/components/auth/SignupForm.tsx",
		"src/components/habits/HabitForm.tsx",
		"src/components/habits/HabitCard.tsx",
		"src/components/habits/HabitList.tsx",
		"src/components/shared/ProtectedRoute.tsx",
		"src/app/page.tsx",
		"src/app/dashboard/page.tsx",
	} {
		content, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return fail("source behavior markers", err.Error())
		}
		files[rel] = string(content)
	}

	checks := []struct {
		label   string
		text    string
		pattern string
	}{
		{"signup stores a user array", allSource, `(?s)writeUsers\s*\(|localStorage\.setItem\s*\([^)]*habit-tracker-users`},
		{"signup/login stores a session", allSource, `(?s)writeSession\s*\(|localStorage\.setItem\s*\([^)]*habit-tracker-session`},
		{"logout clears session", allSource, `(?s)(clearSession\s*\(|writeSession\s*\(\s*null\s*\)|removeItem\s*\([^)]*habit-tracker-session)`},
		{"duplicate signup is rejected", files["src/lib/auth.ts"], `(?s)(some|find)\s*\([^)]*email[^)]*\).*User already exists|User already exists.*(some|find)\s*\([^)]*email`},
		{"invalid login is rejected", files["src/lib/auth.ts"], `Invalid email or password`},
		{"dashboard protects unauthenticated users", allSource, `/login`},
		{"dashboard filters habits by userId", allSource, `(?s)(filter\s*\([^)]*userId|habit\.userId\s*={2,3}|habit\.userId\s*===|userId\s*===\s*habit\.userId)`},
		{"habit creation assigns owner", allSource, `(?s)userId\s*[:,]`},
		{"habit frequency defaults to daily", allSource, `frequency\s*:\s*['"]daily['"]`},
		{"habit edit preserves immutable fields", allSource, `(?s)\.\.\.\s*habit|id\s*:\s*habit\.id|createdAt\s*:\s*habit\.createdAt|completions\s*:\s*habit\.completions`},
		{"delete requires explicit confirmation", allSource, `confirm-delete-button`},
		{"completion uses today's date", allSource, `(?s)(getToday|today|new Date\(\))`},
		{"completion de-duplicates dates", files["src/lib/habits.ts"], `new Set\s*\(`},
		{"streak checks today before counting", files["src/lib/streaks.ts"], `(?s)includes\s*\(\s*today\s*\).*return\s+0|return\s+0.*includes\s*\(\s*today\s*\)`},
		{"validator trims habit names", files["src/lib/validators.ts"], `\.trim\s*\(`},
	}

	var missing []string
	for _, check := range checks {
		if ok, err := regexp.MatchString(check.pattern, check.text); err != nil {
			return fail("source behavior markers", fmt.Sprintf("invalid matcher for %s: %v", check.label, err))
		} else if !ok {
			missing = append(missing, check.label)
		}
	}

	if len(missing) > 0 {
		return fail("source behavior markers", "missing behavior evidence: "+strings.Join(missing, ", "))
	}

	return pass("source behavior markers", "auth, storage, habit ownership, CRUD, completion, and streak behavior evidence found")
}

func readSourceFiles(root string) (string, error) {
	var builder strings.Builder
	srcRoot := filepath.Join(root, "src")
	err := filepath.WalkDir(srcRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".tsx") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		builder.WriteString(filepath.ToSlash(path))
		builder.WriteByte('\n')
		builder.Write(content)
		builder.WriteByte('\n')
		return nil
	})
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}
