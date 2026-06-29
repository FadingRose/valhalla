package secretcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

// defaultIdentityPath returns ~/.valhalla/identity.txt by convention.
// Override via VALHALLA_IDENTITY env if needed.
func defaultIdentityPath() string {
	if v := os.Getenv("VALHALLA_IDENTITY"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".valhalla", "identity.txt")
}

// resolveStoreDir returns the absolute secrets dir, anchored at the repo root.
func resolveStoreDir(repo string) string {
	if v := os.Getenv("VALHALLA_SECRETS_DIR"); v != "" {
		return v
	}
	if filepath.IsAbs(flagStoreDir) {
		return flagStoreDir
	}
	return filepath.Join(repo, flagStoreDir)
}

var secretInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a new age identity and recipients file",
	Long: `init creates a fresh age keypair for this device (pure-Go, no age-keygen binary needed).

Writes:
  ~/.valhalla/identity.txt   - private key (NEVER commit)
  <repo>/secrets/recipients.txt - public key (appended if exists)`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		return initKeys(secretsDir)
	},
}

func init() {
	secretCmd.AddCommand(secretInitCmd)
}

func initKeys(secretsDir string) error {
	identPath := defaultIdentityPath()
	if _, err := os.Stat(identPath); err == nil {
		return fmt.Errorf("identity already exists at %s (delete it first to regenerate)", identPath)
	}

	// Generate via pure-Go age library.
	pub, err := secretage.GenerateIdentity(identPath)
	if err != nil {
		return err
	}

	// Create secrets dir + append pubkey.
	if err := os.MkdirAll(secretsDir, 0o700); err != nil {
		return err
	}
	if err := secretage.AppendRecipient(secretsDir, pub, "device: "+hostname()); err != nil {
		return err
	}

	fmt.Printf("✓ generated identity: %s\n", identPath)
	fmt.Printf("✓ added recipient:    %s/%s\n", secretsDir, secretage.RecipientsFile)
	fmt.Println()
	fmt.Printf("public key: %s\n", pub)
	fmt.Println()
	fmt.Println("⚠️  Never commit the identity file. Add ~/.valhalla/ to .gitignore.")
	return nil
}
