package secretcmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/FadingRose/valhalla/internal/secretage"
	"github.com/spf13/cobra"
)

var secretAddFromStdin bool

var secretAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add or update a secret",
	Long: `add stores a value at the given path (e.g. llm/openrouter).

Value is read interactively unless --stdin is given (then piped).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := resolveRepo()
		if err != nil {
			return err
		}
		secretsDir := resolveStoreDir(repo)
		recips, err := secretage.LoadRecipients(secretsDir)
		if err != nil {
			return err
		}
		value, err := readValue()
		if err != nil {
			return err
		}
		if err := secretage.EncryptToFile(secretsDir, args[0], value, recips); err != nil {
			return err
		}
		fmt.Printf("stored: %s (encrypted, %d bytes)\n", args[0], len(value))
		return nil
	},
}

func init() {
	secretAddCmd.Flags().BoolVar(&secretAddFromStdin, "stdin", false,
		"read value from stdin instead of prompting")
	secretCmd.AddCommand(secretAddCmd)
}

func readValue() ([]byte, error) {
	if secretAddFromStdin {
		return readAll(os.Stdin)
	}
	// Prompt on stderr so stdout stays clean for pipes.
	fmt.Fprint(os.Stderr, "value: ")
	// Use `read -s` via stty to hide input, mimicking password entry.
	out, err := hiddenRead()
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr)
	return out, nil
}

func hiddenRead() ([]byte, error) {
	// Disable echo using stty; this is POSIX and portable across macOS/Linux.
	tty, err := os.Open("/dev/tty")
	if err != nil {
		// Fallback: read from stdin without echo suppression.
		return readAll(os.Stdin)
	}
	defer tty.Close()
	_ = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	return readUntilNewline(tty)
}

func readAll(r *os.File) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		buf = append(buf, tmp[:n]...)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if n == 0 {
			break
		}
	}
	// Trim trailing newline for stdin-fed values.
	return trimTrailingNL(buf), nil
}

func readUntilNewline(f *os.File) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	one := make([]byte, 1)
	for {
		n, err := f.Read(one)
		if err != nil || n == 0 {
			return buf, nil
		}
		if one[0] == '\n' {
			return buf, nil
		}
		buf = append(buf, one[0])
	}
}

func trimTrailingNL(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return b
}

// hostname for tagging recipients.
func hostname() string {
	host, _ := os.Hostname()
	if i := strings.IndexByte(host, '.'); i > 0 {
		host = host[:i]
	}
	return host
}

// resolveRepo mirrors cmd.resolveRepo without import cycle.
func resolveRepo() (string, error) {
	if v := os.Getenv("VALHALLA_REPO"); v != "" {
		return v, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out, err := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("not inside a git repo (set VALEHALLA_REPO): %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
