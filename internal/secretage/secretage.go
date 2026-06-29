// Package secretage wraps filippo.io/age for our use case.
package secretage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

// RecipientsFile is the canonical name for the public-key list.
const RecipientsFile = "recipients.txt"

// IdentityHeader / PubkeyHeader are the file formats age-keygen uses.
const (
	identityFormat = "# created: %s\n# public key: %s\nAGE-SECRET-KEY-%s\n"
)

// ─── Recipients ────────────────────────────────────────────────────────────

// LoadRecipients reads and parses secrets/recipients.txt.
// Each line is an age public key (age1...); lines starting with # are ignored.
func LoadRecipients(secretsDir string) ([]age.Recipient, error) {
	path := filepath.Join(secretsDir, RecipientsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read recipients: %w (run `valhalla secret init` first)", err)
	}
	var recips []age.Recipient
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r, err := age.ParseX25519Recipient(line)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient %q: %w", line, err)
		}
		recips = append(recips, r)
	}
	if len(recips) == 0 {
		return nil, fmt.Errorf("no recipients in %s", path)
	}
	return recips, nil
}

// AppendRecipient adds a public key line with an optional comment.
func AppendRecipient(secretsDir, pubkey, comment string) error {
	path := filepath.Join(secretsDir, RecipientsFile)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if comment != "" {
		fmt.Fprintf(f, "# %s\n", comment)
	}
	fmt.Fprintf(f, "%s\n", pubkey)
	return nil
}

// ─── Identity generation (pure Go, no age-keygen binary) ───────────────────

// GenerateIdentity creates a new X25519 identity and writes it to identPath
// in the standard age-keygen format. Returns the public key string.
func GenerateIdentity(identPath string) (pubkey string, err error) {
	sc, err := age.GenerateX25519Identity()
	if err != nil {
		return "", fmt.Errorf("generate identity: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(identPath), 0o700); err != nil {
		return "", err
	}
	pub := sc.Recipient().String()
	// Serialize private key in age-keygen compatible format.
	// AGE-SECRET-KEY-1...  (base32-no-pad, bech32-style)
	priv := sc.String()
	content := fmt.Sprintf("# created: %s\n# public key: %s\n%s\n",
		timestamp(), pub, priv)
	if err := os.WriteFile(identPath, []byte(content), 0o600); err != nil {
		return "", err
	}
	return pub, nil
}

// timestamp is a helper to avoid pulling in time at package init.
func timestamp() string { return nowString() }

// ─── Encrypt / Decrypt ─────────────────────────────────────────────────────

// Encrypt encrypts plaintext for the given recipients.
func Encrypt(plaintext []byte, recips []age.Recipient) ([]byte, error) {
	var buf bytes.Buffer
	aw, err := age.Encrypt(&buf, recips...)
	if err != nil {
		return nil, err
	}
	if _, err := aw.Write(plaintext); err != nil {
		return nil, err
	}
	if err := aw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decrypt decrypts ciphertext using an identity loaded from the given path.
func Decrypt(ciphertext []byte, identityPath string) ([]byte, error) {
	ids, err := loadIdentities(identityPath)
	if err != nil {
		return nil, err
	}
	ar, err := age.Decrypt(bytes.NewReader(ciphertext), ids...)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return io.ReadAll(ar)
}

func loadIdentities(path string) ([]age.Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read identity %s: %w", path, err)
	}
	return age.ParseIdentities(bytes.NewReader(data))
}

// ─── File-level helpers ────────────────────────────────────────────────────

// EncryptToFile writes an encrypted value to secretsDir/<path>.age.
func EncryptToFile(secretsDir, secretPath string, plaintext []byte, recips []age.Recipient) error {
	full := filepath.Join(secretsDir, secretPath+".age")
	if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
		return err
	}
	ciphertext, err := Encrypt(plaintext, recips)
	if err != nil {
		return err
	}
	return os.WriteFile(full, ciphertext, 0o600)
}

// DecryptFromFile reads secretsDir/<path>.age and returns decrypted bytes.
func DecryptFromFile(secretsDir, secretPath, identityPath string) ([]byte, error) {
	full := filepath.Join(secretsDir, secretPath+".age")
	ciphertext, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", full, err)
	}
	return Decrypt(ciphertext, identityPath)
}

// ListAll walks secretsDir and returns all stored secret paths (without .age).
func ListAll(secretsDir string) ([]string, error) {
	var paths []string
	err := filepath.Walk(secretsDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(secretsDir, p)
		if err != nil {
			return err
		}
		if rel == RecipientsFile {
			return nil
		}
		if !strings.HasSuffix(rel, ".age") {
			return nil
		}
		clean := strings.TrimSuffix(rel, ".age")
		paths = append(paths, filepath.ToSlash(clean))
		return nil
	})
	return paths, err
}

// Exists reports whether a given secret path is present.
func Exists(secretsDir, secretPath string) bool {
	full := filepath.Join(secretsDir, secretPath+".age")
	_, err := os.Stat(full)
	return err == nil
}
