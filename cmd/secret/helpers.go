package secretcmd

import (
	"bytes"
	"os/exec"
	"strings"
)

// newBytesReader is a tiny helper to avoid importing "bytes" at the top
// while keeping the file self-contained for clipboard ops.
func newBytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }

// helperExec runs a command and returns trimmed combined output.
func helperExec(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
