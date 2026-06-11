# valhalla

Unified dotfiles. Nvim + OpenCode + Kitty + Fish, managed through symlinks from a single repo.

## Structure

```
valhalla/
├── nvim/                      # Neovim config (LazyVim + custom plugins + alter-avenger theme)
├── opencode/
│   ├── themes/alter-avenger.json   # OpenCode TUI theme (matches nvim colorscheme)
│   └── tui.json                    # OpenCode TUI config (theme selection)
├── kitty/
│   └── kitty.conf             # Kitty terminal config
├── fish/                      # Fish shell config (functions, completions, themes)
├── scripts/
│   ├── install.sh             # Create all symlinks
│   └── sync.sh                # Git add + commit + push
└── .gitignore
```

## Quick Start

### Prerequisites

| Tool | macOS | Arch Linux |
|------|-------|------------|
| Git | `xcode-select --install` | `sudo pacman -S git` |
| Neovim ≥ 0.10 | `brew install neovim` | `sudo pacman -S neovim` |
| Kitty | `brew install --cask kitty` | `sudo pacman -S kitty` |
| Fish | `brew install fish` | `sudo pacman -S fish` |
| Nerd Font | `brew install --cask font-maple-mono-nf-cn` | `sudo pacman -S ttf-maple-mono-cn` |
| Rust (for some nvim plugins) | `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs \| sh` | `sudo pacman -S rustup` |
| Node.js | `brew install node` | `sudo pacman -S nodejs npm` |
| A clipboard tool | built-in (`pbcopy`/`pbpaste`) | `sudo pacman -S xclip` |
| `fzf` (optional) | `brew install fzf` | `sudo pacman -S fzf` |
| `ripgrep` | `brew install ripgrep` | `sudo pacman -S ripgrep` |

### Install

```bash
# 1. Clone the repo
git clone git@github.com:FadingRose/valhalla.git ~/valhalla

# 2. Run the installer (creates symlinks, backs up existing configs)
bash ~/valhalla/scripts/install.sh
```

The installer will:
- Back up any existing config to `.bak`
- Create symlinks: `~/.config/nvim` → `~/valhalla/nvim`, etc.
- Keep your OpenCode `agent/`, `skills/`, `package.json` intact (only links theme + tui.json)

### First Launch

Open `nvim`. LazyVim will auto-install all plugins. After plugins install, the **alter-avenger** colorscheme activates automatically.

## Daily Workflow

### Sync Changes

After editing any config:

```bash
bash ~/valhalla/scripts/sync.sh
```

This commits and pushes all changes with a timestamp.

### Set Fish as Default Shell

```bash
# macOS
echo $(which fish) | sudo tee -a /etc/shells
chsh -s $(which fish)

# Arch Linux
echo $(which fish) | sudo tee -a /etc/shells
chsh -s $(which fish)
```

## What's Inside

### Nvim

- **LazyVim** as the base framework
- **alter-avenger** colorscheme — dark purple + crimson + gold, inspired by Jeanne d'Arc Alter
- **Custom plugins**: blaming (git blame heatmap), solidity-state-variable-highlighter, local-variable-highlighter, auditscope/glance (attention heatmap), sol-import, source-code-fetcher (Etherscan), blockchain-explorers
- **Flash.nvim** for treesitter-aware jumps, **Hop** for word/line jumping
- **Neoscroll** for smooth scrolling, **Markview** for rendered markdown
- **Lualine** with dynamic character voice lines and themed buffer tabline

### OpenCode Theme

`alter-avenger.json` maps the nvim colorscheme palette to all OpenCode TUI tokens (syntax, diff, markdown, UI). Same visual identity across both tools.

### Kitty

- Maple Mono NF CN font
- Fade-style centered tab bar
- Cursor trail effect
- Cell height 120%

### Fish

- Tide prompt
- fzf integration
- nvm for Node.js version management

## Adding a New Machine

```bash
git clone git@github.com:FadingRose/valhalla.git ~/valhalla
bash ~/valhalla/scripts/install.sh

# Install Neovim plugins (auto on first launch, or manual):
nvim --headless "+Lazy! sync" +qa

# Install Fish plugins
fish -c "fisher update"
```

## Modify and Distribute

1. Edit files directly under `~/valhalla/`
2. Run `bash ~/valhalla/scripts/sync.sh` to push
3. On other machines: `cd ~/valhalla && git pull`

## Theme Integration

The alter-avenger palette is defined in three places, kept in sync:

| File | Purpose |
|------|---------|
| `nvim/colors/alter-avenger.lua` | Nvim colorscheme + exports `M.palette`, `M.character`, `M.neural_scramble()` |
| `opencode/themes/alter-avenger.json` | OpenCode TUI theme (defs → theme tokens) |
| `nvim/lua/plugins/lualine.lua` | Reads palette at runtime for statusline colors |
| `nvim/lua/plugins/dashboard.lua` | Reads `neural_scramble()` for glitch effects |

When changing colors, update `M.palette` in `alter-avenger.lua` first, then sync the hex values to the OpenCode JSON theme.
