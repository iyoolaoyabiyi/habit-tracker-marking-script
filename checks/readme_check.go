package checks

import "strings"

func checkReadme(root string) Result {
	var issues []string
	content, ok := readFileIssue(root, "README.md", &issues)
	if !ok {
		return failWithIssues("readme requirements", issues)
	}

	text := strings.ToLower(content)
	requiredPhrases := []string{
		"project overview",
		"setup instructions",
		"run instructions",
		"test instructions",
		"local persistence",
		"pwa support",
		"trade-offs",
		"test file mapping",
	}

	for _, phrase := range requiredPhrases {
		if !strings.Contains(text, phrase) {
			issues = append(issues, "README.md missing section containing "+phrase)
		}
	}

	if len(issues) > 0 {
		return failWithIssues("readme requirements", issues)
	}

	return pass("readme requirements", "required README sections found")
}
