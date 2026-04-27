package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkTypeContracts(root string) Result {
	authContent, err := os.ReadFile(filepath.Join(root, "src/types/auth.ts"))
	if err != nil {
		return fail("type contracts", err.Error())
	}
	habitContent, err := os.ReadFile(filepath.Join(root, "src/types/habit.ts"))
	if err != nil {
		return fail("type contracts", err.Error())
	}

	authChecks := []string{
		"export type User = {",
		"id: string;",
		"email: string;",
		"password: string;",
		"createdAt: string;",
		"export type Session = {",
		"userId: string;",
	}
	for _, needle := range authChecks {
		if !strings.Contains(string(authContent), needle) {
			return fail("type contracts", "src/types/auth.ts missing "+needle)
		}
	}

	habitChecks := []string{
		"export type Habit = {",
		"userId: string;",
		"name: string;",
		"description: string;",
		"frequency: 'daily';",
		"createdAt: string;",
		"completions: string[];",
	}
	for _, needle := range habitChecks {
		if !strings.Contains(string(habitContent), needle) {
			return fail("type contracts", "src/types/habit.ts missing "+needle)
		}
	}

	return pass("type contracts", "required exported types found")
}
