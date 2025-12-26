#!/usr/bin/env bash
#
# Uninstall Git Guardian hooks
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

echo "Uninstalling Git Guardian hooks..."
echo ""

# Remove pre-commit hook
if [ -f "$HOOKS_DIR/pre-commit" ]; then
    rm "$HOOKS_DIR/pre-commit"
    echo -e "${GREEN}✓ pre-commit hook removed${NC}"
else
    echo "pre-commit hook not found"
fi

# Remove pre-push hook
if [ -f "$HOOKS_DIR/pre-push" ]; then
    rm "$HOOKS_DIR/pre-push"
    echo -e "${GREEN}✓ pre-push hook removed${NC}"
else
    echo "pre-push hook not found"
fi

echo ""
echo -e "${GREEN}Uninstallation complete!${NC}"

