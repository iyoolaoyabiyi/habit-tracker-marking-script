package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkNoRemoteBackendMarkers(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return fail("no remote backend markers", err.Error())
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
	for _, marker := range forbidden {
		if strings.Contains(lower, marker) {
			return fail("no remote backend markers", "package.json contains forbidden backend/auth marker "+marker)
		}
	}

	return pass("no remote backend markers", "no obvious remote database or external auth packages declared")
}
