package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Commit represents a Git commit
type Commit struct {
	Hash    string   `json:"hash"`
	Author  string   `json:"author"`
	Date    string   `json:"date"`
	Message string   `json:"message"`
	Files   []string `json:"files"`
	Diff    string   `json:"diff,omitempty"`
}

// Analyzer handles Git repository analysis
type Analyzer struct {
	repoPath string
}

// NewAnalyzer creates a new Git analyzer
func NewAnalyzer(repoPath string) *Analyzer {
	return &Analyzer{
		repoPath: repoPath,
	}
}

// GetUnpushedCommits retrieves commits that haven't been pushed to remote
func (a *Analyzer) GetUnpushedCommits(remote, branch string) ([]Commit, error) {
	// Get current branch if not specified
	if branch == "" {
		cmd := exec.Command("git", "-C", a.repoPath, "branch", "--show-current")
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get current branch: %w", err)
		}
		branch = strings.TrimSpace(string(output))
	}

	// Check if remote branch exists
	remoteBranch := fmt.Sprintf("%s/%s", remote, branch)
	cmd := exec.Command("git", "-C", a.repoPath, "rev-parse", "--verify", remoteBranch)
	if err := cmd.Run(); err != nil {
		// Remote branch doesn't exist, return all commits
		return a.getAllCommits()
	}

	// Get unpushed commits
	cmd = exec.Command("git", "-C", a.repoPath, "log", fmt.Sprintf("%s..HEAD", remoteBranch), "--format=%H|||%an|||%ai|||%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get unpushed commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|||")
		if len(parts) < 4 {
			continue
		}

		commit := Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		}

		// Get files changed in this commit
		files, err := a.getCommitFiles(commit.Hash)
		if err == nil {
			commit.Files = files
		}

		// Get diff for this commit
		diff, err := a.getCommitDiff(commit.Hash)
		if err == nil {
			commit.Diff = diff
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// getAllCommits gets all commits (used when remote doesn't exist)
func (a *Analyzer) getAllCommits() ([]Commit, error) {
	cmd := exec.Command("git", "-C", a.repoPath, "log", "--format=%H|||%an|||%ai|||%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get all commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|||")
		if len(parts) < 4 {
			continue
		}

		commit := Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Date:    parts[2],
			Message: parts[3],
		}

		files, err := a.getCommitFiles(commit.Hash)
		if err == nil {
			commit.Files = files
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// getCommitFiles returns the list of files changed in a commit
func (a *Analyzer) getCommitFiles(hash string) ([]string, error) {
	cmd := exec.Command("git", "-C", a.repoPath, "diff-tree", "--no-commit-id", "--name-only", "-r", hash)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// getCommitDiff returns the diff for a commit
func (a *Analyzer) getCommitDiff(hash string) (string, error) {
	cmd := exec.Command("git", "-C", a.repoPath, "show", hash, "--format=", "--no-color")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetChangedFiles extracts all unique files from a list of commits
func (a *Analyzer) GetChangedFiles(commits []Commit) []string {
	fileSet := make(map[string]bool)
	for _, commit := range commits {
		for _, file := range commit.Files {
			// Convert to absolute path
			absPath := filepath.Join(a.repoPath, file)
			fileSet[absPath] = true
		}
	}

	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	return files
}

// GetStagedFiles returns list of staged files
func (a *Analyzer) GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "-C", a.repoPath, "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			absPath := filepath.Join(a.repoPath, line)
			files = append(files, absPath)
		}
	}
	return files, nil
}
