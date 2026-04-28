package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkNoRemoteBackendMarkers(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return failWithIssues("no remote backend markers", []string{"package.json could not be read: " + err.Error()})
	}
	lower := strings.ToLower(string(content))
	forbidden := []string{
		"firebase",
		"supabase",
		"mongoose",
		"prisma",
		"aws-amplify",
		"next-auth",
		"auth0",
	}
	var issues []string
	for _, marker := range forbidden {
		if strings.Contains(lower, marker) {
			issues = append(issues, "package.json contains forbidden backend/auth marker "+marker)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("no remote backend markers", issues)
	}

	return pass("no remote backend markers", "no obvious remote database or external auth packages declared")
}
