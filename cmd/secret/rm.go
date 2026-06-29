package secretcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

var rmForce bool
var rmVerbose bool

var secretRmCmd = &cobra.Command{
	Use:     "rm <path>",
	Aliases: []string{"remove", "delete"},
	Short:   "Remove a secret from the store",
	Long: `rm deletes the encrypted file at the given path.

By default rm refuses to act if the path appears in manifest.toml's
[[secrets.inject]] list — removing an entry that's still declared would
cause the next ` + "`valhalla secret inject`" + ` to fail. Use --force to override.

The manifest.toml file itself is never edited automatically. After removing
a secret, you should also remove the corresponding [[secrets.inject]] entry.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		path := args[0]

		if !secretage.Exists(secretsDir, path) {
			if rmForce {
				return nil
			}
			return fmt.Errorf("no such secret: %s", path)
		}

		// Check manifest usage unless --force.
		if !rmForce {
			if used, by := inManifest(repo, path); used {
				return fmt.Errorf("secret %q is still declared in manifest.toml as var %s;\n"+
					"  remove that [[secrets.inject]] entry first, or use --force",
					path, strings.Join(by, ", "))
			}
		}

		full := filepath.Join(secretsDir, path+".age")
		if err := os.Remove(full); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
		if rmVerbose {
			fmt.Printf("removed: %s\n", path)
		} else {
			fmt.Printf("removed: %s\n", path)
		}

		// Clean up empty parent directories (but never the secrets root itself).
		cleanupEmptyDirs(secretsDir, full)
		return nil
	},
}

func init() {
	secretRmCmd.Flags().BoolVar(&rmForce, "force", false,
		"remove even if the path is still referenced in manifest.toml")
	secretCmd.AddCommand(secretRmCmd)
}

// inManifest reports whether the given secret path is declared in
// manifest.toml's [[secrets.inject]]. Returns the variable names that
// reference it.
func inManifest(repo, secretPath string) (bool, []string) {
	// We load manifest via the manifest package to avoid import cycle,
	// but since manifest doesn't expose AllEntries without a Profile,
	// we parse TOML directly here. Keep it simple.
	manifestPath := filepath.Join(repo, "manifest.toml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return false, nil
	}
	// Naive line-based scan — good enough for our simple manifest format.
	var referenced []string
	lines := strings.Split(string(data), "\n")
	inInject := false
	currentVar := ""
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[[secrets.inject]]") {
			inInject = true
			currentVar = ""
			continue
		}
		if strings.HasPrefix(trimmed, "[[") || strings.HasPrefix(trimmed, "[") {
			inInject = false
			continue
		}
		if !inInject {
			continue
		}
		if strings.HasPrefix(trimmed, "var ") || strings.HasPrefix(trimmed, "var=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				currentVar = strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
		}
		if strings.HasPrefix(trimmed, "path ") || strings.HasPrefix(trimmed, "path=") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				p := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				if p == secretPath && currentVar != "" {
					referenced = append(referenced, currentVar)
				}
			}
		}
	}
	return len(referenced) > 0, referenced
}

// cleanupEmptyDirs removes empty parent directories of removedPath up to
// (but not including) secretsRoot. Keeps the tree tidy.
func cleanupEmptyDirs(secretsRoot, removedPath string) {
	dir := filepath.Dir(removedPath)
	for dir != secretsRoot && dir != "/" && dir != "." {
		entries, err := os.ReadDir(dir)
		if err != nil || len(entries) > 0 {
			break
		}
		_ = os.Remove(dir)
		dir = filepath.Dir(dir)
	}
}
