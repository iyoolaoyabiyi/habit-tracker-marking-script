package checks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func checkPackageScripts(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return fail("package scripts", err.Error())
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fail("package scripts", "invalid package.json: "+err.Error())
	}

	required := []string{"dev", "build", "start", "test:unit", "test:integration", "test:e2e", "test"}
	var missing []string
	for _, name := range required {
		if _, ok := pkg.Scripts[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fail("package scripts", "missing script(s): "+strings.Join(missing, ", "))
	}

	return pass("package scripts", "all required script names present")
}
