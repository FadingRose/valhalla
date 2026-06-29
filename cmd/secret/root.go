// Package secretcmd implements `valhalla secret ...` subcommands.
// It wraps an age-encrypted store of key/value secrets.
//
// Storage layout (inside repo):
//
//	secrets/
//	  llm/
//	    openrouter.age
//	    deepseek.age
//	  wallets/
//	    alice.age
//
// Each .age file is independently encrypted with every recipient listed in
// secrets/recipients.txt. To add a new device, append its public key and run
// `valhalla secret reencrypt` (TODO).
package secretcmd

import (
	"github.com/spf13/cobra"
)

var flagStoreDir string

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage age-encrypted secrets",
	Long: `secret subcommands operate on an encrypted directory inside the repo.

Secrets are stored as individual age-encrypted files keyed by path, e.g.
"llm/openrouter" lives at secrets/llm/openrouter.age.`,
	Example: `  # Store a value (will prompt)
  valhalla secret add wallets/new

  # Pipe a value instead
  echo "sk-..." | valhalla secret add llm/openrouter --stdin

  # Read a value
  valhalla secret get llm/openrouter

  # Copy to clipboard (macOS)
  valhalla secret get wallets/alice --copy

  # List all stored secrets
  valhalla secret list

  # Render all secrets to ~/.local/share/valhalla/secrets.fish
  valhalla secret inject

  # Re-encrypt after adding/removing a device
  valhalla secret reencrypt`,
}

func init() {
	secretCmd.PersistentFlags().StringVar(&flagStoreDir, "store", "secrets",
		"path to secrets directory (relative to repo root)")
}

// Cmd returns the secret subcommand so cmd root can attach it.
func Cmd() *cobra.Command {
	return secretCmd
}

// DevicesCmd returns the devices subcommand. Exported so the root package
// can attach it as a top-level command (alongside `secret`).
func DevicesCmd() *cobra.Command {
	return devicesCmd
}
