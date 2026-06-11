#!/bin/bash
set -euo pipefail

REPO_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_DIR="$HOME/.config"

link() {
  local src="$1"
  local dest="$2"
  if [ -L "$dest" ]; then
    existing=$(readlink "$dest")
    if [ "$existing" = "$src" ]; then
      echo "  ✓ $dest → $src (already linked)"
      return
    fi
    echo "  ! $dest → $existing (updating to $src)"
    rm "$dest"
  elif [ -d "$dest" ]; then
    echo "  ! $dest exists as directory (backing up → ${dest}.bak)"
    mv "$dest" "${dest}.bak"
  elif [ -f "$dest" ]; then
    echo "  ! $dest exists as file (backing up → ${dest}.bak)"
    mv "$dest" "${dest}.bak"
  fi
  ln -s "$src" "$dest"
  echo "  ✓ $dest → $src"
}

echo "valhalla dotfiles installer"
echo "repo: $REPO_DIR"
echo ""

mkdir -p "$CONFIG_DIR"

echo "[nvim]"
link "$REPO_DIR/nvim" "$CONFIG_DIR/nvim"

echo ""
echo "[opencode]"
mkdir -p "$CONFIG_DIR/opencode/themes"
if [ -L "$CONFIG_DIR/opencode/themes/alter-avenger.json" ]; then
  rm "$CONFIG_DIR/opencode/themes/alter-avenger.json"
fi
link "$REPO_DIR/opencode/themes/alter-avenger.json" "$CONFIG_DIR/opencode/themes/alter-avenger.json"
link "$REPO_DIR/opencode/tui.json" "$CONFIG_DIR/opencode/tui.json"

echo ""
echo "[kitty]"
link "$REPO_DIR/kitty" "$CONFIG_DIR/kitty"

echo ""
echo "[fish]"
link "$REPO_DIR/fish" "$CONFIG_DIR/fish"

echo ""
echo "done. restart your shell/terminal to apply."
