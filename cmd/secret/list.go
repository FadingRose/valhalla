package secretcmd

import (
	"fmt"
	"os"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

	var secretListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all stored secret paths",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			repo, err := resolveRepo()
			if err != nil {
				return err
			}
			secretsDir := resolveStoreDir(repo)
			// Friendly handling of missing store.
			if info, err := os.Stat(secretsDir); err != nil || !info.IsDir() {
				fmt.Println("(no secrets stored — run `valhalla secret init` first)")
				return nil
			}
			paths, err := secretage.ListAll(secretsDir)
			if err != nil {
				return err
			}
			if len(paths) == 0 {
				fmt.Println("(no secrets stored)")
				return nil
			}
			for _, p := range paths {
				fmt.Println(p)
			}
			return nil
		},
	}

func init() {
	secretCmd.AddCommand(secretListCmd)
}
