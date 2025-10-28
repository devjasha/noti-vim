package main

import (
	"encoding/json"
	"fmt"

	"github.com/devjasha/noti-vim/internal/config"
	"github.com/devjasha/noti-vim/internal/notes"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <title>",
	Short: "Create a new note",
	Long:  `Create a new note with the specified title`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNew,
}

var (
	newFolder string
	newTags   []string
)

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&newFolder, "folder", "f", "", "folder for the new note")
	newCmd.Flags().StringSliceVarP(&newTags, "tags", "t", []string{}, "tags for the new note")
}

func runNew(cmd *cobra.Command, args []string) error {
	title := args[0]

	// Use default folder/tags from config if not specified
	cfg := config.Get()
	if newFolder == "" {
		newFolder = cfg.DefaultFolder
	}
	if len(newTags) == 0 && len(cfg.DefaultTags) > 0 {
		newTags = cfg.DefaultTags
	}

	note, err := notes.CreateNote(title, newFolder, newTags)
	if err != nil {
		return fmt.Errorf("could not create note: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	if jsonOutput {
		data, err := json.MarshalIndent(note, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		fmt.Println(note.FilePath)
		return nil
	}

	fmt.Printf("Created note: %s\n", note.Title)
	fmt.Printf("  slug: %s\n", note.Slug)
	fmt.Printf("  path: %s\n", note.FilePath)

	return nil
}
