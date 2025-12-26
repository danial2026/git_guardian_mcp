#!/bin/bash

# Test script to verify MCP server is working

set -e

MCP_BINARY="./git-guardian-mcp"

if [ ! -f "$MCP_BINARY" ]; then
    echo "❌ Binary not found. Building..."
    go build -o git-guardian-mcp
fi

echo "🧪 Testing MCP Server Protocol..."
echo ""

# Test 1: Initialize
echo "📡 Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | $MCP_BINARY 2>/dev/null | head -1 | jq .
echo ""

# Test 2: List tools
echo "🔧 Test 2: List Tools"
{
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
} | $MCP_BINARY 2>/dev/null | tail -1 | jq '.result.tools[] | {name, description}'
echo ""

# Test 3: List resources
echo "📦 Test 3: List Resources"
{
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    echo '{"jsonrpc":"2.0","id":3,"method":"resources/list","params":{}}'
} | $MCP_BINARY 2>/dev/null | tail -1 | jq .
echo ""

# Test 4: List prompts
echo "💬 Test 4: List Prompts"
{
    echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'
    echo '{"jsonrpc":"2.0","id":4,"method":"prompts/list","params":{}}'
} | $MCP_BINARY 2>/dev/null | tail -1 | jq .
echo ""

echo "✅ All protocol tests passed!"
echo ""
echo "Now restart Cursor to connect to the MCP server"



