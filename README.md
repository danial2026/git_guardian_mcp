# Git Guardian MCP

> AI-powered Git hook validation using Model Context Protocol

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8)](https://golang.org/dl/)

## What is this?

Git Guardian MCP is a developer tool that validates code quality before commits and pushes. It works with AI assistants like Cursor through the Model Context Protocol (MCP).

**Key Features:**
- 🔍 Analyzes unpushed commits
- ✅ Multi-language static analysis (Go, Dart, Bash, JS/TS)
- 🧪 Configurable test runner
- 🤖 AI integration via MCP
- 🚀 Git hooks automation

## Quick Start

```bash
# Build
go build -o git-guardian-mcp

# Connect to Cursor
./setup-cursor.sh

# Install Git hooks (optional)
./scripts/install-hooks.sh
```

## Usage

### With Cursor AI

Once connected, just ask Cursor naturally:
```
"Check my code for issues before I push"
"Validate my unpushed commits"
"Run static analysis"
```

### Command Line

```bash
# Analyze unpushed commits
./git-guardian-mcp  # runs as MCP server

# Install hooks
./scripts/install-hooks.sh

# Remove hooks
./scripts/uninstall-hooks.sh
```

## Configuration

Create `.mcp.yml` in your repository:

```yaml
tests:
  - name: go-tests
    command: go test ./...
    blocking: true
    timeout: 300

  - name: lint
    command: golangci-lint run
    blocking: true
    timeout: 120
```

## MCP Tools

The server exposes these tools to AI assistants:

| Tool | Description |
|------|-------------|
| `analyze_commits` | Analyze unpushed commits |
| `run_checks` | Run static analysis |
| `run_tests` | Execute test suite |
| `explain_failure` | Get error explanations |
| `validate_push` | Full pre-push validation |

## Static Analysis

Automatically runs for:
- **Go**: `gofmt`, `go vet`, `golangci-lint`
- **Dart**: `dart analyze`, `flutter analyze`
- **Bash**: `shellcheck`
- **JS/TS**: `eslint`

## Project Structure

```
.
├── main.go              # Entry point
├── pkg/
│   ├── mcp/            # MCP server implementation
│   ├── git/            # Git operations
│   ├── analyzer/       # Static analysis
│   ├── config/         # Configuration loading
│   └── tests/          # Test runner
├── hooks/              # Git hook templates
├── scripts/            # Setup scripts
└── testcases/          # Example error cases
```

## Development

```bash
# Build
make build

# Run tests
make test

# Format code
make fmt

# Install locally
make install
```

## Requirements

- Go 1.21+
- Git
- Optional: `golangci-lint`, `shellcheck`, `eslint`

## License

MIT - See [LICENSE](LICENSE)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_STANDARDS.md](CODE_STANDARDS.md)

---

**Made for developers who care about code quality 🛡️**
