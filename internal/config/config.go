package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigName = "config.yml"
	userConfigDir     = ".config/janitor"
)

// Config holds the structure of the YAML config file
type Config struct {
	// Folders is a map where the key is the directory name (e.g., "Videos")
	// and the value is a slice of extensions (e.g., [".mp4", ".mov"])
	Folders map[string][]string `yaml:"folders"`
}

// defaultMappings provides the out-of-the-box experience
var defaultMappings = Config{
	Folders: map[string][]string{
		"Videos":       {".mp4", ".mov", ".avi", ".mkv", ".wmv", ".webm", ".flv", ".mpg", ".mpeg", ".av1", ".opus", ".ts"},
		"Music":        {".mp3", ".wav", ".flac", ".aac", ".ogg", ".mpa", ".m4a", ".wma", ".midi"},
		"Pictures":     {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"},
		"Documents":    {".pdf", ".docx", ".doc", ".txt", ".md", ".xls", ".xlsx", ".ppt", ".pptx", ".cbz", ".cbr"},
		"Archives":     {".zip", ".tar", ".gz", ".rar", ".7z", ".xz", ".bz2"},
		"Applications": {".exe", ".msi", ".apk", ".dmg", ".deb", ".rpm", ".msi", ".appx", ".msix"},
	},
}

// LoadConfig finds and loads the configuration file.
// It searches in:
// 1. Current directory (./config.yml)
// 2. User config directory (~/.config/janitor/config.yml)
// 3. Returns default config if no file is found.
func LoadConfig() (*Config, error) {
	// 1. Check current directory
	if cfg, err := loadFromFile(defaultConfigName); err == nil {
		fmt.Println("Using config from current directory.")
		return cfg, nil
	}

	// 2. Check user config directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(homeDir, userConfigDir, defaultConfigName)
		if cfg, err := loadFromFile(userConfigPath); err == nil {
			fmt.Println("Using config from user home directory.")
			return cfg, nil
		}
	}

	// 3. Use default config
	fmt.Println("No config file found, using default mappings.")
	return &defaultMappings, nil
}

// loadFromFile attempts to read and parse a config file from a given path.
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
