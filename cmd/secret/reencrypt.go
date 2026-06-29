package secretcmd

import (
	"fmt"
	"os"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

var secretReencryptCmd = &cobra.Command{
	Use:   "reencrypt",
	Short: "Re-encrypt all secrets for the current set of devices",
	Long: `reencrypt walks every .age file in the store, decrypts it with the
local identity, and re-encrypts it using all pubkeys currently listed in
secrets/recipients.txt.

Use this after:
  - ` + "`valhalla devices add`" + ` (so the new device can decrypt existing secrets)
  - ` + "`valhalla devices remove`" + ` (so the removed device can no longer decrypt)
  - manual edits to recipients.txt

Files that the local identity cannot decrypt are skipped (with a warning).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		identPath := defaultIdentityPath()
		if _, err := os.Stat(identPath); err != nil {
			return fmt.Errorf("identity missing at %s: %w", identPath, err)
		}

		// Show current recipients before re-encrypting.
		devs, _ := secretage.LoadDevices(secretsDir)
		fmt.Printf("Re-encrypting for %d device(s):\n", len(devs))
		for _, d := range devs {
			fmt.Printf("  • %s\n", d.Name)
		}
		fmt.Println()

		count, err := secretage.ReencryptAll(secretsDir, identPath)
		if err != nil {
			return err
		}
		fmt.Printf("\n✓ re-encrypted %d secret(s)\n", count)
		fmt.Println()
		fmt.Println("Next: push the updated ciphertexts:")
		fmt.Println("  valhalla sync")
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretReencryptCmd)
}
