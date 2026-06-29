// Package cmd implements the valhalla CLI surface.
//
// Commands are organized by file; each file registers its subcommand on init.
// root.go owns global flags (--repo, --profile) and common helpers.
package cmd

import (
	"fmt"
	"os"

	secretcmd "github.com/FadingRose/valhalla/cmd/secret"
	"github.com/spf13/cobra"
)

var (
	flagRepo    string
	flagProfile string
)

// rootCmd is the top-level entry. Run via `valhalla` or `go run .`.
var rootCmd = &cobra.Command{
	Use:   "valhalla",
	Short: "Unified config + secrets manager",
	Long: `valhalla manages dotfiles via symlinks and secrets via age-encrypted files.

It replaces scripts/install.sh and scripts/sync.sh with a single binary that
reads manifest.toml for declarative configuration.`,
	Example: `  # Apply all config symlinks on a new machine
  valhalla apply

  # Sync local changes to git
  valhalla sync

  # Look up a secret value
  valhalla secret get llm/openrouter

  # Rotate a key, then reload
  valhalla secret add llm/openrouter
  valhalla secret inject

  # See full docs: https://github.com/FadingRose/valhalla`,
	SilenceUsage: true,
}

// Execute is the single entry point called from main.
func Execute() error {
	rootCmd.AddCommand(secretcmd.Cmd())
	rootCmd.AddCommand(secretcmd.DevicesCmd())
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagRepo, "repo", "",
		"path to valhalla repo (default: auto-detect via git rev-parse)")
	rootCmd.PersistentFlags().StringVar(&flagProfile, "profile", "",
		"name of device profile to activate (default: auto-match hostname)")
}

// resolveRepo returns the repo path from --repo, or falls back to git top-level
// of CWD, or errors.
func resolveRepo() (string, error) {
	if flagRepo != "" {
		if _, err := os.Stat(flagRepo); err != nil {
			return "", fmt.Errorf("--repo: %w", err)
		}
		return flagRepo, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	top, err := gitTopLevel(cwd)
	if err != nil {
		return "", fmt.Errorf("not inside a git repo and --repo not set: %w", err)
	}
	return top, nil
}

func gitTopLevel(dir string) (string, error) {
	return runGit(dir, "rev-parse", "--show-toplevel")
}
