package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devjasha/noti-vim/internal/config"
	"github.com/devjasha/noti-vim/pkg/frontmatter"
)

// Note represents a markdown note
type Note struct {
	Slug     string    `json:"slug"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Tags     []string  `json:"tags"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Folder   string    `json:"folder"`
	FilePath string    `json:"file_path"`
}

// ParseNote reads and parses a note from a file
func ParseNote(path string) (*Note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	fm, content, err := frontmatter.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("could not parse frontmatter: %w", err)
	}

	// Get file info for modified time
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not stat file: %w", err)
	}

	// Calculate slug and folder from path
	cfg := config.Get()
	relPath, err := filepath.Rel(cfg.NotesDir, path)
	if err != nil {
		return nil, fmt.Errorf("could not get relative path: %w", err)
	}

	// Remove .md extension
	slug := strings.TrimSuffix(relPath, ".md")
	// Convert path separator to forward slash
	slug = filepath.ToSlash(slug)

	// Extract folder (everything except filename)
	folder := filepath.ToSlash(filepath.Dir(relPath))
	if folder == "." {
		folder = ""
	}

	note := &Note{
		Slug:     slug,
		Title:    fm.Title,
		Content:  content,
		Tags:     fm.Tags,
		Created:  fm.Created,
		Modified: fileInfo.ModTime(),
		Folder:   folder,
		FilePath: path,
	}

	return note, nil
}

// SaveNote saves a note to disk
func SaveNote(note *Note) error {
	cfg := config.Get()

	// Build full file path
	fullPath := filepath.Join(cfg.NotesDir, note.Slug+".md")

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	// Format frontmatter and content
	fm := &frontmatter.Frontmatter{
		Title:   note.Title,
		Tags:    note.Tags,
		Created: note.Created,
	}

	data, err := frontmatter.Format(fm, note.Content)
	if err != nil {
		return fmt.Errorf("could not format note: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

// ListNotes returns all notes, optionally filtered by folder and/or tag
func ListNotes(folder, tag string) ([]*Note, error) {
	cfg := config.Get()

	var notes []*Note

	err := filepath.Walk(cfg.NotesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-markdown files
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Skip hidden files and .templates directory
		if strings.HasPrefix(info.Name(), ".") || strings.Contains(path, "/.templates/") {
			return nil
		}

		// Parse note
		note, err := ParseNote(path)
		if err != nil {
			// Skip files that can't be parsed
			return nil
		}

		// Apply filters
		if folder != "" && note.Folder != folder {
			return nil
		}

		if tag != "" {
			hasTag := false
			for _, t := range note.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				return nil
			}
		}

		notes = append(notes, note)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not walk notes directory: %w", err)
	}

	return notes, nil
}

// GetNote retrieves a single note by slug
func GetNote(slug string) (*Note, error) {
	cfg := config.Get()
	path := filepath.Join(cfg.NotesDir, slug+".md")

	return ParseNote(path)
}

// DeleteNote deletes a note by slug
func DeleteNote(slug string) error {
	cfg := config.Get()
	path := filepath.Join(cfg.NotesDir, slug+".md")

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("could not delete note: %w", err)
	}

	return nil
}

// CreateNote creates a new note with the given title and optional parameters
func CreateNote(title string, folder string, tags []string) (*Note, error) {
	// Generate slug from title
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Add folder to slug if specified
	if folder != "" {
		slug = folder + "/" + slug
	}

	note := &Note{
		Slug:     slug,
		Title:    title,
		Content:  "",
		Tags:     tags,
		Created:  time.Now(),
		Modified: time.Now(),
		Folder:   folder,
	}

	// Set filepath
	cfg := config.Get()
	note.FilePath = filepath.Join(cfg.NotesDir, slug+".md")

	// Save note
	if err := SaveNote(note); err != nil {
		return nil, err
	}

	return note, nil
}
