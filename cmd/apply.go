package cmd

import (
	"fmt"

	"github.com/FadingRose/valhalla/internal/linker"
	"github.com/FadingRose/valhalla/internal/manifest"
	"github.com/spf13/cobra"
)

var applyDryRun bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Create symlinks declared in manifest.toml",
	Long: `apply reads manifest.toml and creates every declared symlink.
If a profile matches the current hostname (or --profile is given), its extra
links are added on top.

apply is idempotent: running it twice does nothing the second time if the
links are already correct.`,
	Example: `  # First-time setup on a new machine
  valhalla apply

  # Preview without changes
  valhalla apply --dry-run

  # Use a specific device profile
  valhalla apply --profile macbook`,
	Args: cobra.NoArgs,
	RunE: runApply,
}

func init() {
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false,
		"show what would happen without making changes")
	rootCmd.AddCommand(applyCmd)
}

func runApply(cmd *cobra.Command, _ []string) error {
	repo, err := resolveRepo()
	if err != nil {
		return err
	}
	m, err := manifest.Load(repo)
	if err != nil {
		return err
	}
	profile := m.SelectProfile(flagProfile)
	links := m.AllLinks(profile)

	fmt.Printf("repo:    %s\n", repo)
	if profile != nil {
		fmt.Printf("profile: %s (host=%s)\n", profile.Name, profile.MatchHost)
	}
	fmt.Printf("links:   %d\n\n", len(links))

	for _, l := range links {
		source := manifest.ResolveSource(repo, l.Source)
		target, err := manifest.ExpandTarget(l.Target)
		if err != nil {
			return fmt.Errorf("target %q: %w", l.Target, err)
		}
		if applyDryRun {
			fmt.Printf("  [dry-run] %s -> %s\n", target, source)
			continue
		}
		r, err := linker.Link(source, target)
		if err != nil {
			return fmt.Errorf("link %s: %w", l.Source, err)
		}
		printLinkResult(r)
	}
	return nil
}

func printLinkResult(r *linker.Result) {
	switch r.Action {
	case linker.ActionSkipped:
		fmt.Printf("  ✓ %s (already linked)\n", r.Target)
	case linker.ActionCreated:
		fmt.Printf("  + %s -> %s\n", r.Target, r.Source)
	case linker.ActionUpdated:
		fmt.Printf("  ~ %s -> %s (%s)\n", r.Target, r.Source, r.Detail)
	case linker.ActionBackedUp:
		fmt.Printf("  + %s -> %s (backup: %s)\n", r.Target, r.Source, r.Detail)
	}
}
