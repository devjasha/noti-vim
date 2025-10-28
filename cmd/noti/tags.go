package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/devjasha/noti-vim/internal/notes"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	Long:  `List all tags used across notes with usage counts`,
	RunE:  runTags,
}

var showCounts bool

func init() {
	rootCmd.AddCommand(tagsCmd)
	tagsCmd.Flags().BoolVarP(&showCounts, "count", "c", true, "show usage counts")
}

func runTags(cmd *cobra.Command, args []string) error {
	allNotes, err := notes.ListNotes("", "")
	if err != nil {
		return fmt.Errorf("could not list notes: %w", err)
	}

	// Count tag occurrences
	tagCounts := make(map[string]int)
	for _, note := range allNotes {
		for _, tag := range note.Tags {
			tagCounts[tag]++
		}
	}

	// Convert to sorted slice
	type tagInfo struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}

	var tags []tagInfo
	for tag, count := range tagCounts {
		tags = append(tags, tagInfo{Tag: tag, Count: count})
	}

	// Sort by count (descending), then by name
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].Count != tags[j].Count {
			return tags[i].Count > tags[j].Count
		}
		return tags[i].Tag < tags[j].Tag
	})

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	if jsonOutput {
		data, err := json.MarshalIndent(tags, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		for _, tag := range tags {
			fmt.Println(tag.Tag)
		}
		return nil
	}

	// Human-readable output
	if len(tags) == 0 {
		fmt.Println("No tags found")
		return nil
	}

	fmt.Printf("Found %d tag(s):\n\n", len(tags))
	for _, tag := range tags {
		if showCounts {
			fmt.Printf("  %-20s (%d note%s)\n", tag.Tag, tag.Count, plural(tag.Count))
		} else {
			fmt.Printf("  %s\n", tag.Tag)
		}
	}

	return nil
}

func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
