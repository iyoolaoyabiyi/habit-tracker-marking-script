package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkStorageKeys(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "src/lib/constants.ts"))
	if err != nil {
		return fail("storage keys", err.Error())
	}

	required := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	for _, key := range required {
		if !strings.Contains(string(content), key) {
			return fail("storage keys", "missing key "+key+" in src/lib/constants.ts")
		}
	}

	return pass("storage keys", "required localStorage keys found")
}
