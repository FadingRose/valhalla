package cmd

import (
	"fmt"

	"github.com/FadingRose/valhalla/internal/manifest"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show which declared links are missing or broken",
	Args:  cobra.NoArgs,
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(_ *cobra.Command, _ []string) error {
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

	missing, broken, ok := 0, 0, 0
	for _, l := range links {
		target, _ := manifest.ExpandTarget(l.Target)
		source := manifest.ResolveSource(repo, l.Source)
		status := classifyLink(target, source)
		switch status {
		case "ok":
			ok++
		case "missing":
			missing++
			fmt.Printf("  - %s (not linked)\n", l.Source)
		case "broken":
			broken++
			fmt.Printf("  ! %s (points elsewhere)\n", l.Source)
		}
	}
	fmt.Printf("\n%d ok, %d missing, %d broken\n", ok, missing, broken)
	return nil
}
