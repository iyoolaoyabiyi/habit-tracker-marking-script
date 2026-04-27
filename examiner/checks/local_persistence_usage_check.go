package checks

import (
	"os"
	"path/filepath"
	"strings"
)

func checkLocalPersistenceUsage(root string) Result {
	storageContent, err := os.ReadFile(filepath.Join(root, "src/lib/storage.ts"))
	if err != nil {
		return fail("local persistence usage", err.Error())
	}
	constantsContent, err := os.ReadFile(filepath.Join(root, "src/lib/constants.ts"))
	if err != nil {
		return fail("local persistence usage", err.Error())
	}

	storageText := string(storageContent)
	storageRequired := []string{
		"window.localStorage.getItem",
		"window.localStorage.setItem",
	}
	for _, marker := range storageRequired {
		if !strings.Contains(storageText, marker) {
			return fail("local persistence usage", "src/lib/storage.ts missing "+marker)
		}
	}

	constantsText := string(constantsContent)
	keyRequired := []string{
		"habit-tracker-users",
		"habit-tracker-session",
		"habit-tracker-habits",
	}
	for _, marker := range keyRequired {
		if !strings.Contains(constantsText, marker) {
			return fail("local persistence usage", "src/lib/constants.ts missing "+marker)
		}
	}

	return pass("local persistence usage", "localStorage persistence is clearly implemented")
}
