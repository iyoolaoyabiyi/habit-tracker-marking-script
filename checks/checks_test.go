package checks

import (
	"os"
	"path/filepath"
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
