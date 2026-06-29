// Package secretcmd — devices subcommand.
//
// `valhalla devices list/add/remove` manages the recipients.txt file.
// Adding or removing a device does NOT re-encrypt existing secrets —
// that's done by `valhalla secret reencrypt`. The typical flow is:
//
//	valhalla devices add ...    (new device)
//	valhalla secret reencrypt   (re-encrypt all secrets for new recipients)
//	valhalla sync               (push)
package secretcmd

import (
	"fmt"
	"os"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage devices authorized to decrypt secrets",
	Long: `devices operates on secrets/recipients.txt.

List, add, or remove devices (age public keys). Adding a new device does not
automatically re-encrypt existing secrets — run ` + "`valhalla secret reencrypt`" + `
after adding/removing devices.`,
}

// devicesListCmd — `valhalla devices list`
var devicesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all registered devices",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		devs, err := secretage.LoadDevices(secretsDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("(no recipients.txt — run `valhalla secret init` first)")
				return nil
			}
			return err
		}
		if len(devs) == 0 {
			fmt.Println("(no devices registered)")
			return nil
		}
		fmt.Printf("%-20s %-30s %-12s %s\n", "NAME", "PUBKEY", "ADDED", "HOSTNAME")
		for _, d := range devs {
			pub := d.Pubkey
			if len(pub) > 28 {
				pub = pub[:28] + "…"
			}
			added := d.AddedAt
			if len(added) >= 10 {
				added = added[:10] // date only
			}
			fmt.Printf("%-20s %-30s %-12s %s\n", d.Name, pub, added, d.Hostname)
		}
		fmt.Printf("\n%d device(s)\n", len(devs))
		return nil
	},
}

// devicesAddCmd — `valhalla devices add <pubkey> [--name NAME]`
var (
	devicesAddName     string
	devicesAddHostname string
)

var devicesAddCmd = &cobra.Command{
	Use:   "add <pubkey>",
	Short: "Register a device's public key",
	Long: `add appends a pubkey to recipients.txt.

This is used when a new device has generated its identity (via
` + "`valhalla secret init`" + `) and shared its pubkey with you. After adding,
run ` + "`valhalla secret reencrypt`" + ` so all existing secrets are encrypted
for the new device too.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		name := devicesAddName
		if name == "" {
			name = "unnamed-device"
		}
		d := secretage.Device{
			Name:     name,
			Hostname: devicesAddHostname,
			Pubkey:   args[0],
		}
		if err := secretage.AddDevice(secretsDir, d); err != nil {
			return err
		}
		fmt.Printf("✓ added device: %s (%s)\n", d.Name, shortKey(d.Pubkey))
		fmt.Println()
		fmt.Println("Next: re-encrypt all secrets so this device can decrypt them:")
		fmt.Println("  valhalla secret reencrypt")
		return nil
	},
}

// devicesRemoveCmd — `valhalla devices remove <name|pubkey>`
var devicesRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a device and forbid it from decrypting new secrets",
	Long: `remove deletes a device entry from recipients.txt.

After removal, this device's identity can no longer decrypt NEWLY encrypted
secrets (after ` + "`valhalla secret reencrypt`" + `). However, the device can
still decrypt the OLD git history — for a full eviction, rotate all secret
values themselves.

The argument can be a device name or a pubkey prefix.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		removed, ok, err := secretage.RemoveDevice(secretsDir, args[0])
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no device matching %q", args[0])
		}
		fmt.Printf("✓ removed device: %s (%s)\n", removed.Name, shortKey(removed.Pubkey))
		fmt.Println()
		fmt.Println("Next: re-encrypt all secrets WITHOUT this device's pubkey:")
		fmt.Println("  valhalla secret reencrypt")
		fmt.Println()
		fmt.Println("⚠️  This only affects future commits. The removed device can still")
		fmt.Println("   decrypt secrets in older git history. To fully evict, rotate")
		fmt.Println("   secret values with `valhalla secret add` for each entry.")
		return nil
	},
}

func init() {
	devicesAddCmd.Flags().StringVar(&devicesAddName, "name", "",
		"friendly device name")
	devicesAddCmd.Flags().StringVar(&devicesAddHostname, "hostname", "",
		"hostname (informational)")

	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesAddCmd)
	devicesCmd.AddCommand(devicesRemoveCmd)
}

// shortKey returns the first 14 chars of a pubkey plus ellipsis, or
// the whole key if shorter. Never panics on short/empty input.
func shortKey(k string) string {
	if len(k) <= 14 {
		return k
	}
	return k[:14] + "…"
}
