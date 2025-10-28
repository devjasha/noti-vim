package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/devjasha/noti-vim/internal/notes"
	"github.com/spf13/cobra"
)

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List all folders",
	Long:  `List all folders used to organize notes`,
	RunE:  runFolders,
}

var showTree bool

type folderInfo struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

func init() {
	rootCmd.AddCommand(foldersCmd)
	foldersCmd.Flags().BoolVarP(&showTree, "tree", "t", false, "show as tree structure")
}

func runFolders(cmd *cobra.Command, args []string) error {
	allNotes, err := notes.ListNotes("", "")
	if err != nil {
		return fmt.Errorf("could not list notes: %w", err)
	}

	// Collect unique folders
	folderSet := make(map[string]int)
	for _, note := range allNotes {
		if note.Folder != "" {
			folderSet[note.Folder]++

			// Also count parent folders
			parts := strings.Split(note.Folder, "/")
			for i := 1; i < len(parts); i++ {
				parent := strings.Join(parts[:i], "/")
				if parent != "" {
					folderSet[parent]++
				}
			}
		}
	}

	// Convert to sorted slice

	var folders []folderInfo
	for folder, count := range folderSet {
		folders = append(folders, folderInfo{Path: folder, Count: count})
	}

	// Sort alphabetically
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Path < folders[j].Path
	})

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	if jsonOutput {
		data, err := json.MarshalIndent(folders, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		for _, folder := range folders {
			fmt.Println(folder.Path)
		}
		return nil
	}

	// Human-readable output
	if len(folders) == 0 {
		fmt.Println("No folders found (all notes in root)")
		return nil
	}

	fmt.Printf("Found %d folder(s):\n\n", len(folders))

	if showTree {
		printTree(folders)
	} else {
		for _, folder := range folders {
			fmt.Printf("  %-30s (%d note%s)\n", folder.Path, folder.Count, plural(folder.Count))
		}
	}

	return nil
}

func printTree(folders []folderInfo) {
	// Build tree structure
	for _, folder := range folders {
		depth := strings.Count(folder.Path, "/")
		indent := strings.Repeat("  ", depth)
		name := folder.Path
		if depth > 0 {
			parts := strings.Split(folder.Path, "/")
			name = parts[len(parts)-1]
		}
		fmt.Printf("%s📁 %s (%d)\n", indent, name, folder.Count)
	}
}
