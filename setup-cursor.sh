#!/bin/bash

# Setup script to connect Git Guardian MCP to Cursor

set -e

# Get the absolute path to the binary
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MCP_BINARY="$SCRIPT_DIR/git-guardian-mcp"

# Build if not exists
if [ ! -f "$MCP_BINARY" ]; then
    echo "🔨 Building Git Guardian MCP..."
    cd "$SCRIPT_DIR"
    go build -o git-guardian-mcp
    echo "✅ Built successfully"
fi

# Cursor MCP settings file
CURSOR_MCP_SETTINGS="$HOME/.cursor/mcp_settings.json"

# Create .cursor directory if it doesn't exist
mkdir -p "$HOME/.cursor"

# Check if settings file exists
if [ ! -f "$CURSOR_MCP_SETTINGS" ]; then
    echo "📝 Creating new MCP settings file..."
    cat > "$CURSOR_MCP_SETTINGS" <<EOF
{
  "mcpServers": {
    "git-guardian": {
      "command": "$MCP_BINARY",
      "args": [],
      "env": {}
    }
  }
}
EOF
    echo "✅ Created $CURSOR_MCP_SETTINGS"
else
    echo "⚠️  MCP settings file already exists at:"
    echo "   $CURSOR_MCP_SETTINGS"
    echo ""
    echo "Add this configuration manually:"
    echo ""
    cat <<EOF
{
  "mcpServers": {
    "git-guardian": {
      "command": "$MCP_BINARY",
      "args": [],
      "env": {}
    }
  }
}
EOF
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Restart Cursor"
echo "2. Git Guardian MCP tools will be available"
echo ""
echo "Available tools:"
echo "  - analyze_commits: Analyze unpushed commits"
echo "  - run_checks: Run static analysis"
echo "  - run_tests: Execute test suite"
echo "  - explain_failure: Get detailed error info"
echo "  - validate_push: Full pre-push validation"
echo ""
echo "Binary location: $MCP_BINARY"


