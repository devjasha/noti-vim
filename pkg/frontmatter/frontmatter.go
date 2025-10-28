package frontmatter

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter of a markdown note
type Frontmatter struct {
	Title   string    `yaml:"title"`
	Tags    []string  `yaml:"tags"`
	Created time.Time `yaml:"created"`
}

// Parse parses a markdown file with YAML frontmatter
// Returns the frontmatter and the content separately
func Parse(data []byte) (*Frontmatter, string, error) {
	// Check if file starts with frontmatter delimiter
	if !bytes.HasPrefix(data, []byte("---\n")) && !bytes.HasPrefix(data, []byte("---\r\n")) {
		return nil, string(data), nil
	}

	// Find the end of frontmatter
	delimiter := []byte("---")
	parts := bytes.SplitN(data, delimiter, 3)
	if len(parts) < 3 {
		return nil, string(data), nil
	}

	// Parse YAML frontmatter
	var fm Frontmatter
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, "", fmt.Errorf("could not parse frontmatter: %w", err)
	}

	// Get content (everything after second ---)
	content := strings.TrimSpace(string(parts[2]))

	return &fm, content, nil
}

// Format formats frontmatter and content into a complete markdown file
func Format(fm *Frontmatter, content string) ([]byte, error) {
	// Marshal frontmatter to YAML
	fmData, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("could not marshal frontmatter: %w", err)
	}

	// Build complete file
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(fmData)
	buf.WriteString("---\n\n")
	buf.WriteString(content)

	return buf.Bytes(), nil
}
