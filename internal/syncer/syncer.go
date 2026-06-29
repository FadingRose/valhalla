// Package syncer wraps git operations for the valhalla repo.
//
// Syncer is intentionally generic: it does add/commit/pull/push but knows
// nothing about manifests. The caller decides what "syncing" means at the
// command layer (commit-all vs selective, what message to use).
package syncer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Syncer operates on a fixed repo path. All git commands run with -C <path>.
type Syncer struct {
	RepoPath string
}

// New returns a Syncer rooted at repoPath. The path must be a git work tree.
func New(repoPath string) (*Syncer, error) {
	s := &Syncer{RepoPath: repoPath}
	if err := s.ensureGit(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Syncer) ensureGit() error {
	out, err := s.run("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return fmt.Errorf("not a git repo: %s (%v)", s.RepoPath, strings.TrimSpace(out))
	}
	return nil
}

// run executes a git command and returns combined output. Any non-zero exit
// produces an error that includes stderr.
func (s *Syncer) run(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", s.RepoPath}, args...)...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// HasChanges returns true if there are staged or unstaged changes.
func (s *Syncer) HasChanges() (bool, error) {
	out, err := s.run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// Branch reports the current branch name.
func (s *Syncer) Branch() (string, error) {
	out, err := s.run("rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(out), err
}

// Upstream returns the tracking ref (e.g. origin/main) or "" if unset.
func (s *Syncer) Upstream() (string, error) {
	out, err := s.run("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		return "", nil // no upstream configured — not an error
	}
	return strings.TrimSpace(out), nil
}

// AddAll stages all changes (equivalent to `git add -A`).
func (s *Syncer) AddAll() error {
	_, err := s.run("add", "-A")
	return err
}

// Add stages specific paths.
func (s *Syncer) Add(paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := s.run(args...)
	return err
}

// Commit creates a commit with the given message. Returns ErrNothingToCommit
// if the index is clean.
func (s *Syncer) Commit(message string) error {
	has, err := s.HasChanges()
	if err != nil {
		return err
	}
	if !has {
		return ErrNothingToCommit
	}
	_, err = s.run("commit", "-m", message)
	return err
}

// Pull runs `git pull --rebase`. Returns the combined output for reporting.
func (s *Syncer) Pull() (string, error) {
	return s.run("pull", "--rebase")
}

// Push runs `git push`. Returns the combined output for reporting.
func (s *Syncer) Push() (string, error) {
	return s.run("push")
}

// Summary generates a short, meaningful commit message by inspecting what
// changed. It groups changes by top-level dir to avoid dumping 20 filenames.
//
// Example output:
//   "update nvim, fish, manifest"
//   "update fish"
//   "update go"
func (s *Syncer) Summary() (string, error) {
	out, err := s.run("status", "--porcelain")
	if err != nil {
		return "update config", err
	}
	seen := map[string]bool{}
	var dirs []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		path := strings.TrimSpace(line[3:])
		// Strip quotes (git porcelain quotes paths with spaces).
		path = strings.Trim(path, `"`)
		top := path
		if i := strings.Index(path, "/"); i >= 0 {
			top = path[:i]
		}
		if !seen[top] {
			seen[top] = true
			dirs = append(dirs, top)
		}
	}
	if len(dirs) == 0 {
		return "update config", nil
	}
	if len(dirs) > 5 {
		return fmt.Sprintf("update %s and %d more", dirs[0], len(dirs)-1), nil
	}
	return "update " + strings.Join(dirs, ", "), nil
}

// ErrNothingToCommit is returned by Commit when the working tree is clean.
var ErrNothingToCommit = fmt.Errorf("nothing to commit")
