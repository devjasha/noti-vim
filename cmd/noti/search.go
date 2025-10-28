package main

import (
	"encoding/json"
	"fmt"

	"github.com/devjasha/noti-vim/internal/search"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search notes by content",
	Long:  `Search for notes containing the specified query in title, content, or tags`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	results, err := search.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	quietOutput, _ := cmd.Flags().GetBool("quiet")

	if jsonOutput {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("could not marshal JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	if quietOutput {
		for _, result := range results {
			fmt.Println(result.Note.Slug)
		}
		return nil
	}

	// Human-readable output
	if len(results) == 0 {
		fmt.Println("No matches found")
		return nil
	}

	fmt.Printf("Found %d note(s) matching '%s':\n\n", len(results), query)

	for _, result := range results {
		fmt.Printf("📄 %s\n", result.Note.Title)
		fmt.Printf("   slug: %s\n", result.Note.Slug)
		if len(result.Note.Tags) > 0 {
			fmt.Printf("   tags: %v\n", result.Note.Tags)
		}
		fmt.Printf("   %d match(es):\n", len(result.Matches))

		for _, match := range result.Matches {
			if match.Context == "title" {
				fmt.Printf("     • in title: %s\n", match.Line)
			} else if match.Context == "tag" {
				fmt.Printf("     • in tag: %s\n", match.Line)
			} else {
				fmt.Printf("     • line %d: %s\n", match.LineNumber, match.Line)
			}
		}
		fmt.Println()
	}

	return nil
}
