package checks

import (
	"fmt"
)

func checkRouteFiles(root string) Result {
	type routeCheck struct {
		file       string
		contains   [][]string
		labelNames []string
	}

	checks := []routeCheck{
		{
			file:       "src/app/page.tsx",
			contains:   [][]string{{"SplashScreen"}, {"/dashboard", "DASHBOARD"}, {"/login", "LOGIN"}},
			labelNames: []string{"SplashScreen", "/dashboard", "/login"},
		},
		{
			file:       "src/app/login/page.tsx",
			contains:   [][]string{{"LoginForm"}},
			labelNames: []string{"LoginForm"},
		},
		{
			file:       "src/app/signup/page.tsx",
			contains:   [][]string{{"SignupForm"}},
			labelNames: []string{"SignupForm"},
		},
		{
			file:       "src/app/dashboard/page.tsx",
			contains:   [][]string{{"/login", "LOGIN"}},
			labelNames: []string{"/login"},
		},
	}

	var issues []string
	for _, check := range checks {
		text, ok := readFileIssue(root, check.file, &issues)
		if !ok {
			continue
		}
		for index, alternatives := range check.contains {
			if !containsAny(text, alternatives...) {
				issues = append(issues, fmt.Sprintf("%s missing %q", check.file, check.labelNames[index]))
			}
		}
	}

	if len(issues) > 0 {
		return failWithIssues("route contract files", issues)
	}

	return pass("route contract files", "route files and core route markers found")
}
