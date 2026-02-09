#!/usr/bin/env bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_NAME=netxp

echo "Building netxp..."
cd "$ROOT_DIR"
# fetch dependencies (liner)
echo "fetching dependencies..."
go get github.com/peterh/liner
# build (show only first 50 lines)
go build -o "$OUT_NAME" main.go 2>&1 | head -50

INSTALL_DIRS=("$HOME/.local/bin" "/usr/local/bin" "$HOME/bin")
installed="false"
for d in "${INSTALL_DIRS[@]}"; do
  if [ -w "$(dirname "$d")" ] || [ ! -e "$d" ]; then
    mkdir -p "$d" 2>/dev/null || true
    mv -f "$OUT_NAME" "$d/" 2>/dev/null || cp -f "$OUT_NAME" "$d/" 2>/dev/null || true
    if [ -x "$d/$OUT_NAME" ]; then
      echo "Installed to $d/$OUT_NAME"
      installed="true"
      break
    fi
  fi
done

if [ "$installed" = "false" ]; then
  echo "Could not install to a system path automatically."
  echo "You can move '$OUT_NAME' to a directory in your PATH, for example:"
  echo "  mkdir -p \$HOME/bin && mv $OUT_NAME \$HOME/bin/"
  echo "Or run: sudo mv $OUT_NAME /usr/local/bin/"
fi

echo "Done."
