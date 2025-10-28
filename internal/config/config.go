package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NotesDir      string   `yaml:"notes_dir"`
	DefaultFolder string   `yaml:"default_folder"`
	DefaultTags   []string `yaml:"default_tags"`
	Editor        string   `yaml:"editor"`
	GitAutoCommit bool     `yaml:"git_auto_commit"`
	GitAutoPush   bool     `yaml:"git_auto_push"`
}

var current *Config

// Load loads the configuration from the specified file or default location
func Load(cfgFile string) error {
	if cfgFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get home directory: %w", err)
		}
		cfgFile = filepath.Join(homeDir, ".config", "noti", "config.yaml")
	}

	// Create default config if file doesn't exist
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if err := createDefaultConfig(cfgFile); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Set defaults
	if cfg.Editor == "" {
		cfg.Editor = os.Getenv("EDITOR")
		if cfg.Editor == "" {
			cfg.Editor = "vim"
		}
	}

	current = &cfg
	return nil
}

// Get returns the current configuration
func Get() *Config {
	if current == nil {
		// Return sensible defaults if config not loaded
		homeDir, _ := os.UserHomeDir()
		return &Config{
			NotesDir:      filepath.Join(homeDir, "notes"),
			DefaultFolder: "",
			DefaultTags:   []string{},
			Editor:        "vim",
			GitAutoCommit: false,
			GitAutoPush:   false,
		}
	}
	return current
}

// Save saves the current configuration
func Save() error {
	if current == nil {
		return fmt.Errorf("no configuration loaded")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}

	cfgFile := filepath.Join(homeDir, ".config", "noti", "config.yaml")
	cfgDir := filepath.Dir(cfgFile)

	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := yaml.Marshal(current)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(cfgFile string) error {
	cfgDir := filepath.Dir(cfgFile)
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get home directory: %w", err)
	}

	defaultCfg := Config{
		NotesDir:      filepath.Join(homeDir, "notes"),
		DefaultFolder: "",
		DefaultTags:   []string{},
		Editor:        os.Getenv("EDITOR"),
		GitAutoCommit: false,
		GitAutoPush:   false,
	}

	if defaultCfg.Editor == "" {
		defaultCfg.Editor = "vim"
	}

	data, err := yaml.Marshal(defaultCfg)
	if err != nil {
		return fmt.Errorf("could not marshal default config: %w", err)
	}

	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
