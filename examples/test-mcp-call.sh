#!/usr/bin/env bash
#
# Test MCP server with sample requests
#

set -e

BINARY="./git-guardian-mcp"

if [ ! -f "$BINARY" ]; then
    echo "Building binary..."
    go build -o git-guardian-mcp
fi

echo "Testing MCP Server..."
echo ""

# Test 1: Initialize
echo "Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | $BINARY | head -1
echo ""

# Test 2: List tools
echo "Test 2: List Tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | $BINARY | head -1
echo ""

# Test 3: Analyze commits (will work if in a git repo)
echo "Test 3: Analyze Commits"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"analyze_commits","arguments":{"repo_path":"."}}}' | $BINARY | head -1
echo ""

echo "✓ MCP server tests complete"

