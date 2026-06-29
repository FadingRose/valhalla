package cmd

import (
	"os"
	"path/filepath"
)

// classifyLink inspects target and reports its relationship to source.
// Returns one of: "ok" (correct symlink), "missing" (target absent),
// "broken" (target exists but doesn't point to source).
func classifyLink(target, source string) string {
	got, err := os.Readlink(target)
	if err != nil {
		return "missing"
	}
	abs, _ := filepath.Abs(source)
	if got == source || got == abs {
		return "ok"
	}
	return "broken"
}
