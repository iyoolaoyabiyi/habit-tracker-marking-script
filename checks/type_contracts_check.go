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
		return failWithIssues("type contracts", []string{"src/types/auth.ts could not be read: " + err.Error()})
	}
	habitContent, err := os.ReadFile(filepath.Join(root, "src/types/habit.ts"))
	if err != nil {
		return failWithIssues("type contracts", []string{"src/types/habit.ts could not be read: " + err.Error()})
	}

	authChecks := []struct {
		name   string
		fields []string
	}{
		{name: "User", fields: []string{"id: string;", "email: string;", "password: string;", "createdAt: string;"}},
		{name: "Session", fields: []string{"userId: string;", "email: string;"}},
	}
	var issues []string
	for _, check := range authChecks {
		body, ok := exportedTypeBody(string(authContent), check.name)
		if !ok {
			issues = append(issues, "src/types/auth.ts missing exported "+check.name+" type")
			continue
		}
		for _, field := range check.fields {
			if !typeBodyHasField(body, field) {
				issues = append(issues, fmt.Sprintf("src/types/auth.ts %s missing %s", check.name, field))
			}
		}
	}

	habitBody, ok := exportedTypeBody(string(habitContent), "Habit")
	if !ok {
		issues = append(issues, "src/types/habit.ts missing exported Habit type")
	} else {
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
			if !typeBodyHasField(habitBody, needle) {
				issues = append(issues, "src/types/habit.ts Habit missing "+needle)
			}
		}
	}

	if len(issues) > 0 {
		return failWithIssues("type contracts", issues)
	}

	return pass("type contracts", "required exported types found")
}

func exportedTypeBody(content, name string) (string, bool) {
	pattern := regexp.MustCompile(`export\s+(?:type\s+` + regexp.QuoteMeta(name) + `\s*=\s*|interface\s+` + regexp.QuoteMeta(name) + `\s*)\{`)
	match := pattern.FindStringIndex(content)
	if match == nil {
		return "", false
	}
	rest := content[match[1]:]
	depth := 1
	for index, char := range rest {
		switch char {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return rest[:index], true
			}
		}
	}
	return "", false
}

func typeBodyHasField(body, field string) bool {
	parts := strings.SplitN(strings.TrimSuffix(strings.TrimSpace(field), ";"), ":", 2)
	if len(parts) != 2 {
		return strings.Contains(body, field)
	}
	name := strings.TrimSpace(parts[0])
	fieldType := strings.TrimSpace(parts[1])
	typePattern := regexp.QuoteMeta(fieldType)
	typePattern = strings.ReplaceAll(typePattern, `string\[\]`, `string\s*\[\]`)
	typePattern = strings.ReplaceAll(typePattern, `'daily'`, `['"]daily['"]`)
	pattern := regexp.MustCompile(`(?m)\b` + regexp.QuoteMeta(name) + `\s*:\s*` + typePattern + `\s*[,;]?`)
	return pattern.MatchString(body)
}
