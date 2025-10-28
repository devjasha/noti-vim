package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/devjasha/noti-vim/internal/config"
)

// IsGitRepo checks if the notes directory is a git repository
func IsGitRepo() bool {
	cfg := config.Get()
	gitDir := filepath.Join(cfg.NotesDir, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

// Init initializes a git repository in the notes directory
func Init() error {
	cfg := config.Get()

	if IsGitRepo() {
		return fmt.Errorf("git repository already exists")
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not initialize git repository: %w\n%s", err, output)
	}

	return nil
}

// Status returns the git status of the notes directory
func Status() (string, error) {
	cfg := config.Get()

	if !IsGitRepo() {
		return "", fmt.Errorf("not a git repository (use 'noti git init' to initialize)")
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not get git status: %w\n%s", err, output)
	}

	return string(output), nil
}

// StatusShort returns a short git status summary
func StatusShort() (string, error) {
	cfg := config.Get()

	if !IsGitRepo() {
		return "not a git repository", nil
	}

	cmd := exec.Command("git", "status", "--short", "--branch")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not get git status: %w\n%s", err, output)
	}

	return strings.TrimSpace(string(output)), nil
}

// Add stages all changes in the notes directory
func Add() error {
	cfg := config.Get()

	if !IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not stage changes: %w\n%s", err, output)
	}

	return nil
}

// Commit creates a git commit with the given message
func Commit(message string) error {
	cfg := config.Get()

	if !IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	// Stage all changes first
	if err := Add(); err != nil {
		return err
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if there's nothing to commit
		if strings.Contains(string(output), "nothing to commit") {
			return fmt.Errorf("nothing to commit, working tree clean")
		}
		return fmt.Errorf("could not create commit: %w\n%s", err, output)
	}

	return nil
}

// Push pushes commits to the remote repository
func Push() error {
	cfg := config.Get()

	if !IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "push")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not push to remote: %w\n%s", err, output)
	}

	return nil
}

// Pull pulls changes from the remote repository
func Pull() error {
	cfg := config.Get()

	if !IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not pull from remote: %w\n%s", err, output)
	}

	return nil
}

// Sync performs a full sync: add, commit, pull, and push
func Sync(message string) error {
	if !IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}

	// Check if there are any changes to commit
	status, err := Status()
	if err != nil {
		return err
	}

	hasChanges := strings.TrimSpace(status) != ""

	// Only commit if there are changes
	if hasChanges {
		if message == "" {
			message = "Auto-sync notes"
		}

		if err := Commit(message); err != nil {
			return fmt.Errorf("failed to commit changes: %w", err)
		}
	}

	// Pull changes from remote
	if err := Pull(); err != nil {
		return fmt.Errorf("failed to pull changes: %w", err)
	}

	// Push if we made a commit
	if hasChanges {
		if err := Push(); err != nil {
			return fmt.Errorf("failed to push changes: %w", err)
		}
	}

	return nil
}

// Log returns the git log
func Log(maxCount int) (string, error) {
	cfg := config.Get()

	if !IsGitRepo() {
		return "", fmt.Errorf("not a git repository")
	}

	args := []string{"log", "--oneline", "--decorate", "--color=always"}
	if maxCount > 0 {
		args = append(args, fmt.Sprintf("-n%d", maxCount))
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not get git log: %w\n%s", err, output)
	}

	return string(output), nil
}

// HasRemote checks if a git remote is configured
func HasRemote() (bool, error) {
	cfg := config.Get()

	if !IsGitRepo() {
		return false, nil
	}

	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = cfg.NotesDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("could not check git remotes: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}
