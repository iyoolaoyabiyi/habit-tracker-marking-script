package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runExaminerRuntimeTests(root string) (string, error) {
	testPath := filepath.Join(root, "tests", "unit", "examiner-runtime.test.ts")
	cleanup, err := writeTemporaryFile(testPath, examinerRuntimeTestSource)
	if err != nil {
		return "", err
	}
	defer cleanup()

	return runCommand(root, "npm", "exec", "--", "vitest", "run", "tests/unit/examiner-runtime.test.ts", "--environment", "jsdom")
}

func runExaminerBrowserTests(root string, options RuntimeOptions) (string, error) {
	testPath := filepath.Join(root, "tests", "e2e", "examiner-browser.spec.ts")
	source := examinerBrowserTestSource
	source = strings.ReplaceAll(source, "__RUN_OFFLINE__", fmt.Sprintf("%t", options.RunOffline))
	source = strings.ReplaceAll(source, "__RUN_ACCESSIBILITY__", fmt.Sprintf("%t", options.RunAccessibility))
	source = strings.ReplaceAll(source, "__RUN_RESPONSIVE__", fmt.Sprintf("%t", options.RunResponsive))

	cleanup, err := writeTemporaryFile(testPath, source)
	if err != nil {
		return "", err
	}
	defer cleanup()

	return runCommand(root, "npm", "exec", "--", "playwright", "test", "tests/e2e/examiner-browser.spec.ts")
}

func writeTemporaryFile(path, content string) (func(), error) {
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("refusing to overwrite existing file %s", path)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return nil, err
	}

	return func() {
		_ = os.Remove(path)
	}, nil
}

const examinerRuntimeTestSource = `import { describe, expect, it } from 'vitest';
import { getHabitSlug } from '@/src/lib/slug';
import { validateHabitName } from '@/src/lib/validators';
import { calculateCurrentStreak } from '@/src/lib/streaks';
import { toggleHabitCompletion } from '@/src/lib/habits';
import type { Habit } from '@/src/types/habit';

describe('examiner-owned utility contracts', () => {
  const baseHabit: Habit = {
    id: 'habit-1',
    userId: 'user-1',
    name: 'Drink Water',
    description: '',
    frequency: 'daily',
    createdAt: '2026-04-19T00:00:00.000Z',
    completions: [],
  };

  it('verifies slug normalization rules independently', () => {
    expect(getHabitSlug('  Drink   Water! #1  ')).toBe('drink-water-1');
  });

  it('verifies habit name validation independently', () => {
    expect(validateHabitName('   ')).toMatchObject({ valid: false, value: '', error: 'Habit name is required' });
    expect(validateHabitName('a'.repeat(61))).toMatchObject({ valid: false, error: 'Habit name must be 60 characters or fewer' });
    expect(validateHabitName('  Read Books  ')).toEqual({ valid: true, value: 'Read Books', error: null });
  });

  it('verifies current streak rules independently', () => {
    expect(calculateCurrentStreak([], '2026-04-19')).toBe(0);
    expect(calculateCurrentStreak(['2026-04-18'], '2026-04-19')).toBe(0);
    expect(calculateCurrentStreak(['2026-04-17', '2026-04-18', '2026-04-19'], '2026-04-19')).toBe(3);
    expect(calculateCurrentStreak(['2026-04-17', '2026-04-19'], '2026-04-19')).toBe(1);
  });

  it('verifies completion toggling independently', () => {
    const completed = toggleHabitCompletion(baseHabit, '2026-04-19');
    expect(completed.completions).toEqual(['2026-04-19']);
    expect(baseHabit.completions).toEqual([]);
    expect(toggleHabitCompletion({ ...completed, completions: ['2026-04-19', '2026-04-19'] }, '2026-04-20').completions).toEqual(['2026-04-19', '2026-04-20']);
    expect(toggleHabitCompletion(completed, '2026-04-19').completions).toEqual([]);
  });
});
`

const examinerBrowserTestSource = `import { expect, test } from '@playwright/test';

const existingUser = {
  id: 'user-1',
  email: 'existing@example.com',
  password: 'password123',
  createdAt: '2026-04-19T00:00:00.000Z',
};

const otherUser = {
  id: 'user-2',
  email: 'other@example.com',
  password: 'password123',
  createdAt: '2026-04-19T00:00:00.000Z',
};

const existingHabit = {
  id: 'habit-1',
  userId: 'user-1',
  name: 'Drink Water',
  description: 'Eight glasses',
  frequency: 'daily',
  createdAt: '2026-04-19T00:00:00.000Z',
  completions: [],
};

async function seedStorage(page, values) {
  await page.addInitScript((storageValues) => {
    window.localStorage.setItem('habit-tracker-users', JSON.stringify(storageValues.users ?? []));
    window.localStorage.setItem('habit-tracker-session', JSON.stringify(storageValues.session ?? null));
    window.localStorage.setItem('habit-tracker-habits', JSON.stringify(storageValues.habits ?? []));
  }, values);
}

test.describe('examiner-owned browser verification', () => {
  test('verifies auth routing, user isolation, CRUD, completion, reload, and logout', async ({ page, context }) => {
    await context.clearCookies();
    await seedStorage(page, {
      users: [existingUser, otherUser],
      session: null,
      habits: [
        existingHabit,
        { ...existingHabit, id: 'habit-2', userId: 'user-2', name: 'Other Habit' },
      ],
    });

    await page.goto('/');
    await expect(page.getByTestId('splash-screen')).toBeVisible();
    await page.waitForURL('**/login');

    await page.goto('/dashboard');
    await page.waitForURL('**/login');

    await page.getByTestId('auth-login-email').fill(existingUser.email);
    await page.getByTestId('auth-login-password').fill(existingUser.password);
    await page.getByTestId('auth-login-submit').click();
    await page.waitForURL('**/dashboard');
    await expect(page.getByTestId('dashboard-page')).toBeVisible();
    await expect(page.getByTestId('habit-card-drink-water')).toBeVisible();
    await expect(page.getByTestId('habit-card-other-habit')).toHaveCount(0);

    await page.getByTestId('create-habit-button').click();
    await page.getByTestId('habit-name-input').fill('Read Books');
    await page.getByTestId('habit-description-input').fill('Ten pages');
    await page.getByTestId('habit-save-button').click();
    await expect(page.getByTestId('habit-card-read-books')).toBeVisible();

    await page.getByTestId('habit-edit-read-books').click();
    await page.getByTestId('habit-name-input').fill('Read More');
    await page.getByTestId('habit-save-button').click();
    await expect(page.getByTestId('habit-card-read-more')).toBeVisible();

    await page.getByTestId('habit-complete-read-more').click();
    await expect(page.getByTestId('habit-streak-read-more')).toContainText('Streak: 1');
    await page.getByTestId('habit-complete-read-more').click();
    await expect(page.getByTestId('habit-streak-read-more')).toContainText('Streak: 0');

    await page.reload();
    await expect(page.getByTestId('dashboard-page')).toBeVisible();
    await expect(page.getByTestId('habit-card-read-more')).toBeVisible();

    await page.getByTestId('habit-delete-read-more').click();
    await expect(page.getByTestId('confirm-delete-button')).toBeVisible();
    await page.getByTestId('confirm-delete-button').click();
    await expect(page.getByTestId('habit-card-read-more')).toHaveCount(0);

    await page.getByTestId('auth-logout-button').click();
    await page.waitForURL('**/login');
    await expect(page.getByTestId('auth-login-submit')).toBeVisible();
  });

  test('verifies PWA manifest, icons, service worker, and offline shell', async ({ page, context, request }) => {
    test.skip(!__RUN_OFFLINE__, 'offline/PWA examiner check disabled');
    await context.clearCookies();
    await page.goto('/');
    await expect(page.locator('link[rel="manifest"]')).toHaveAttribute('href', /manifest\.json/);

    const manifest = await request.get('/manifest.json');
    expect(manifest.ok()).toBeTruthy();
    const manifestBody = await manifest.json();
    expect(manifestBody.icons.some((icon) => String(icon.sizes).includes('192x192'))).toBeTruthy();
    expect(manifestBody.icons.some((icon) => String(icon.sizes).includes('512x512'))).toBeTruthy();
    for (const icon of manifestBody.icons) {
      const response = await request.get(icon.src);
      expect(response.ok()).toBeTruthy();
    }

    await page.waitForFunction(() => 'serviceWorker' in navigator);
    await page.evaluate(async () => Boolean(await navigator.serviceWorker.ready));
    await context.setOffline(true);
    await page.goto('/', { waitUntil: 'domcontentloaded' });
    await expect(page.getByTestId('splash-screen')).toBeVisible();
    await context.setOffline(false);
  });

  test('verifies accessible controls', async ({ page }) => {
    test.skip(!__RUN_ACCESSIBILITY__, 'accessibility examiner check disabled');
    await seedStorage(page, {
      users: [existingUser],
      session: { userId: existingUser.id, email: existingUser.email },
      habits: [existingHabit],
    });

    await page.goto('/dashboard');
    await expect(page.getByTestId('dashboard-page')).toBeVisible();
    await page.getByTestId('create-habit-button').click();
    const unlabeledControls = await page.locator('input, textarea, select').evaluateAll((controls) =>
      controls
        .filter((control) => {
          const id = control.getAttribute('id');
          const ariaLabel = control.getAttribute('aria-label');
          const ariaLabelledBy = control.getAttribute('aria-labelledby');
          const hasLabel = id ? Boolean(document.querySelector('label[for="' + CSS.escape(id) + '"]')) : false;
          return !ariaLabel && !ariaLabelledBy && !hasLabel;
        })
        .map((control) => control.getAttribute('data-testid') || control.tagName),
    );
    expect(unlabeledControls).toEqual([]);

    await page.keyboard.press('Tab');
    const activeElementTag = await page.evaluate(() => document.activeElement?.tagName);
    expect(['A', 'BUTTON', 'INPUT', 'TEXTAREA', 'SELECT']).toContain(activeElementTag);
  });

  test('verifies responsive layout at required widths', async ({ page }) => {
    test.skip(!__RUN_RESPONSIVE__, 'responsive examiner check disabled');
    await seedStorage(page, {
      users: [existingUser],
      session: { userId: existingUser.id, email: existingUser.email },
      habits: [existingHabit],
    });

    for (const width of [320, 768, 1280]) {
      await page.setViewportSize({ width, height: 720 });
      await page.goto('/dashboard');
      await expect(page.getByTestId('dashboard-page')).toBeVisible();
      const hasHorizontalOverflow = await page.evaluate(() => document.documentElement.scrollWidth > window.innerWidth);
      expect(hasHorizontalOverflow).toBeFalsy();
    }
  });
});
`
