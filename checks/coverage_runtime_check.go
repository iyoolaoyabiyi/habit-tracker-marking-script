package checks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type coverageSummary struct {
	Total struct {
		Lines struct {
			Pct float64 `json:"pct"`
		} `json:"lines"`
	} `json:"total"`
}

func verifyCoverageSummary(root string, options RuntimeOptions) (string, error) {
	output, err := runCommandWithEnv(root, npmEnv(options), "npm", "exec", "--", "vitest", "run", "--coverage", "--coverage.reporter=json-summary")
	if err != nil {
		return output, err
	}

	summaryPath := filepath.Join(root, "coverage", "coverage-summary.json")
	content, err := os.ReadFile(summaryPath)
	if err != nil {
		return output, fmt.Errorf("could not read coverage summary at %s: %w", summaryPath, err)
	}

	summary, err := parseCoverageSummary(content)
	if err != nil {
		return output, fmt.Errorf("invalid coverage summary JSON: %w", err)
	}
	if summary.Total.Lines.Pct < 80 {
		return output, fmt.Errorf("src/lib line coverage is %.2f%%, below required 80%%", summary.Total.Lines.Pct)
	}

	return output, nil
}

func parseCoverageSummary(content []byte) (coverageSummary, error) {
	var summary coverageSummary
	err := json.Unmarshal(content, &summary)
	return summary, err
}
