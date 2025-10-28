package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devjasha/noti-vim/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a notes directory",
	Long:  `Initialize a notes directory and create configuration`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	var notesDir string

	if len(args) > 0 {
		notesDir = args[0]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get home directory: %w", err)
		}
		notesDir = filepath.Join(homeDir, "notes")
	}

	// Make absolute path
	absPath, err := filepath.Abs(notesDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	// Update config
	cfg := config.Get()
	cfg.NotesDir = absPath

	if err := config.Save(); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}

	fmt.Printf("Initialized notes directory at: %s\n", absPath)
	fmt.Println("\nConfiguration saved to: ~/.config/noti/config.yaml")
	fmt.Println("\nYou can now create notes with: noti new \"My First Note\"")

	return nil
}
