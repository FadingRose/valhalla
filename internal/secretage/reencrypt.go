// ReencryptAll decrypts every .age file in the store and re-encrypts it
// using the current recipients list. Used after adding/removing devices.
package secretage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReencryptAll walks the store, decrypts each entry with identityPath,
// and re-encrypts with the current recipients loaded from recipients.txt.
//
// Returns the number of files re-encrypted. If any file fails, returns
// immediately with the error.
//
// If identityPath cannot decrypt a file (e.g. it belongs to a device not
// yet in recipients), that file is skipped with a warning, not an error.
func ReencryptAll(secretsDir, identityPath string) (int, error) {
	recips, err := LoadRecipients(secretsDir)
	if err != nil {
		return 0, fmt.Errorf("load recipients: %w", err)
	}

	count := 0
	err = filepath.Walk(secretsDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(secretsDir, p)
		if rel == RecipientsFile {
			return nil
		}
		if !strings.HasSuffix(rel, ".age") {
			return nil
		}

		// Decrypt with old identity.
		plaintext, err := DecryptFromFile(secretsDir,
			strings.TrimSuffix(rel, ".age"), identityPath)
		if err != nil {
			// Best-effort: skip files we can't decrypt (they might belong
			// to a device that hasn't been added to recipients yet).
			fmt.Fprintf(os.Stderr, "  ! skip %s (cannot decrypt with this identity)\n", rel)
			return nil
		}

		// Re-encrypt with new recipients.
		if err := EncryptToFile(secretsDir,
			strings.TrimSuffix(rel, ".age"), plaintext, recips); err != nil {
			return fmt.Errorf("re-encrypt %s: %w", rel, err)
		}
		count++
		return nil
	})
	return count, err
}
