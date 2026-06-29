// Package linker abstracts the symlink lifecycle.
//
// Linker knows nothing about manifests or profiles — it operates purely on
// (source, target) pairs and reports what it did. This keeps the apply
// command small and makes linker trivially testable.
package linker

import (
	"fmt"
	"os"
	"path/filepath"
)

// Action describes what happened for a single link.
type Action string

const (
	ActionCreated  Action = "created"  // new symlink established
	ActionUpdated  Action = "updated"  // replaced an existing symlink
	ActionBackedUp Action = "backed up" // moved existing file/dir aside
	ActionSkipped  Action = "skipped"  // already correct, no-op
)

// Result is the outcome of linking one source to one target.
type Result struct {
	Source  string
	Target  string
	Action  Action
	Detail  string // human-readable extra info (backup path, etc.)
}

// Link ensures target points to source.
//
// Behavior:
//   - If target is already a symlink to source: skip (idempotent).
//   - If target is a symlink to somewhere else: replace.
//   - If target is a real file/dir: back up to <target>.bak, then link.
//   - Parent directories of target are created as needed.
//
// Source must exist (returns error otherwise).
func Link(source, target string) (*Result, error) {
	// Source must exist.
	if _, err := os.Lstat(source); err != nil {
		return nil, fmt.Errorf("source missing: %w", err)
	}

	// Ensure parent dir exists.
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir parent: %w", err)
	}

	// Case 1: target already a symlink.
	if existing, err := os.Readlink(target); err == nil {
		if existing == source {
			return &Result{Source: source, Target: target, Action: ActionSkipped}, nil
		}
		if err := os.Remove(target); err != nil {
			return nil, fmt.Errorf("remove old symlink: %w", err)
		}
		if err := os.Symlink(source, target); err != nil {
			return nil, err
		}
		return &Result{Source: source, Target: target, Action: ActionUpdated,
			Detail: fmt.Sprintf("was -> %s", existing)}, nil
	}

	// Case 2: target exists as real file/dir — back it up.
	if _, err := os.Lstat(target); err == nil {
		backup := target + ".bak"
		// Avoid clobbering an existing .bak.
		if _, err := os.Lstat(backup); err == nil {
			return nil, fmt.Errorf("target %s exists and %s already exists; resolve manually", target, backup)
		}
		if err := os.Rename(target, backup); err != nil {
			return nil, fmt.Errorf("backup: %w", err)
		}
		if err := os.Symlink(source, target); err != nil {
			return nil, err
		}
		return &Result{Source: source, Target: target, Action: ActionBackedUp,
			Detail: backup}, nil
	}

	// Case 3: nothing at target — fresh symlink.
	if err := os.Symlink(source, target); err != nil {
		return nil, err
	}
	return &Result{Source: source, Target: target, Action: ActionCreated}, nil
}

// Unlink removes a symlink if (and only if) it points at the expected source.
// Real files/dirs are never touched. Useful for `valhalla apply --dry-run`
// inverse or a future `valhalla revert` command.
func Unlink(target, expectedSource string) error {
	got, err := os.Readlink(target)
	if err != nil {
		return fmt.Errorf("not a symlink: %s", target)
	}
	if got != expectedSource {
		return fmt.Errorf("symlink points to %s, not %s; refusing to remove", got, expectedSource)
	}
	return os.Remove(target)
}
