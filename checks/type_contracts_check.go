package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	authChecks := []struct {
		name   string
		fields []string
	}{
		{name: "User", fields: []string{"id: string;", "email: string;", "password: string;", "createdAt: string;"}},
		{name: "Session", fields: []string{"userId: string;", "email: string;"}},
	}
	for _, check := range authChecks {
		body, ok := exportedTypeBody(string(authContent), check.name)
		if !ok {
			return fail("type contracts", "src/types/auth.ts missing export type "+check.name)
		}
		for _, field := range check.fields {
			if !strings.Contains(body, field) {
				return fail("type contracts", fmt.Sprintf("src/types/auth.ts %s missing %s", check.name, field))
			}
		}
	}

	habitBody, ok := exportedTypeBody(string(habitContent), "Habit")
	if !ok {
		return fail("type contracts", "src/types/habit.ts missing export type Habit")
	}
	habitChecks := []string{
		"id: string;",
		"userId: string;",
		"name: string;",
		"description: string;",
		"frequency: 'daily';",
		"createdAt: string;",
		"completions: string[];",
	}
	for _, needle := range habitChecks {
		if !strings.Contains(habitBody, needle) {
			return fail("type contracts", "src/types/habit.ts Habit missing "+needle)
		}
	}

	return pass("type contracts", "required exported types found")
}

func exportedTypeBody(content, name string) (string, bool) {
	pattern := regexp.MustCompile(`export\s+type\s+` + regexp.QuoteMeta(name) + `\s*=\s*\{`)
	match := pattern.FindStringIndex(content)
	if match == nil {
		return "", false
	}
	rest := content[match[1]:]
	end := strings.Index(rest, "};")
	if end == -1 {
		return "", false
	}
	return rest[:end], true
}
