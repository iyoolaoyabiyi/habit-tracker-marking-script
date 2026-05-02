package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func checkTests(root string) Result {
	type testSpec struct {
		file     string
		patterns []string
		intents  []string
	}

	specs := []testSpec{
		{
			file: "tests/unit/slug.test.ts",
			patterns: []string{
				`describe('getHabitSlug'`,
				`returns lowercase hyphenated slug for a basic habit name`,
				`trims outer spaces and collapses repeated internal spaces`,
				`removes non alphanumeric characters except hyphens`,
			},
			intents: []string{`getHabitSlug`, `expect\s*\(`, `toBe\s*\(`},
		},
		{
			file: "tests/unit/validators.test.ts",
			patterns: []string{
				`describe('validateHabitName'`,
				`returns an error when habit name is empty`,
				`returns an error when habit name exceeds 60 characters`,
				`returns a trimmed value when habit name is valid`,
			},
			intents: []string{`validateHabitName`, `expect\s*\(`, `Habit name is required`, `60 characters`},
		},
		{
			file: "tests/unit/streaks.test.ts",
			patterns: []string{
				`describe('calculateCurrentStreak'`,
				`returns 0 when completions is empty`,
				`returns 0 when today is not completed`,
				`returns the correct streak for consecutive completed days`,
				`ignores duplicate completion dates`,
				`breaks the streak when a calendar day is missing`,
			},
			intents: []string{`calculateCurrentStreak`, `expect\s*\(`, `2026-[0-9]{2}-[0-9]{2}`},
		},
		{
			file: "tests/unit/habits.test.ts",
			patterns: []string{
				`describe('toggleHabitCompletion'`,
				`adds a completion date when the date is not present`,
				`removes a completion date when the date already exists`,
				`does not mutate the original habit object`,
				`does not return duplicate completion dates`,
			},
			intents: []string{`toggleHabitCompletion`, `expect\s*\(`, `completions`},
		},
		{
			file: "tests/integration/auth-flow.test.tsx",
			patterns: []string{
				`describe('auth flow'`,
				`submits the signup form and creates a session`,
				`shows an error for duplicate signup email`,
				`submits the login form and stores the active session`,
				`shows an error for invalid login credentials`,
			},
			intents: []string{`@testing-library/react`, `(?:userEvent|fireEvent)`, `localStorage`, `auth-login-email`, `auth-signup-email`, `expect\s*\(`},
		},
		{
			file: "tests/integration/habit-form.test.tsx",
			patterns: []string{
				`describe('habit form'`,
				`shows a validation error when habit name is empty`,
				`creates a new habit and renders it in the list`,
				`edits an existing habit and preserves immutable fields`,
				`deletes a habit only after explicit confirmation`,
				`toggles completion and updates the streak display`,
			},
			intents: []string{`@testing-library/react`, `(?:userEvent|fireEvent)`, `habit-name-input`, `confirm-delete-button`, `habit-complete`, `expect\s*\(`},
		},
		{
			file: "tests/e2e/app.spec.ts",
			patterns: []string{
				`test.describe('Habit Tracker app'`,
				`shows the splash screen and redirects unauthenticated users to /login`,
				`redirects authenticated users from / to /dashboard`,
				`prevents unauthenticated access to /dashboard`,
				`signs up a new user and lands on the dashboard`,
				`logs in an existing user and loads only that user's habits`,
				`creates a habit from the dashboard`,
				`completes a habit for today and updates the streak`,
				`persists session and habits after page reload`,
				`logs out and redirects to /login`,
				`loads the cached app shell when offline after the app has been loaded once`,
			},
			intents: []string{`@playwright/test`, `page\.goto`, `(?:waitForURL|toHaveURL)`, `getByTestId`, `localStorage`, `setOffline`, `expect\s*\(`},
		},
	}

	var issues []string
	for _, spec := range specs {
		content, err := os.ReadFile(filepath.Join(root, spec.file))
		if err != nil {
			issues = append(issues, spec.file+" could not be read: "+err.Error())
			continue
		}
		text := string(content)
		sourceText := normalizeJSEscapedText(text)
		for _, pattern := range spec.patterns {
			if !strings.Contains(sourceText, pattern) {
				issues = append(issues, fmt.Sprintf("%s missing %q", spec.file, pattern))
			}
		}
		for _, pattern := range spec.intents {
			ok, err := regexp.MatchString(pattern, text)
			if err != nil {
				return fail("test suite contract", fmt.Sprintf("%s has invalid examiner regex %q: %v", spec.file, pattern, err))
			}
			if !ok {
				issues = append(issues, fmt.Sprintf("%s missing intent evidence matching %q", spec.file, pattern))
			}
		}
		if assertions := regexp.MustCompile(`expect\s*\(`).FindAllStringIndex(text, -1); len(assertions) < 2 {
			issues = append(issues, fmt.Sprintf("%s has too few assertions to be meaningful", spec.file))
		}
	}

	if len(issues) > 0 {
		return failWithIssues("test suite contract", issues)
	}

	return pass("test suite contract", "required test files, titles, assertions, and intent markers found")
}
