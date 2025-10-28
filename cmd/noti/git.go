package main

import (
	"fmt"

	"github.com/devjasha/noti-vim/internal/git"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git operations for notes",
	Long:  `Manage version control for your notes using git`,
}

var gitInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize git repository",
	Long:  `Initialize a git repository in the notes directory`,
	RunE:  runGitInit,
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git status",
	Long:  `Show the git status of your notes directory`,
	RunE:  runGitStatus,
}

var gitCommitCmd = &cobra.Command{
	Use:   "commit <message>",
	Short: "Commit changes",
	Long:  `Stage and commit all changes in the notes directory`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGitCommit,
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push commits to remote",
	Long:  `Push commits to the remote git repository`,
	RunE:  runGitPush,
}

var gitPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from remote",
	Long:  `Pull changes from the remote git repository`,
	RunE:  runGitPull,
}

var gitSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync with remote",
	Long:  `Commit local changes, pull from remote, and push to remote`,
	RunE:  runGitSync,
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show git log",
	Long:  `Show the git commit history`,
	RunE:  runGitLog,
}

var (
	syncMessage string
	logLimit    int
)

func init() {
	rootCmd.AddCommand(gitCmd)

	gitCmd.AddCommand(gitInitCmd)
	gitCmd.AddCommand(gitStatusCmd)
	gitCmd.AddCommand(gitCommitCmd)
	gitCmd.AddCommand(gitPushCmd)
	gitCmd.AddCommand(gitPullCmd)
	gitCmd.AddCommand(gitSyncCmd)
	gitCmd.AddCommand(gitLogCmd)

	gitSyncCmd.Flags().StringVarP(&syncMessage, "message", "m", "", "commit message for sync")
	gitLogCmd.Flags().IntVarP(&logLimit, "limit", "n", 10, "number of commits to show")
}

func runGitInit(cmd *cobra.Command, args []string) error {
	if err := git.Init(); err != nil {
		return err
	}

	fmt.Println("Initialized git repository in notes directory")
	return nil
}

func runGitStatus(cmd *cobra.Command, args []string) error {
	status, err := git.StatusShort()
	if err != nil {
		return err
	}

	if status == "" {
		fmt.Println("No changes (working tree clean)")
		return nil
	}

	fmt.Println(status)
	return nil
}

func runGitCommit(cmd *cobra.Command, args []string) error {
	message := args[0]

	if err := git.Commit(message); err != nil {
		return err
	}

	fmt.Printf("Committed changes: %s\n", message)
	return nil
}

func runGitPush(cmd *cobra.Command, args []string) error {
	hasRemote, err := git.HasRemote()
	if err != nil {
		return err
	}

	if !hasRemote {
		return fmt.Errorf("no remote repository configured")
	}

	if err := git.Push(); err != nil {
		return err
	}

	fmt.Println("Pushed commits to remote")
	return nil
}

func runGitPull(cmd *cobra.Command, args []string) error {
	hasRemote, err := git.HasRemote()
	if err != nil {
		return err
	}

	if !hasRemote {
		return fmt.Errorf("no remote repository configured")
	}

	if err := git.Pull(); err != nil {
		return err
	}

	fmt.Println("Pulled changes from remote")
	return nil
}

func runGitSync(cmd *cobra.Command, args []string) error {
	hasRemote, err := git.HasRemote()
	if err != nil {
		return err
	}

	if !hasRemote {
		return fmt.Errorf("no remote repository configured")
	}

	fmt.Println("Syncing with remote...")

	if err := git.Sync(syncMessage); err != nil {
		return err
	}

	fmt.Println("✓ Sync completed successfully")
	return nil
}

func runGitLog(cmd *cobra.Command, args []string) error {
	log, err := git.Log(logLimit)
	if err != nil {
		return err
	}

	if log == "" {
		fmt.Println("No commits yet")
		return nil
	}

	fmt.Print(log)
	return nil
}
