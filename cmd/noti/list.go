package main

import (
	"encoding/json"
	"fmt"

	"github.com/devjasha/noti-vim/internal/notes"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Long:  `List all notes, optionally filtered by folder or tag`,
	RunE:  runList,
}

var (
	listFolder string
	listTag    string
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listFolder, "folder", "f", "", "filter by folder")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "filter by tag")
}

func runList(cmd *cobra.Command, args []string) error {
	notesList, err := notes.ListNotes(listFolder, listTag)
	if err != nil {
		return fmt.Errorf("could not list notes: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	if jsonOutput {
		data, err := json.MarshalIndent(notesList, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		for _, note := range notesList {
			fmt.Println(note.Slug)
		}
		return nil
	}

	// Human-readable output
	if len(notesList) == 0 {
		fmt.Println("No notes found")
		return nil
	}

	fmt.Printf("Found %d note(s):\n\n", len(notesList))
	for _, note := range notesList {
		fmt.Printf("  %s\n", note.Title)
		fmt.Printf("    slug: %s\n", note.Slug)
		if len(note.Tags) > 0 {
			fmt.Printf("    tags: %v\n", note.Tags)
		}
		if note.Folder != "" {
			fmt.Printf("    folder: %s\n", note.Folder)
		}
		fmt.Println()
	}

	return nil
}
