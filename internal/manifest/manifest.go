// Package manifest defines the declarative structure that drives valhalla.
// It is the single source of truth for what valhalla manages.
//
// The manifest answers two questions:
//   - "What should be linked where?" (links)
//   - "Which secrets should be injected and where?" (secrets)
//
// Optional device-specific overrides live under [profiles.<name>].
package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Manifest is the top-level config structure.
//
// Example manifest.toml:
//
//	[[links]]
//	source = "nvim"
//	target = "~/.config/nvim"
//
//	[secrets]
//	output = "~/.config/fish/conf.d/secrets.fish"
//	[[secrets.inject]]
//	var = "OPENROUTER_API_KEY"
//	path = "llm/openrouter"
type Manifest struct {
	Links    []Link    `toml:"links"`
	Secrets  Secrets   `toml:"secrets"`
	Profiles []Profile `toml:"profiles"`
}

// Link declares one symlink to be created by `valhalla apply`.
// Source is relative to repo root; Target may contain ~ and is expanded.
type Link struct {
	Source string `toml:"source"`
	Target string `toml:"target"`
}

// Secrets declares the injection target and the variables to materialize.
type Secrets struct {
	Output string  `toml:"output"`
	Inject []Entry `toml:"inject"`
}

// Entry binds a shell variable name to a secret path in the encrypted store.
type Entry struct {
	Var  string `toml:"var"`
	Path string `toml:"path"`
}

// Profile is a named bundle of extra links/secrets for a specific host.
// Selecting a profile via --profile merges its fields on top of the base.
type Profile struct {
	Name       string  `toml "name"`
	MatchHost  string  `toml:"match_host"`
	ExtraLinks []Link  `toml:"extra_links"`
	ExtraEntry []Entry `toml:"extra_inject"`
}

// Load reads and parses manifest.toml from the given repo root.
// It also expands ~ in targets and resolves source to absolute paths.
func Load(repoRoot string) (*Manifest, error) {
	path := filepath.Join(repoRoot, "manifest.toml")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("manifest not found at %s: %w", path, err)
	}

	var m Manifest
	if _, err := toml.DecodeFile(path, &m); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Manifest) validate() error {
	if len(m.Links) == 0 && len(m.Secrets.Inject) == 0 {
		return fmt.Errorf("manifest is empty: define at least one [[links]] or [[secrets.inject]]")
	}
	for i, l := range m.Links {
		if l.Source == "" || l.Target == "" {
			return fmt.Errorf("links[%d]: source and target are required", i)
		}
	}
	return nil
}

// ExpandTarget returns the absolute target path with ~ expanded.
func ExpandTarget(target string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("empty target")
	}
	if target[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, target[1:]), nil
	}
	return target, nil
}

// ResolveSource returns absolute path of a link source within repoRoot.
func ResolveSource(repoRoot, source string) string {
	return filepath.Join(repoRoot, source)
}

// Hostname returns the current machine hostname for profile matching.
func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		return runtime.GOOS
	}
	return host
}

// SelectProfile returns a merged manifest: base + the profile that matches
// either matchHost or the given --profile flag. Nil if no match.
func (m *Manifest) SelectProfile(name string) *Profile {
	// Explicit --profile wins.
	if name != "" {
		for i := range m.Profiles {
			if m.Profiles[i].Name == name {
				return &m.Profiles[i]
			}
		}
	}
	// Auto-match by hostname.
	host := Hostname()
	for i := range m.Profiles {
		if m.Profiles[i].MatchHost == host {
			return &m.Profiles[i]
		}
	}
	return nil
}

// AllLinks returns base links + profile extras (if any).
func (m *Manifest) AllLinks(profile *Profile) []Link {
	out := make([]Link, len(m.Links))
	copy(out, m.Links)
	if profile != nil {
		out = append(out, profile.ExtraLinks...)
	}
	return out
}

// AllEntries returns base secret entries + profile extras.
func (m *Manifest) AllEntries(profile *Profile) []Entry {
	out := make([]Entry, len(m.Secrets.Inject))
	copy(out, m.Secrets.Inject)
	if profile != nil {
		out = append(out, profile.ExtraEntry...)
	}
	return out
}
