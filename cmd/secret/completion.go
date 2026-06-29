// Package secretcmd — dynamic completion functions.
//
// These power `<TAB>` completion for paths and device names. Cobra calls
// them via `__complete` under the hood. Each returns a list of candidates
// plus a directive (NoFileComp, Default, etc.).
package secretcmd

import (
	"strings"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

// completeSecretPath returns all stored secret paths. Used for:
//   valhalla secret get <TAB>
//   valhalla secret rm <TAB>
func completeSecretPath(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	repo, err := resolveRepo()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	secretsDir := resolveStoreDir(repo)
	paths, err := secretage.ListAll(secretsDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var matched []string
	for _, p := range paths {
		if strings.HasPrefix(p, toComplete) {
			matched = append(matched, p)
		}
	}
	return matched, cobra.ShellCompDirectiveNoFileComp
}

// completeDeviceName returns all registered device names. Used for:
//   valhalla devices remove <TAB>
func completeDeviceName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	repo, err := resolveRepo()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	secretsDir := resolveStoreDir(repo)
	devs, err := secretage.LoadDevices(secretsDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	var matched []string
	for _, d := range devs {
		if strings.HasPrefix(d.Name, toComplete) {
			matched = append(matched, d.Name)
		}
	}
	return matched, cobra.ShellCompDirectiveNoFileComp
}

// completeShellFormat returns the supported --shell values.
// Used for:
//   valhalla secret inject --shell <TAB>
func completeShellFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	shells := []string{"fish", "posix", "env", "json"}
	var matched []string
	for _, s := range shells {
		if strings.HasPrefix(s, toComplete) {
			matched = append(matched, s)
		}
	}
	return matched, cobra.ShellCompDirectiveNoFileComp
}
