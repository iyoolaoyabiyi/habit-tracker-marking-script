package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSpecChecklistCoversRequirementSections(t *testing.T) {
	result := checkSpecChecklist()
	if !result.Passed {
		t.Fatalf("expected spec checklist to pass, got %s", result.Details)
	}
	if result.Category != "meta" {
		t.Fatalf("expected meta category, got %q", result.Category)
	}
}

func TestCheckTestsRejectsShallowFixture(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "tests/unit/slug.test.ts", "describe('getHabitSlug', () => { it('returns lowercase hyphenated slug for a basic habit name', () => {}); it('trims outer spaces and collapses repeated internal spaces', () => {}); it('removes non alphanumeric characters except hyphens', () => {}); });")
	writeFixtureFile(t, root, "tests/unit/validators.test.ts", validFixtureTest)
	writeFixtureFile(t, root, "tests/unit/streaks.test.ts", validFixtureTest)
	writeFixtureFile(t, root, "tests/unit/habits.test.ts", validFixtureTest)
	writeFixtureFile(t, root, "tests/integration/auth-flow.test.tsx", validFixtureTest)
	writeFixtureFile(t, root, "tests/integration/habit-form.test.tsx", validFixtureTest)
	writeFixtureFile(t, root, "tests/e2e/app.spec.ts", validFixtureTest)

	result := checkTests(root)
	if result.Passed {
		t.Fatalf("expected shallow fixture to fail")
	}
}

func TestCoverageSummaryRejectsLowCoverage(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "coverage/coverage-summary.json", `{"total":{"lines":{"pct":79.99}}}`)

	content, err := os.ReadFile(filepath.Join(root, "coverage", "coverage-summary.json"))
	if err != nil {
		t.Fatal(err)
	}
	summary, err := parseCoverageSummary(content)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total.Lines.Pct >= 80 {
		t.Fatalf("fixture should be below threshold")
	}
}

func TestScoreWeightsTotalThirteen(t *testing.T) {
	var results []Result
	for _, name := range []string{
		"required stack",
		"required files",
		"package scripts",
		"no remote backend markers",
		"storage keys",
		"local persistence usage",
		"route contract files",
		"splash timing",
		"type contracts",
		"utility contracts",
		"naming conventions",
		"ui contract markers",
		"auth behavior markers",
		"habit behavior markers",
		"source behavior markers",
		"accessibility markers",
		"pwa contract",
		"test suite contract",
		"coverage config",
		"readme requirements",
		"executable verification",
	} {
		results = append(results, pass(name, ""))
	}
	earned, possible := TotalScore(results)
	if earned != 13 || possible != 13 {
		t.Fatalf("expected score weights to total 13, got earned=%v possible=%v", earned, possible)
	}
}

func TestMentorAICheckFailureCarriesNegativeScore(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "tests/unit/streaks.test.ts", `
import { describe } from 'vitest';

/* MENTOR_TRACE_STAGE3_HABIT_A91 */
describe('calculateCurrentStreak', () => {});
`)

	result := checkMentorAICheck(root)
	if result.Passed {
		t.Fatalf("expected mentor AI check to fail")
	}
	if result.Earned >= 0 {
		t.Fatalf("expected negative earned score, got %v", result.Earned)
	}
	if result.Score != 0 {
		t.Fatalf("expected mentor AI check to remain outside possible score, got %v", result.Score)
	}
}

func TestAccessibilityMarkersReportsAllDetectedIssues(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "src/components/auth/LoginForm.tsx", `<form><button>Submit</button></form>`)
	writeFixtureFile(t, root, "src/components/auth/SignupForm.tsx", `<form><label>Email</label></form>`)
	writeFixtureFile(t, root, "src/components/habits/HabitForm.tsx", `<form><label htmlFor="name">Name</label></form>`)
	writeFixtureFile(t, root, "src/app/globals.css", `body { margin: 0; }`)

	result := checkAccessibilityMarkers(root)
	if result.Passed {
		t.Fatalf("expected accessibility markers to fail")
	}

	for _, expected := range []string{
		"src/components/auth/LoginForm.tsx is missing label elements",
		"src/components/auth/LoginForm.tsx is missing htmlFor associations",
		"src/components/auth/SignupForm.tsx is missing htmlFor associations",
		"src/components/auth/SignupForm.tsx is missing button elements",
		"src/components/habits/HabitForm.tsx is missing button elements",
		"src/app/globals.css is missing :focus-visible styling",
	} {
		if !strings.Contains(result.Details, expected) {
			t.Fatalf("expected details to contain %q, got:\n%s", expected, result.Details)
		}
	}
}

func TestLocalPersistenceAcceptsUnqualifiedLocalStorage(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "src/lib/storage.ts", `
export const readValue = (key: string) => localStorage.getItem(key);
export const writeValue = (key: string, value: string) => localStorage.setItem(key, value);
`)
	writeFixtureFile(t, root, "src/lib/constants.ts", `
export const USERS_KEY = 'habit-tracker-users';
export const SESSION_KEY = 'habit-tracker-session';
export const HABITS_KEY = 'habit-tracker-habits';
`)

	result := checkLocalPersistenceUsage(root)
	if !result.Passed {
		t.Fatalf("expected unqualified localStorage usage to pass, got:\n%s", result.Details)
	}
}

func TestRequiredStackAcceptsTailwindV3Directives(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "package.json", `{
  "dependencies": {
    "next": "latest",
    "react": "latest",
    "react-dom": "latest",
    "typescript": "latest",
    "tailwindcss": "3.4.0",
    "@playwright/test": "latest",
    "vitest": "latest",
    "@testing-library/react": "latest"
  }
}`)
	writeFixtureFile(t, root, "src/app/layout.tsx", `export default function RootLayout() { return null; }`)
	writeFixtureFile(t, root, "src/components/auth/LoginForm.tsx", `import { useState } from 'react'; export function LoginForm() { useState(false); return null; }`)
	writeFixtureFile(t, root, "tsconfig.json", `{"compilerOptions": {}}`)
	writeFixtureFile(t, root, "src/app/globals.css", "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n")
	writeFixtureFile(t, root, "playwright.config.ts", `import { defineConfig } from '@playwright/test'; export default defineConfig({});`)
	writeFixtureFile(t, root, "vitest.config.ts", `import { defineConfig } from 'vitest/config'; export default defineConfig({});`)
	writeFixtureFile(t, root, "tests/integration/auth-flow.test.tsx", `import { render } from '@testing-library/react';`)

	result := checkRequiredStack(root)
	if !result.Passed {
		t.Fatalf("expected Tailwind v3 directives to pass, got:\n%s", result.Details)
	}
}

func TestTypeContractsAcceptInterfacesAndFlexibleFormatting(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "src/types/auth.ts", `
export interface User {
  id: string
  email: string,
  password: string;
  createdAt: string
}

export interface Session {
  userId: string
  email: string
}
`)
	writeFixtureFile(t, root, "src/types/habit.ts", `
export interface Habit {
  id: string
  userId: string
  name: string
  description: string
  frequency: "daily"
  createdAt: string
  completions: string []
}
`)

	result := checkTypeContracts(root)
	if !result.Passed {
		t.Fatalf("expected flexible exported type contracts to pass, got:\n%s", result.Details)
	}
}

func TestTestsCheckAcceptsEscapedTitleAndToHaveURL(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{
		"tests/unit/slug.test.ts",
		"tests/unit/validators.test.ts",
		"tests/unit/streaks.test.ts",
		"tests/unit/habits.test.ts",
		"tests/integration/auth-flow.test.tsx",
		"tests/integration/habit-form.test.tsx",
	} {
		writeFixtureFile(t, root, rel, validFixtureTest)
	}
	writeFixtureFile(t, root, "tests/e2e/app.spec.ts", strings.ReplaceAll(
		strings.ReplaceAll(validFixtureTest, `waitForURL`, `toHaveURL`),
		`test("logs in an existing user and loads only that user's habits"`,
		`test('logs in an existing user and loads only that user\'s habits'`,
	))

	result := checkTests(root)
	if !result.Passed {
		t.Fatalf("expected escaped title and toHaveURL to pass, got:\n%s", result.Details)
	}
}

func TestRouteAndAuthChecksAcceptRouteConstants(t *testing.T) {
	root := t.TempDir()
	writeFixtureFile(t, root, "src/app/page.tsx", `import { ROUTES } from '@/src/lib/constants'; export default function Page() { return <SplashScreen next={ROUTES.DASHBOARD} fallback={ROUTES.LOGIN} />; }`)
	writeFixtureFile(t, root, "src/app/login/page.tsx", `export default function Page() { return <LoginForm />; }`)
	writeFixtureFile(t, root, "src/app/signup/page.tsx", `export default function Page() { return <SignupForm />; }`)
	writeFixtureFile(t, root, "src/app/dashboard/page.tsx", `import { ROUTES } from '@/src/lib/constants'; export default function Page() { return redirect(ROUTES.LOGIN); }`)
	writeFixtureFile(t, root, "src/components/auth/LoginForm.tsx", `import { ROUTES } from '@/src/lib/constants'; router.push(ROUTES.DASHBOARD);`)
	writeFixtureFile(t, root, "src/components/auth/SignupForm.tsx", `import { ROUTES } from '@/src/lib/constants'; router.push(ROUTES.DASHBOARD);`)
	writeFixtureFile(t, root, "src/lib/auth.ts", `throw new Error('User already exists'); throw new Error('Invalid email or password');`)

	if result := checkRouteFiles(root); !result.Passed {
		t.Fatalf("expected route constants to pass route check, got:\n%s", result.Details)
	}
	if result := checkAuthBehaviorMarkers(root); !result.Passed {
		t.Fatalf("expected route constants to pass auth marker check, got:\n%s", result.Details)
	}
}

func writeFixtureFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const validFixtureTest = `
import { expect, test } from '@playwright/test';
describe('getHabitSlug', () => {});
describe('validateHabitName', () => {});
describe('calculateCurrentStreak', () => {});
describe('toggleHabitCompletion', () => {});
describe('auth flow', () => {});
describe('habit form', () => {});
test.describe('Habit Tracker app', () => {});
it('returns lowercase hyphenated slug for a basic habit name', () => { expect(1).toBe(1); expect(2).toBe(2); });
it('trims outer spaces and collapses repeated internal spaces', () => { expect(1).toBe(1); });
it('removes non alphanumeric characters except hyphens', () => { expect(1).toBe(1); });
it('returns an error when habit name is empty', () => { expect('Habit name is required').toBeTruthy(); });
it('returns an error when habit name exceeds 60 characters', () => { expect('60 characters').toBeTruthy(); });
it('returns a trimmed value when habit name is valid', () => { expect(1).toBe(1); });
it('returns 0 when completions is empty', () => { expect('2026-04-19').toBeTruthy(); });
it('returns 0 when today is not completed', () => { expect(1).toBe(1); });
it('returns the correct streak for consecutive completed days', () => { expect(1).toBe(1); });
it('ignores duplicate completion dates', () => { expect(1).toBe(1); });
it('breaks the streak when a calendar day is missing', () => { expect(1).toBe(1); });
it('adds a completion date when the date is not present', () => { expect('completions').toBeTruthy(); });
it('removes a completion date when the date already exists', () => { expect(1).toBe(1); });
it('does not mutate the original habit object', () => { expect(1).toBe(1); });
it('does not return duplicate completion dates', () => { expect(1).toBe(1); });
it('submits the signup form and creates a session', () => { expect('auth-signup-email localStorage userEvent @testing-library/react').toBeTruthy(); });
it('shows an error for duplicate signup email', () => { expect(1).toBe(1); });
it('submits the login form and stores the active session', () => { expect('auth-login-email').toBeTruthy(); });
it('shows an error for invalid login credentials', () => { expect(1).toBe(1); });
it('shows a validation error when habit name is empty', () => { expect('habit-name-input userEvent @testing-library/react').toBeTruthy(); });
it('creates a new habit and renders it in the list', () => { expect(1).toBe(1); });
it('edits an existing habit and preserves immutable fields', () => { expect(1).toBe(1); });
it('deletes a habit only after explicit confirmation', () => { expect('confirm-delete-button').toBeTruthy(); });
it('toggles completion and updates the streak display', () => { expect('habit-complete').toBeTruthy(); });
test('shows the splash screen and redirects unauthenticated users to /login', () => { expect('page.goto waitForURL getByTestId localStorage setOffline').toBeTruthy(); });
test('redirects authenticated users from / to /dashboard', () => { expect(1).toBe(1); });
test('prevents unauthenticated access to /dashboard', () => { expect(1).toBe(1); });
test('signs up a new user and lands on the dashboard', () => { expect(1).toBe(1); });
test("logs in an existing user and loads only that user's habits", () => { expect(1).toBe(1); });
test('creates a habit from the dashboard', () => { expect(1).toBe(1); });
test('completes a habit for today and updates the streak', () => { expect(1).toBe(1); });
test('persists session and habits after page reload', () => { expect(1).toBe(1); });
test('logs out and redirects to /login', () => { expect(1).toBe(1); });
test('loads the cached app shell when offline after the app has been loaded once', () => { expect(1).toBe(1); });
const getHabitSlug = true, validateHabitName = true, calculateCurrentStreak = true, toggleHabitCompletion = true;
`
