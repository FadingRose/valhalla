// Device represents one authorized device in the recipients file.
// The file format is line-based:
//
//	# device: <name> (hostname: <host>, added <RFC3339 timestamp>)
//	age1...
//
// These helpers parse and mutate that file atomically.
package secretage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Device ties a friendly name to an age public key.
type Device struct {
	Name      string // user-chosen name, e.g. "mac-mini"
	Hostname  string // hostname at time of addition (informational)
	Pubkey    string // age1... recipient string
	AddedAt   string // RFC3339 timestamp
}

// LoadDevices reads recipients.txt and returns all registered devices.
// The comment line preceding each pubkey carries the metadata.
// Lines without a preceding device comment get Name = hostname (best effort).
func LoadDevices(secretsDir string) ([]Device, error) {
	path := filepath.Join(secretsDir, RecipientsFile)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var devices []Device
	var pendingMeta *Device // metadata from last # device: comment

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			pendingMeta = nil
			continue
		}
		if strings.HasPrefix(line, "#") {
			// Try to parse "# device: name (hostname: h, added TIME)".
			rest := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if d, ok := parseDeviceComment(rest); ok {
				meta := d
				pendingMeta = &meta
			}
			// Other comments are ignored but don't clear pendingMeta
			// (e.g. multi-line notes).
			continue
		}
		// Non-comment, non-empty: treat as pubkey line.
		d := Device{Pubkey: line}
		if pendingMeta != nil {
			d.Name = pendingMeta.Name
			d.Hostname = pendingMeta.Hostname
			d.AddedAt = pendingMeta.AddedAt
			pendingMeta = nil
		}
		devices = append(devices, d)
	}
	return devices, scanner.Err()
}

// parseDeviceComment parses "device: name (hostname: h, added TIME)".
// Returns ok=false if the line doesn't match the expected format.
func parseDeviceComment(s string) (Device, bool) {
	const prefix = "device:"
	if !strings.HasPrefix(s, prefix) {
		return Device{}, false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(s, prefix))

	var d Device
	// Extract parenthesized metadata if present.
	if i := strings.Index(rest, "("); i >= 0 {
		d.Name = strings.TrimSpace(rest[:i])
		if j := strings.LastIndex(rest, ")"); j > i {
			meta := rest[i+1 : j]
			for _, kv := range strings.Split(meta, ",") {
				kv = strings.TrimSpace(kv)
				if strings.HasPrefix(kv, "hostname:") {
					d.Hostname = strings.TrimSpace(strings.TrimPrefix(kv, "hostname:"))
				} else if strings.HasPrefix(kv, "added") {
					d.AddedAt = strings.TrimSpace(strings.TrimPrefix(kv, "added"))
				}
			}
		}
	} else {
		d.Name = rest
	}
	if d.Name == "" {
		return Device{}, false
	}
	return d, true
}

// FormatDeviceComment produces "# device: name (hostname: h, added TIME)".
func FormatDeviceComment(d Device) string {
	if d.AddedAt == "" {
		d.AddedAt = time.Now().UTC().Format(time.RFC3339)
	}
	return fmt.Sprintf("device: %s (hostname: %s, added %s)", d.Name, d.Hostname, d.AddedAt)
}

// AddDevice appends a new device to recipients.txt.
// Returns an error if the pubkey or name already exists.
func AddDevice(secretsDir string, d Device) error {
	existing, err := LoadDevices(secretsDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, e := range existing {
		if e.Pubkey == d.Pubkey {
			return fmt.Errorf("pubkey already registered as %q", e.Name)
		}
		if d.Name != "" && e.Name == d.Name {
			return fmt.Errorf("name %q already in use (pubkey %s...)", d.Name, e.Pubkey[:14])
		}
	}

	path := filepath.Join(secretsDir, RecipientsFile)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "# %s\n%s\n", FormatDeviceComment(d), d.Pubkey); err != nil {
		return err
	}
	return nil
}

// RemoveDevice deletes a device by name or by pubkey prefix.
// Returns the removed device and whether anything was removed.
// If pubkeyOnly is false, matches by name first, then pubkey prefix.
func RemoveDevice(secretsDir, identifier string) (Device, bool, error) {
	path := filepath.Join(secretsDir, RecipientsFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return Device{}, false, err
	}

	var out []string
	var removed Device
	didRemove := false

	lines := strings.Split(string(data), "\n")
	skipNextPubkey := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Detect "# device: ..." comment lines.
		if strings.HasPrefix(trimmed, "#") {
			rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			if d, ok := parseDeviceComment(rest); ok {
				// Lookahead: is the next non-empty line this device's pubkey?
				nextPub := ""
				for j := i + 1; j < len(lines); j++ {
					t := strings.TrimSpace(lines[j])
					if t == "" {
						continue
					}
					nextPub = t
					break
				}
				matchName := d.Name == identifier
				matchKey := strings.HasPrefix(nextPub, identifier)
				if matchName || matchKey {
					removed = d
					removed.Pubkey = nextPub
					didRemove = true
					// Skip this comment line AND the following pubkey line.
					skipNextPubkey = true
					continue
				}
			}
		}

		if skipNextPubkey && strings.HasPrefix(trimmed, "age1") {
			skipNextPubkey = false
			continue
		}
		skipNextPubkey = false

		// Preserve the line (drop trailing empties handled at end).
		out = append(out, line)
	}

	if !didRemove {
		return Device{}, false, nil
	}

	// Re-serialize, trimming any leading/trailing blank lines.
	result := strings.Join(out, "\n")
	result = strings.TrimRight(result, "\n") + "\n"
	return removed, true, os.WriteFile(path, []byte(result), 0o644)
}
