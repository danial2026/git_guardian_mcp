# Contributing to Git Guardian MCP

Thanks for your interest! Here's how to contribute.

## Quick Start

1. Fork the repo
2. Create a branch: `git checkout -b feature/your-feature`
3. Make changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit: `git commit -m "Add feature"`
7. Push: `git push origin feature/your-feature`
8. Open a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/danial2026/git_guardian_mcp
cd git_guardian_mcp

# Build
go build -o git_guardian_mcp

# Run tests
go test ./...

# Format
gofmt -w .
```

## Code Guidelines

- Follow [CODE_STANDARDS.md](CODE_STANDARDS.md)
- Write clear, simple code
- Add tests for new features
- Keep comments short and useful
- Use `gofmt` before committing

## Adding Features

See [ADDING_FEATURES.md](ADDING_FEATURES.md) for detailed guides on:
- Adding new MCP tools
- Supporting new languages
- Adding analyzers
- Extending configuration

## Pull Request Process

1. **Update tests** - Add tests for new features
2. **Run checks** - Ensure `make test` and `make fmt` pass
3. **Update docs** - Add to README if needed
4. **Small PRs** - Keep changes focused
5. **Clear description** - Explain what and why

## Reporting Issues

**Bugs:**
- Go version
- OS
- Steps to reproduce
- Expected vs actual behavior

**Feature Requests:**
- Use case
- Proposed solution
- Why it's needed

## Code Review

We look for:
- ✅ Tests pass
- ✅ Code is formatted
- ✅ Clear commit messages
- ✅ Documentation updated
- ✅ Follows standards

## Questions?

Open an issue with the `question` label.
