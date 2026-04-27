package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func checkSplashTiming(root string) Result {
	content, err := os.ReadFile(filepath.Join(root, "src/app/page.tsx"))
	if err != nil {
		return fail("splash timing", err.Error())
	}
	matches := regexp.MustCompile(`SPLASH_DELAY_MS\s*=\s*(\d+)`).FindStringSubmatch(string(content))
	if len(matches) != 2 {
		return fail("splash timing", "could not find splash delay constant in src/app/page.tsx")
	}
	var delay int
	fmt.Sscanf(matches[1], "%d", &delay)
	if delay < 800 || delay > 2000 {
		return fail("splash timing", fmt.Sprintf("splash delay %dms is outside the required 800-2000ms range", delay))
	}

	return pass("splash timing", fmt.Sprintf("splash delay is %dms", delay))
}
