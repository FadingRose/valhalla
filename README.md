# valhalla

Unified dotfiles + secrets manager. Nvim + OpenCode + Kitty + Fish, managed
through symlinks from a single repo, with age-encrypted secrets.

## Structure

```
valhalla/
├── main.go, go.mod, manifest.toml   ← Go binary + declarative config
├── cmd/                             ← CLI commands (cobra)
├── internal/                        ← manifest, linker, syncer, secretage
├── nvim/                            # Neovim config (LazyVim + custom plugins)
├── opencode/                        # OpenCode themes
├── kitty/                           # Kitty terminal config
├── fish/                            # Fish shell config
│   └── conf.d/valhalla-secrets.fish ← sources ~/.local/share/valhalla/secrets.fish
├── secrets/                         # age-encrypted secrets (safe to commit)
│   ├── recipients.txt               # device public keys
│   ├── llm/*.age
│   ├── wallets/*.age
│   └── ...
├── scripts/                         ← legacy install.sh / sync.sh (fallback)
└── Makefile (TODO)
```

## Quick Start

### Prerequisites

| Tool | macOS | Arch Linux |
|------|-------|------------|
| Git | `xcode-select --install` | `sudo pacman -S git` |
| Go ≥ 1.21 | `brew install go` | `sudo pacman -S go` |
| Neovim ≥ 0.10 | `brew install neovim` | `sudo pacman -S neovim` |
| Kitty | `brew install --cask kitty` | `sudo pacman -S kitty` |
| Fish | `brew install fish` | `sudo pacman -S fish` |
| Nerd Font | `brew install --cask font-maple-mono-nf-cn` | `sudo pacman -S ttf-maple-mono-cn` |
| Rust (nvim deps) | `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs \| sh` | `sudo pacman -S rustup` |
| Node.js | `brew install node` | `sudo pacman -S nodejs npm` |

### Install on a new machine

```bash
# 1. Clone
git clone git@github.com:FadingRose/valhalla.git ~/valhalla
cd ~/valhalla

# 2. Build the binary
go build -o valhalla .

# 3. Create all symlinks (nvim, fish, kitty, opencode, ...)
./valhalla apply

# 4. If this machine should have secrets (optional):
#    a. Copy the age identity from an existing device:
scp existing-device:~/.valhalla/identity.txt ~/.valhalla/identity.txt
chmod 600 ~/.valhalla/identity.txt
#    b. Decrypt and inject secrets into fish:
./valhalla secret inject
```

Open `nvim` — LazyVim auto-installs plugins. Restart fish — secrets are loaded.

## Commands

```bash
valhalla                              # show help
valhalla apply [--dry-run]            # create symlinks from manifest.toml
valhalla diff                         # show missing/broken symlinks
valhalla sync [--message MSG]         # commit + pull --rebase + push

valhalla secret init [--name NAME]    # generate age keypair + register device
valhalla secret add <path> [--stdin]  # encrypt and store a value
valhalla secret get <path> [--copy]   # decrypt and print / copy to clipboard
valhalla secret list                  # list all stored paths
valhalla secret rm <path> [--force]   # remove a secret (manifest-aware)
valhalla secret inject [--shell fish] # render all secrets to a shell file
valhalla secret reencrypt             # re-encrypt for current devices

valhalla devices list                 # show all registered devices
valhalla devices add <pubkey>         # register a new device
valhalla devices remove <name>        # revoke a device
```

### Shell completion

```bash
valhalla completion install           # install fish completions
```

Provides `<TAB>` completion for secret paths, device names, and shell formats:
- `valhalla secret get <TAB>` → lists all stored secrets
- `valhalla secret rm <TAB>` → lists all stored secrets
- `valhalla devices remove <TAB>` → lists registered device names
- `valhalla secret inject --shell <TAB>` → fish|posix|env|json

Restart fish (`exec fish`) after install to activate.

### Common workflows

**Update a secret** (e.g. rotate an API key):
```bash
valhalla secret add llm/openrouter    # prompts for value
# or:
echo "sk-newvalue..." | valhalla secret add llm/openrouter --stdin
valhalla secret inject                # regenerate the fish file
exec fish                             # reload shell
```

**Add a brand new secret**:
```bash
# 1. Store it
valhalla secret add llm/newservice
# 2. Declare it in manifest.toml:
#    [[secrets.inject]]
#    var = "NEWSERVICE_API_KEY"
#    path = "llm/newservice"
# 3. Regenerate
valhalla secret inject
```

## How secrets work

```
┌────────────────────────────────────────────────────────────┐
│ valhalla/secrets/*.age   (ciphertext, git-tracked)         │
│   encrypted with age (X25519), safe to push to GitHub      │
└─────────────────┬──────────────────────────────────────────┘
                  │ valhalla secret inject
                  ▼
┌────────────────────────────────────────────────────────────┐
│ ~/.local/share/valhalla/secrets.fish  (plaintext, mode 600)│
│   NEVER committed — outside the repo                        │
└─────────────────┬──────────────────────────────────────────┘
                  │ sourced at shell startup by
                  ▼
┌────────────────────────────────────────────────────────────┐
│ fish/conf.d/valhalla-secrets.fish  (wrapper, git-tracked)  │
│   safe to commit — only contains `source` command          │
└────────────────────────────────────────────────────────────┘
```

**Security properties**:
- Encryption: age (X25519, curated by Filippo Valsorda)
- Repo contains only ciphertext — safe on public GitHub
- Plaintext file lives outside repo with `0600` permissions
- Private key at `~/.valhalla/identity.txt` — never in git
- Removing a device: revoke its pubkey, re-encrypt all secrets (TODO)

## manifest.toml

The single source of truth. Declares:
- `[[links]]` — symlinks to create (`valhalla apply`)
- `[[profiles]]` — per-device overrides (matched by hostname)
- `[secrets]` + `[[secrets.inject]]` — what to materialize and where

See [`manifest.toml`](manifest.toml) for the full example.

## Daily workflow

After editing any config (nvim, fish, etc.):
```bash
valhalla sync          # smart commit message + push
```

On other machines:
```bash
cd ~/valhalla && git pull
valhalla apply         # pick up any new symlinks
```

## Legacy fallback

The original bash scripts still work if you don't want to build Go:
```bash
bash scripts/install.sh     # same as `valhalla apply`
bash scripts/sync.sh        # same as `valhalla sync`
```

## Theme integration

The alter-avenger palette is defined in three places, kept in sync:

| File | Purpose |
|------|---------|
| `nvim/colors/alter-avenger.lua` | Nvim colorscheme + exports `M.palette`, `M.character`, `M.neural_scramble()` |
| `opencode/themes/alter-avenger.json` | OpenCode TUI theme |
| `nvim/lua/plugins/lualine.lua` | Reads palette at runtime |

## Set Fish as default shell

```bash
echo $(which fish) | sudo tee -a /etc/shells
chsh -s $(which fish)
```
