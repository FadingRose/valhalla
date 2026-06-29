package secretcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

var (
	initName     string
	initNoAddDev bool
)

var secretInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a new age identity and register this device",
	Long: `init creates a fresh age keypair for this device (pure-Go, no age-keygen binary needed).

Writes:
  ~/.valhalla/identity.txt       - private key (NEVER commit, mode 0600)
  <repo>/secrets/recipients.txt  - public key appended (with device name)

After init on a new device, you MUST run ` + "`valhalla secret reencrypt`" + ` on
an existing device (one that can already decrypt the store) so all secrets
get encrypted for the newcomer too. Then push and pull on the new device.`,
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
	secretInitCmd.Flags().StringVar(&initName, "name", "",
		"friendly device name (default: short hostname)")
	secretInitCmd.Flags().BoolVar(&initNoAddDev, "no-register", false,
		"generate identity only, do not add to recipients.txt")
	secretCmd.AddCommand(secretInitCmd)
}

func initKeys(secretsDir string) error {
	identPath := defaultIdentityPath()
	if _, err := os.Stat(identPath); err == nil {
		return fmt.Errorf("identity already exists at %s\n"+
			"  (to register an existing device, use `valhalla devices add`)", identPath)
	}

	// Generate via pure-Go age library.
	pub, err := secretage.GenerateIdentity(identPath)
	if err != nil {
		return err
	}

	// Determine device name.
	name := initName
	if name == "" {
		name = hostname()
	}

	fmt.Printf("✓ generated identity: %s\n", identPath)
	fmt.Printf("  public key: %s\n", pub)
	fmt.Println()

	if initNoAddDev {
		fmt.Println("(skipped recipients.txt registration --no-register)")
		fmt.Println()
		fmt.Println("⚠️  Never commit the identity file. Add ~/.valhalla/ to .gitignore.")
		return nil
	}

	// Create secrets dir + add device entry.
	if err := os.MkdirAll(secretsDir, 0o700); err != nil {
		return err
	}
	d := secretage.Device{
		Name:     name,
		Hostname: hostname(),
		Pubkey:   pub,
	}
	if err := secretage.AddDevice(secretsDir, d); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  could not add to recipients.txt: %v\n", err)
		fmt.Fprintln(os.Stderr, "    (you can add it manually with `valhalla devices add`)")
	} else {
		fmt.Printf("✓ registered device: %s\n", name)
		fmt.Printf("  in: %s/%s\n", secretsDir, secretage.RecipientsFile)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Run this on a device that can already decrypt secrets:")
		fmt.Println("       valhalla secret reencrypt")
		fmt.Println("       valhalla sync")
		fmt.Println("  2. Back here, pull and decrypt:")
		fmt.Println("       valhalla sync")
		fmt.Println("       valhalla secret inject")
	}
	fmt.Println()
	fmt.Println("⚠️  Never commit the identity file. Add ~/.valhalla/ to .gitignore.")
	return nil
}

// hostname for tagging recipients.
func hostname() string {
	host, _ := os.Hostname()
	if i := strings.IndexByte(host, '.'); i > 0 {
		host = host[:i]
	}
	return host
}
