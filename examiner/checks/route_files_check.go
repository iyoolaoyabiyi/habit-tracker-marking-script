package checks

import (
	"fmt"
	"os"
	"path/filepath"
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

	for _, check := range checks {
		content, err := os.ReadFile(filepath.Join(root, check.file))
		if err != nil {
			return fail("route contract files", err.Error())
		}
		text := string(content)
		for _, needle := range check.contains {
			if !strings.Contains(text, needle) {
				return fail("route contract files", fmt.Sprintf("%s missing %q", check.file, needle))
			}
		}
	}

	return pass("route contract files", "route files and core route markers found")
}
