#!/bin/bash
# Cross-compiles the Gender-API CLI tool for Linux, Windows, and macOS.
# Automatically increments the patch version and outputs to public/cli-client

set -e

# Version calculation
VERSION_FILE=".version"
if [ ! -f "$VERSION_FILE" ]; then
    echo "1.0.0" > "$VERSION_FILE"
fi

CURRENT_VERSION=$(cat "$VERSION_FILE")
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"

# Increment patch version
NEW_PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$NEW_PATCH"

# Save new version
echo "$NEW_VERSION" > "$VERSION_FILE"
echo "==> Bumping version from $CURRENT_VERSION to v$NEW_VERSION"

# Set output directory to the public folder so users can download it
OUT_DIR="../public/cli-client"
APP_NAME="gender-api-cli"

mkdir -p "$OUT_DIR"

# Define the target OS and architectures
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "==> Starting cross-compilation for v$NEW_VERSION..."

for target in "${TARGETS[@]}"; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    
    OUTPUT_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "    Building $GOOS/$GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-X 'main.Version=$NEW_VERSION'" -o "$OUT_DIR/$OUTPUT_NAME" main.go
done

echo ""
echo "==> Build complete. Public binaries are located in '$OUT_DIR/'."
