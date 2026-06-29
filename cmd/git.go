package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// runGit runs a git command in the given dir and returns trimmed stdout.
func runGit(dir string, args ...string) (string, error) {
	full := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", full...)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), string(ee.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
