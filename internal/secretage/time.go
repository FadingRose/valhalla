package secretage

import "time"

// nowString returns an RFC3339 timestamp for identity file headers.
// Isolated in its own file so the core secretage.go stays import-light.
func nowString() string {
	return time.Now().UTC().Format(time.RFC3339)
}
