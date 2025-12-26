#!/usr/bin/env bash
#
# Install Git Guardian hooks
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
HOOKS_DIR="$REPO_ROOT/.git/hooks"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOURCE_HOOKS_DIR="$SCRIPT_DIR/../hooks"

echo "Installing Git Guardian hooks..."
echo "Repository: $REPO_ROOT"
echo ""

# Check if git-guardian-mcp binary exists
if [ ! -f "$REPO_ROOT/git-guardian-mcp" ] && ! command -v git-guardian-mcp &> /dev/null; then
    echo -e "${YELLOW}Warning: git-guardian-mcp binary not found${NC}"
    echo "Building binary..."
    cd "$SCRIPT_DIR/.."
    if ! go build -o git-guardian-mcp; then
        echo -e "${RED}Failed to build git-guardian-mcp${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Binary built successfully${NC}"
    cd "$REPO_ROOT"
fi

# Install pre-commit hook
if [ -f "$HOOKS_DIR/pre-commit" ]; then
    echo -e "${YELLOW}pre-commit hook already exists${NC}"
    read -p "Overwrite? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Skipping pre-commit hook"
    else
        cp "$SOURCE_HOOKS_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
        chmod +x "$HOOKS_DIR/pre-commit"
        echo -e "${GREEN}✓ pre-commit hook installed${NC}"
    fi
else
    cp "$SOURCE_HOOKS_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
    chmod +x "$HOOKS_DIR/pre-commit"
    echo -e "${GREEN}✓ pre-commit hook installed${NC}"
fi

# Install pre-push hook
if [ -f "$HOOKS_DIR/pre-push" ]; then
    echo -e "${YELLOW}pre-push hook already exists${NC}"
    read -p "Overwrite? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Skipping pre-push hook"
    else
        cp "$SOURCE_HOOKS_DIR/pre-push" "$HOOKS_DIR/pre-push"
        chmod +x "$HOOKS_DIR/pre-push"
        echo -e "${GREEN}✓ pre-push hook installed${NC}"
    fi
else
    cp "$SOURCE_HOOKS_DIR/pre-push" "$HOOKS_DIR/pre-push"
    chmod +x "$HOOKS_DIR/pre-push"
    echo -e "${GREEN}✓ pre-push hook installed${NC}"
fi

echo ""
echo -e "${GREEN}Installation complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Create .mcp.yml in your repository root (see .mcp.example.yml)"
echo "  2. Commit and push to test the hooks"
echo ""
echo "To uninstall:"
echo "  rm $HOOKS_DIR/pre-commit"
echo "  rm $HOOKS_DIR/pre-push"

