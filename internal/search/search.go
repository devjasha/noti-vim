package search

import (
	"strings"

	"github.com/devjasha/noti-vim/internal/notes"
)

// SearchResult represents a search match
type SearchResult struct {
	Note    *notes.Note `json:"note"`
	Matches []Match     `json:"matches"`
}

// Match represents a single match within a note
type Match struct {
	LineNumber int    `json:"line_number"`
	Line       string `json:"line"`
	Context    string `json:"context"`
}

// Search performs a full-text search across all notes
func Search(query string) ([]*SearchResult, error) {
	// Get all notes
	allNotes, err := notes.ListNotes("", "")
	if err != nil {
		return nil, err
	}

	var results []*SearchResult
	queryLower := strings.ToLower(query)

	for _, note := range allNotes {
		matches := searchInNote(note, queryLower)
		if len(matches) > 0 {
			results = append(results, &SearchResult{
				Note:    note,
				Matches: matches,
			})
		}
	}

	return results, nil
}

// searchInNote searches for query within a single note
func searchInNote(note *notes.Note, queryLower string) []Match {
	var matches []Match

	// Search in title
	if strings.Contains(strings.ToLower(note.Title), queryLower) {
		matches = append(matches, Match{
			LineNumber: 0,
			Line:       note.Title,
			Context:    "title",
		})
	}

	// Search in content
	lines := strings.Split(note.Content, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), queryLower) {
			// Get context (line before and after if available)
			context := getContext(lines, i)
			matches = append(matches, Match{
				LineNumber: i + 1,
				Line:       line,
				Context:    context,
			})
		}
	}

	// Search in tags
	for _, tag := range note.Tags {
		if strings.Contains(strings.ToLower(tag), queryLower) {
			matches = append(matches, Match{
				LineNumber: 0,
				Line:       tag,
				Context:    "tag",
			})
		}
	}

	return matches
}

// getContext returns surrounding lines for context
func getContext(lines []string, index int) string {
	start := index - 1
	if start < 0 {
		start = 0
	}

	end := index + 2
	if end > len(lines) {
		end = len(lines)
	}

	contextLines := lines[start:end]
	return strings.Join(contextLines, " ")
}

// FindByFilename searches for notes by filename pattern
func FindByFilename(pattern string) ([]*notes.Note, error) {
	allNotes, err := notes.ListNotes("", "")
	if err != nil {
		return nil, err
	}

	var results []*notes.Note
	patternLower := strings.ToLower(pattern)

	for _, note := range allNotes {
		// Check if pattern matches slug or title
		if strings.Contains(strings.ToLower(note.Slug), patternLower) ||
			strings.Contains(strings.ToLower(note.Title), patternLower) {
			results = append(results, note)
		}
	}

	return results, nil
}
