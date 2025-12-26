package tests

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/danial2026/git_guardian_mcp/pkg/config"
)

// TestResult represents the result of a test run
type TestResult struct {
	Name     string  `json:"name"`
	Success  bool    `json:"success"`
	Blocking bool    `json:"blocking"`
	Duration float64 `json:"duration"` // in seconds
	Output   string  `json:"output"`
	Error    string  `json:"error,omitempty"`
}

// Runner handles test execution
type Runner struct {
	repoPath string
	config   *config.Config
}

// NewRunner creates a new test runner
func NewRunner(repoPath string, config *config.Config) *Runner {
	return &Runner{
		repoPath: repoPath,
		config:   config,
	}
}

// RunAll executes all configured tests
func (r *Runner) RunAll() []TestResult {
	results := make([]TestResult, 0, len(r.config.Tests))

	for _, testConfig := range r.config.Tests {
		result := r.runTest(testConfig)
		results = append(results, result)
	}

	return results
}

func (r *Runner) runTest(testConfig config.TestConfig) TestResult {
	start := time.Now()

	// Create context with timeout
	timeout := time.Duration(testConfig.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Parse command
	parts := strings.Fields(testConfig.Command)
	if len(parts) == 0 {
		return TestResult{
			Name:     testConfig.Name,
			Success:  false,
			Blocking: testConfig.Blocking,
			Duration: 0,
			Error:    "empty command",
		}
	}

	// Execute command
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = r.repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start).Seconds()

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}
	output = strings.TrimSpace(output)

	result := TestResult{
		Name:     testConfig.Name,
		Success:  err == nil,
		Blocking: testConfig.Blocking,
		Duration: duration,
		Output:   output,
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = fmt.Sprintf("test timed out after %d seconds", testConfig.Timeout)
		} else {
			result.Error = err.Error()
		}
	}

	return result
}
