#!/bin/sh
set -e

REPO="RusilKoirala/fireworks"
BIN="fireworks"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# detect OS
os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $os (on Windows use fireworks-setup.exe)"; exit 1 ;;
esac

# detect CPU
arch=$(uname -m)
case "$arch" in
  x86_64|amd64)  arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "Unsupported architecture: $arch"; exit 1 ;;
esac

url="https://github.com/$REPO/releases/latest/download/$BIN-$os-$arch.tar.gz"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $BIN for $os/$arch..."
curl -fsSL "$url" -o "$tmp/$BIN.tar.gz"
tar -xzf "$tmp/$BIN.tar.gz" -C "$tmp"

mkdir -p "$INSTALL_DIR"
install -m 755 "$tmp/$BIN" "$INSTALL_DIR/$BIN"
echo "Installed to $INSTALL_DIR/$BIN"


case ":$PATH:" in
  *":$INSTALL_DIR:"*) echo "Run: $BIN" ;;
  *)
    case "$(basename "${SHELL:-sh}")" in
      zsh)  rc="$HOME/.zshrc" ;;
      bash) if [ "$os" = "darwin" ]; then rc="$HOME/.bash_profile"; else rc="$HOME/.bashrc"; fi ;;
      fish) rc="$HOME/.config/fish/config.fish" ;;
      *)    rc="$HOME/.profile" ;;
    esac

    if [ "$(basename "${SHELL:-sh}")" = "fish" ]; then
      line="fish_add_path $INSTALL_DIR"
    else
      line="export PATH=\"$INSTALL_DIR:\$PATH\""
    fi

    mkdir -p "$(dirname "$rc")"
    if ! grep -qsF "$line" "$rc"; then
      printf '\n# fireworks\n%s\n' "$line" >> "$rc"
    fi
    echo "Added $INSTALL_DIR to PATH in $rc"
    echo "Restart your terminal (or run: source $rc), then run: $BIN"
    ;;
esac