package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/danial2026/git_guardian_mcp/pkg/analyzer"
	"github.com/danial2026/git_guardian_mcp/pkg/config"
	"github.com/danial2026/git_guardian_mcp/pkg/git"
	"github.com/danial2026/git_guardian_mcp/pkg/mcp"
	"github.com/danial2026/git_guardian_mcp/pkg/tests"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("git-guardian-mcp v1.0.0")
		return
	}

	// Log to file to keep Cursor's output clean
	logFile := "/tmp/git-guardian-mcp.log"
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		f = os.NewFile(0, os.DevNull)
	}
	defer f.Close()

	logger := log.New(f, "[git-guardian] ", log.LstdFlags)
	server := mcp.NewServer(logger)

	// Register MCP tools
	server.RegisterTool("analyze_commits", "Analyze unpushed commits for issues", handleAnalyzeCommits)
	server.RegisterTool("run_checks", "Run syntax and static analysis checks", handleRunChecks)
	server.RegisterTool("run_tests", "Execute configured test suites", handleRunTests)
	server.RegisterTool("explain_failure", "Get detailed explanation of a failure", handleExplainFailure)
	server.RegisterTool("validate_push", "Complete validation before push", handleValidatePush)

	// Start MCP server
	if err := server.Start(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}

func handleAnalyzeCommits(params json.RawMessage) (interface{}, error) {
	var input struct {
		RepoPath string `json:"repo_path"`
		Remote   string `json:"remote"`
		Branch   string `json:"branch"`
	}

	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if input.RepoPath == "" {
		input.RepoPath = "."
	}
	if input.Remote == "" {
		input.Remote = "origin"
	}

	gitAnalyzer := git.NewAnalyzer(input.RepoPath)
	commits, err := gitAnalyzer.GetUnpushedCommits(input.Remote, input.Branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get unpushed commits: %w", err)
	}

	return map[string]interface{}{
		"success":       true,
		"commits":       commits,
		"total_commits": len(commits),
		"changed_files": gitAnalyzer.GetChangedFiles(commits),
	}, nil
}

func handleRunChecks(params json.RawMessage) (interface{}, error) {
	var input struct {
		RepoPath string   `json:"repo_path"`
		Files    []string `json:"files"`
	}

	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if input.RepoPath == "" {
		input.RepoPath = "."
	}

	analyzer := analyzer.NewAnalyzer(input.RepoPath)
	results := analyzer.RunChecks(input.Files)

	hasErrors := false
	for _, result := range results {
		if !result.Success {
			hasErrors = true
			break
		}
	}

	return map[string]interface{}{
		"success": !hasErrors,
		"results": results,
	}, nil
}

func handleRunTests(params json.RawMessage) (interface{}, error) {
	var input struct {
		RepoPath   string `json:"repo_path"`
		ConfigPath string `json:"config_path"`
	}

	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	if input.RepoPath == "" {
		input.RepoPath = "."
	}
	if input.ConfigPath == "" {
		input.ConfigPath = ".mcp.yml"
	}

	cfg, err := config.Load(input.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	runner := tests.NewRunner(input.RepoPath, cfg)
	results := runner.RunAll()

	hasBlockingFailures := false
	for _, result := range results {
		if !result.Success && result.Blocking {
			hasBlockingFailures = true
			break
		}
	}

	return map[string]interface{}{
		"success": !hasBlockingFailures,
		"results": results,
	}, nil
}

func handleExplainFailure(params json.RawMessage) (interface{}, error) {
	var input struct {
		FailureType string `json:"failure_type"`
		Details     string `json:"details"`
	}

	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	explanation := analyzer.ExplainFailure(input.FailureType, input.Details)

	return map[string]interface{}{
		"success":     true,
		"explanation": explanation,
	}, nil
}

func handleValidatePush(params json.RawMessage) (interface{}, error) {
	var input struct {
		RepoPath   string `json:"repo_path"`
		Remote     string `json:"remote"`
		Branch     string `json:"branch"`
		ConfigPath string `json:"config_path"`
	}

	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Set defaults
	if input.RepoPath == "" {
		input.RepoPath = "."
	}
	if input.Remote == "" {
		input.Remote = "origin"
	}
	if input.ConfigPath == "" {
		input.ConfigPath = ".mcp.yml"
	}

	// Get unpushed commits
	gitAnalyzer := git.NewAnalyzer(input.RepoPath)
	commits, err := gitAnalyzer.GetUnpushedCommits(input.Remote, input.Branch)
	if err != nil {
		return nil, fmt.Errorf("failed to get unpushed commits: %w", err)
	}

	if len(commits) == 0 {
		return map[string]interface{}{
			"success": true,
			"message": "No unpushed commits to validate",
		}, nil
	}

	changedFiles := gitAnalyzer.GetChangedFiles(commits)

	// Run static analysis
	analyzer := analyzer.NewAnalyzer(input.RepoPath)
	checkResults := analyzer.RunChecks(changedFiles)

	hasCheckErrors := false
	for _, result := range checkResults {
		if !result.Success {
			hasCheckErrors = true
			break
		}
	}

	// Run tests
	cfg, err := config.Load(input.ConfigPath)
	var testResults []tests.TestResult
	hasBlockingTestFailures := false

	if err == nil {
		runner := tests.NewRunner(input.RepoPath, cfg)
		testResults = runner.RunAll()

		for _, result := range testResults {
			if !result.Success && result.Blocking {
				hasBlockingTestFailures = true
				break
			}
		}
	}

	success := !hasCheckErrors && !hasBlockingTestFailures

	return map[string]interface{}{
		"success":       success,
		"commits":       len(commits),
		"changed_files": len(changedFiles),
		"checks":        checkResults,
		"tests":         testResults,
	}, nil
}
