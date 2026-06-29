package secretcmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

var secretGetCopy bool

var secretGetCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Print a decrypted secret to stdout",
	Long: `get decrypts and prints the value at the given path.

Use --copy to copy to clipboard via pbcopy / xclip / wl-copy instead.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		plaintext, err := secretage.DecryptFromFile(secretsDir, args[0], defaultIdentityPath())
		if err != nil {
			return err
		}
		if secretGetCopy {
			return copyToClipboard(plaintext)
		}
		fmt.Println(string(plaintext))
		return nil
	},
}

func init() {
	secretGetCmd.Flags().BoolVar(&secretGetCopy, "copy", false, "copy to clipboard instead of printing")
	// Dynamic completion: `valhalla secret get <TAB>` lists stored secrets.
	secretGetCmd.ValidArgsFunction = completeSecretPath
	secretCmd.AddCommand(secretGetCmd)
}

func copyToClipboard(data []byte) error {
	var binName string
	for _, name := range []string{"pbcopy", "wl-copy", "xclip"} {
		if _, err := exec.LookPath(name); err == nil {
			binName = name
			break
		}
	}
	if binName == "" {
		return fmt.Errorf("no clipboard binary found (pbcopy/wl-copy/xclip)")
	}
	cmd := exec.Command(binName)
	cmd.Stdin = newBytesReader(data)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard: %w", err)
	}
	fmt.Fprintln(os.Stderr, "copied to clipboard.")
	return nil
}
