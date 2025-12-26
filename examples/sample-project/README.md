# Sample Project

This is an example project showing how to use Git Guardian MCP.

## Setup

1. Install Git Guardian hooks:
   ```bash
   ../../scripts/install-hooks.sh
   ```

2. The `.mcp.yml` file is already configured with sample tests.

3. Try making a commit:
   ```bash
   echo "// test" >> main.go
   git add main.go
   git commit -m "Test commit"
   ```

## Configuration

See `.mcp.yml` for the test configuration used in this project.

## Expected Behavior

- **Pre-commit**: Runs static analysis on staged files
- **Pre-push**: Runs full test suite before allowing push

## Customization

Edit `.mcp.yml` to add or modify tests for your project needs.

