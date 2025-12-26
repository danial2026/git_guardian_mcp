# Code Standards

## Go Style Guide

### Formatting

```bash
# Always run before committing
gofmt -w .
```

### Comments

**DO:**
```go
// GetUser fetches user by ID
func GetUser(id string) (*User, error) {
```

**DON'T:**
```go
// GetUser is a function that takes a user ID as a string parameter
// and returns a pointer to a User struct and an error if something
// goes wrong during the fetching process
func GetUser(id string) (*User, error) {
```

**Rule:** One line, describe what it does, not how.

### Error Messages

**DO:**
```go
return nil, fmt.Errorf("failed to parse config: %w", err)
```

**DON'T:**
```go
return nil, errors.New("Error!!! Something bad happened")
```

**Rule:** Lowercase, specific, wrap errors with `%w`.

### Function Size

- **Max 50 lines** per function
- If longer, break into smaller functions
- Each function should do **one thing**

### Variable Names

**DO:**
```go
user := getUser()
count := len(items)
isValid := checkValidity()
```

**DON'T:**
```go
u := getUser()
c := len(items)
validityCheckResult := checkValidity()
```

**Rule:** Clear but concise. Abbreviate if obvious.

### Package Organization

```
pkg/
├── analyzer/    # Static analysis
│   ├── analyzer.go
│   └── helpers.go
├── config/      # Configuration
│   └── config.go
└── mcp/         # MCP server
    └── server.go
```

**Rule:** One concept per package.

---

## File Structure

### Order

1. Package declaration
2. Imports (grouped: stdlib, external, internal)
3. Constants
4. Types
5. Functions

### Example

```go
package analyzer

import (
	"fmt"
	"os"
	
	"gopkg.in/yaml.v3"
	
	"github.com/you/project/pkg/config"
)

const (
	DefaultTimeout = 300
)

type Analyzer struct {
	repoPath string
}

func NewAnalyzer(path string) *Analyzer {
	return &Analyzer{repoPath: path}
}
```

---

## Testing

### Test File Naming

```
file.go      → file_test.go
analyzer.go  → analyzer_test.go
```

### Test Function Naming

```go
func TestGetUser(t *testing.T)              // Simple test
func TestGetUser_NotFound(t *testing.T)     // Specific case
func TestGetUser_InvalidID(t *testing.T)    // Error case
```

### Test Structure

```go
func TestFunction(t *testing.T) {
	// Setup
	input := "test"
	
	// Execute
	result, err := Function(input)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
```

---

## MCP Specific

### Tool Registration

```go
server.RegisterTool(
	"tool_name",           // snake_case
	"Short description",   // One sentence
	handlerFunction,       // Handler func
)
```

### Tool Handlers

```go
func handleTool(params json.RawMessage) (interface{}, error) {
	// 1. Parse params
	var input struct {
		Field string `json:"field"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, err
	}
	
	// 2. Validate
	if input.Field == "" {
		return nil, errors.New("field is required")
	}
	
	// 3. Execute
	result := doSomething(input.Field)
	
	// 4. Return
	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}
```

### Response Format

**Always return:**
```go
map[string]interface{}{
	"success": bool,      // Required
	"data":    anything,  // Your data
	"message": string,    // Optional
}
```

---

## Git Commit Messages

### Format

```
Type: Short description (50 chars max)

Optional detailed explanation.
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `refactor`: Code refactoring
- `test`: Tests
- `chore`: Maintenance

### Examples

**Good:**
```
feat: Add security scanning tool
fix: Handle nil pointer in analyzer
docs: Update README installation section
```

**Bad:**
```
updated stuff
fixed things
asdf
```

---

## Error Handling

### Pattern

```go
if err != nil {
	return nil, fmt.Errorf("context: %w", err)
}
```

### Logging Errors

```go
// Log then return
logger.Printf("ERROR: tool %s failed: %v", name, err)
return nil, err
```

---

## Performance

### Do's

```go
// Reuse buffers
var buf bytes.Buffer
buf.WriteString(data)

// Close resources
defer file.Close()

// Use sync.Pool for repeated allocations
var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
```

### Don'ts

```go
// Don't create goroutines without limits
for i := 0; i < 1000000; i++ {
	go doWork()  // Bad!
}

// Don't ignore errors
file, _ := os.Open(path)  // Bad!
```

---

## Documentation

### Package Documentation

```go
// Package analyzer provides static code analysis.
package analyzer
```

### Public Functions

```go
// RunChecks analyzes files and returns results.
func RunChecks(files []string) []Result {
```

### Private Functions

Optional - only if complex.

---

## Security

### Never

```go
// Don't log sensitive data
logger.Printf("Password: %s", password)  // NO!

// Don't execute unsanitized input
exec.Command(userInput)  // NO!

// Don't ignore errors from security operations
crypto.GenerateKey()  // NO! Check the error!
```

### Always

```go
// Validate input
if !isValid(input) {
	return errors.New("invalid input")
}

// Use constants for sensitive config
const maxRetries = 3

// Close resources
defer file.Close()
```

---

## Tools

### Required

- `gofmt` - Code formatting
- `go vet` - Static analysis

### Recommended

- `golangci-lint` - Comprehensive linting
- `staticcheck` - Advanced checks

### Run Before Commit

```bash
gofmt -w .
go vet ./...
golangci-lint run
go test ./...
```

---

## Questions?

If something isn't covered here, follow these principles:

1. **Clarity over cleverness**
2. **Simple over complex**
3. **Explicit over implicit**
4. **Tested over trusted**

Check existing code for examples, or ask!

---

**Keep it simple, keep it clean! ✨**

