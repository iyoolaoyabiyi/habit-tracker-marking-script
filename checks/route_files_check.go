package checks

import (
	"fmt"
	"strings"
)

func checkRouteFiles(root string) Result {
	type routeCheck struct {
		file     string
		contains []string
	}

	checks := []routeCheck{
		{file: "src/app/page.tsx", contains: []string{"SplashScreen", "/dashboard", "/login"}},
		{file: "src/app/login/page.tsx", contains: []string{"LoginForm"}},
		{file: "src/app/signup/page.tsx", contains: []string{"SignupForm"}},
		{file: "src/app/dashboard/page.tsx", contains: []string{"/login"}},
	}

	var issues []string
	for _, check := range checks {
		text, ok := readFileIssue(root, check.file, &issues)
		if !ok {
			continue
		}
		for _, needle := range check.contains {
			if !strings.Contains(text, needle) {
				issues = append(issues, fmt.Sprintf("%s missing %q", check.file, needle))
			}
		}
	}

	if len(issues) > 0 {
		return failWithIssues("route contract files", issues)
	}

	return pass("route contract files", "route files and core route markers found")
}
