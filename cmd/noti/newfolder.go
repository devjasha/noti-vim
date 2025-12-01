package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devjasha/noti-vim/internal/config"
	"github.com/spf13/cobra"
)

var newFolderCmd = &cobra.Command{
	Use:   "newfolder <folder-path>",
	Short: "Create a new folder for organizing notes",
	Long:  `Create a new folder in the notes directory to organize your notes`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNewFolder,
}

func init() {
	rootCmd.AddCommand(newFolderCmd)
}

type folderResult struct {
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

func runNewFolder(cmd *cobra.Command, args []string) error {
	folderPath := args[0]

	cfg := config.Get()
	fullPath := filepath.Join(cfg.NotesDir, folderPath)

	// Create the folder
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("could not create folder: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	result := folderResult{
		Path:     folderPath,
		FullPath: fullPath,
	}

	if jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		fmt.Println(fullPath)
		return nil
	}

	fmt.Printf("Created folder: %s\n", folderPath)
	fmt.Printf("  path: %s\n", fullPath)

	return nil
}
