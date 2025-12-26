package analyzer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckResult represents the result of a static analysis check
type CheckResult struct {
	Tool     string   `json:"tool"`
	File     string   `json:"file,omitempty"`
	Line     int      `json:"line,omitempty"`
	Column   int      `json:"column,omitempty"`
	Severity string   `json:"severity"`
	Message  string   `json:"message"`
	Success  bool     `json:"success"`
	Output   string   `json:"output,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

// Analyzer handles static analysis checks
type Analyzer struct {
	repoPath string
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(repoPath string) *Analyzer {
	return &Analyzer{
		repoPath: repoPath,
	}
}

// RunChecks runs all applicable checks on the given files
func (a *Analyzer) RunChecks(files []string) []CheckResult {
	results := make([]CheckResult, 0)

	// Group files by type
	filesByType := a.groupFilesByType(files)

	// Run Go checks
	if goFiles := filesByType["go"]; len(goFiles) > 0 {
		results = append(results, a.checkGo(goFiles)...)
	}

	// Run Dart/Flutter checks
	if dartFiles := filesByType["dart"]; len(dartFiles) > 0 {
		results = append(results, a.checkDart(dartFiles)...)
	}

	// Run Bash checks
	if bashFiles := filesByType["bash"]; len(bashFiles) > 0 {
		results = append(results, a.checkBash(bashFiles)...)
	}

	// Run JavaScript/TypeScript checks
	if jsFiles := filesByType["js"]; len(jsFiles) > 0 {
		results = append(results, a.checkJavaScript(jsFiles)...)
	}

	return results
}

func (a *Analyzer) groupFilesByType(files []string) map[string][]string {
	groups := make(map[string][]string)

	for _, file := range files {
		// Check if file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file))
		switch ext {
		case ".go":
			groups["go"] = append(groups["go"], file)
		case ".dart":
			groups["dart"] = append(groups["dart"], file)
		case ".sh", ".bash":
			groups["bash"] = append(groups["bash"], file)
		case ".js", ".jsx", ".ts", ".tsx":
			groups["js"] = append(groups["js"], file)
		}
	}

	return groups
}

func (a *Analyzer) checkGo(files []string) []CheckResult {
	results := make([]CheckResult, 0)

	// Check if go is available
	if !commandExists("go") {
		return results
	}

	// Run gofmt
	fmtResult := a.runGoFmt(files)
	results = append(results, fmtResult)

	// Run go vet
	vetResult := a.runGoVet()
	results = append(results, vetResult)

	// Run golangci-lint if available
	if commandExists("golangci-lint") {
		lintResult := a.runGolangciLint(files)
		results = append(results, lintResult)
	}

	return results
}

func (a *Analyzer) runGoFmt(files []string) CheckResult {
	// Check each file individually for formatting issues
	cmd := exec.Command("gofmt", "-l")
	cmd.Args = append(cmd.Args, files...)
	cmd.Dir = a.repoPath

	output, err := cmd.Output()
	outputStr := strings.TrimSpace(string(output))

	if err != nil || outputStr != "" {
		errors := []string{}
		if outputStr != "" {
			errors = append(errors, "Files not formatted:\n"+outputStr)
		}
		return CheckResult{
			Tool:     "gofmt",
			Severity: "error",
			Message:  "Go formatting issues found",
			Success:  false,
			Output:   outputStr,
			Errors:   errors,
		}
	}

	return CheckResult{
		Tool:     "gofmt",
		Severity: "info",
		Message:  "All Go files properly formatted",
		Success:  true,
	}
}

func (a *Analyzer) runGoVet() CheckResult {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = a.repoPath

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		return CheckResult{
			Tool:     "go vet",
			Severity: "error",
			Message:  "Go vet found issues",
			Success:  false,
			Output:   outputStr,
			Errors:   []string{outputStr},
		}
	}

	return CheckResult{
		Tool:     "go vet",
		Severity: "info",
		Message:  "No issues found by go vet",
		Success:  true,
	}
}

func (a *Analyzer) runGolangciLint(files []string) CheckResult {
	cmd := exec.Command("golangci-lint", "run", "--out-format=colored-line-number")
	cmd.Dir = a.repoPath

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		return CheckResult{
			Tool:     "golangci-lint",
			Severity: "error",
			Message:  "Linter found issues",
			Success:  false,
			Output:   outputStr,
			Errors:   []string{outputStr},
		}
	}

	return CheckResult{
		Tool:     "golangci-lint",
		Severity: "info",
		Message:  "No linting issues found",
		Success:  true,
	}
}

func (a *Analyzer) checkDart(files []string) []CheckResult {
	results := make([]CheckResult, 0)

	// Check if dart is available
	if !commandExists("dart") {
		return results
	}

	cmd := exec.Command("dart", "analyze")
	cmd.Dir = a.repoPath

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		results = append(results, CheckResult{
			Tool:     "dart analyze",
			Severity: "error",
			Message:  "Dart analysis found issues",
			Success:  false,
			Output:   outputStr,
			Errors:   []string{outputStr},
		})
	} else {
		results = append(results, CheckResult{
			Tool:     "dart analyze",
			Severity: "info",
			Message:  "No Dart issues found",
			Success:  true,
		})
	}

	// Check for Flutter
	if commandExists("flutter") {
		cmd = exec.Command("flutter", "analyze")
		cmd.Dir = a.repoPath

		output, err = cmd.CombinedOutput()
		outputStr = strings.TrimSpace(string(output))

		if err != nil {
			results = append(results, CheckResult{
				Tool:     "flutter analyze",
				Severity: "error",
				Message:  "Flutter analysis found issues",
				Success:  false,
				Output:   outputStr,
				Errors:   []string{outputStr},
			})
		} else {
			results = append(results, CheckResult{
				Tool:     "flutter analyze",
				Severity: "info",
				Message:  "No Flutter issues found",
				Success:  true,
			})
		}
	}

	return results
}

func (a *Analyzer) checkBash(files []string) []CheckResult {
	results := make([]CheckResult, 0)

	if !commandExists("shellcheck") {
		return results
	}

	for _, file := range files {
		cmd := exec.Command("shellcheck", "-f", "gcc", file)
		output, err := cmd.CombinedOutput()
		outputStr := strings.TrimSpace(string(output))

		if err != nil {
			results = append(results, CheckResult{
				Tool:     "shellcheck",
				File:     file,
				Severity: "error",
				Message:  "Shellcheck found issues",
				Success:  false,
				Output:   outputStr,
				Errors:   []string{outputStr},
			})
		} else {
			results = append(results, CheckResult{
				Tool:     "shellcheck",
				File:     file,
				Severity: "info",
				Message:  "No shell script issues",
				Success:  true,
			})
		}
	}

	return results
}

func (a *Analyzer) checkJavaScript(files []string) []CheckResult {
	results := make([]CheckResult, 0)

	if !commandExists("eslint") {
		return results
	}

	cmd := exec.Command("eslint")
	cmd.Args = append(cmd.Args, files...)
	cmd.Dir = a.repoPath

	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		results = append(results, CheckResult{
			Tool:     "eslint",
			Severity: "error",
			Message:  "ESLint found issues",
			Success:  false,
			Output:   outputStr,
			Errors:   []string{outputStr},
		})
	} else {
		results = append(results, CheckResult{
			Tool:     "eslint",
			Severity: "info",
			Message:  "No ESLint issues found",
			Success:  true,
		})
	}

	return results
}

// ExplainFailure provides detailed explanation for a failure
func ExplainFailure(failureType, details string) string {
	explanations := map[string]string{
		"gofmt":           "Go files must be formatted using 'gofmt'. Run 'gofmt -w .' to fix formatting issues.",
		"go vet":          "Go vet found potential issues in your code. Review the errors and fix them before pushing.",
		"golangci-lint":   "The linter found code quality issues. Address the reported problems or configure exceptions in .golangci.yml.",
		"dart analyze":    "Dart analyzer found issues. Run 'dart fix --apply' to auto-fix some issues.",
		"flutter analyze": "Flutter analyzer found issues in your Flutter code. Review and fix them.",
		"shellcheck":      "Shellcheck found issues in your shell scripts. Review the suggestions and fix critical issues.",
		"eslint":          "ESLint found JavaScript/TypeScript issues. Run 'eslint --fix' to auto-fix some issues.",
		"test":            "Tests failed. Review the test output and fix failing tests before pushing.",
	}

	explanation := explanations[failureType]
	if explanation == "" {
		explanation = fmt.Sprintf("Check '%s' failed. Review the details and fix the issues.", failureType)
	}

	if details != "" {
		explanation += "\n\nDetails:\n" + details
	}

	return explanation
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
