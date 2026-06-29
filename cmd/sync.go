package cmd

import (
	"fmt"

	"github.com/FadingRose/valhalla/internal/syncer"
	"github.com/spf13/cobra"
)

var syncNoPull bool
var syncNoPush bool
var syncMsg string

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Commit local changes and synchronize with remote",
	Long: `sync is the successor to scripts/sync.sh.

It stages all changes, generates a meaningful commit message from the diff,
then pulls (rebase) and pushes. Flags allow skipping either side.`,
	Args: cobra.NoArgs,
	RunE: runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&syncNoPull, "no-pull", false, "skip git pull")
	syncCmd.Flags().BoolVar(&syncNoPush, "no-push", false, "skip git push")
	syncCmd.Flags().StringVar(&syncMsg, "message", "", "override commit message")
	rootCmd.AddCommand(syncCmd)
}

func runSync(_ *cobra.Command, _ []string) error {
	repo, err := resolveRepo()
	if err != nil {
		return err
	}
	s, err := syncer.New(repo)
	if err != nil {
		return err
	}

	has, err := s.HasChanges()
	if err != nil {
		return err
	}

	if has {
		if err := s.AddAll(); err != nil {
			return fmt.Errorf("add: %w", err)
		}
		msg := syncMsg
		if msg == "" {
			msg, err = s.Summary()
			if err != nil {
				msg = "update config"
			}
		}
		if err := s.Commit(msg); err != nil && err != syncer.ErrNothingToCommit {
			return fmt.Errorf("commit: %w", err)
		}
		fmt.Printf("committed: %s\n", msg)
	} else {
		fmt.Println("nothing to commit")
	}

	if !syncNoPull {
		out, err := s.Pull()
		if err != nil {
			return fmt.Errorf("pull: %s", out)
		}
		fmt.Println("pulled (rebase)")
	}

	if !syncNoPush {
		out, err := s.Push()
		if err != nil {
			return fmt.Errorf("push: %s", out)
		}
		fmt.Println("pushed")
	}

	fmt.Println("valhalla synced.")
	return nil
}
