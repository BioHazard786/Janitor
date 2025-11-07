package organizer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"janitor/internal/config"
)

// Organizer holds the state for an organization run
type Organizer struct {
	config     *config.Config
	basePath   string            // Changed from cwd to basePath
	extToDir   map[string]string // The inverted map: ".png" -> "Pictures"
	selfName   string            // Name of the janitor executable
	configName string            // Name of the config file
}

// New creates a new Organizer
// We've renamed 'cwd' to 'basePath' to be more accurate
func New(cfg *config.Config, basePath string) *Organizer {
	// Build the inverted map for efficient lookups
	extToDir := make(map[string]string)
	for dir, exts := range cfg.Folders {
		for _, ext := range exts {
			// Ensure extension starts with a dot and is lowercase
			extToDir[strings.ToLower(ext)] = dir
		}
	}

	// Get the name of the executable to avoid moving it
	selfName := ""
	exePath, err := os.Executable()
	if err == nil {
		selfName = filepath.Base(exePath)
	}

	return &Organizer{
		config:     cfg,
		basePath:   basePath, // Updated this field
		extToDir:   extToDir,
		selfName:   selfName,
		configName: "config.yml", // This should match the name in config.go
	}
}

// Run executes the file organization
func (o *Organizer) Run() error {
	// Use o.basePath instead of o.cwd
	entries, err := os.ReadDir(o.basePath)
	if err != nil {
		return fmt.Errorf("could not read directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// Skip self and config file
		if fileName == o.selfName || fileName == o.configName {
			continue
		}

		// Get extension and find target directory
		ext := strings.ToLower(filepath.Ext(fileName))
		targetDir, found := o.extToDir[ext]
		if !found {
			// No mapping for this file, so we skip it
			continue
		}

		// File needs to be moved
		if err := o.moveFile(fileName, targetDir); err != nil {
			log.Printf("Failed to move %s: %v", fileName, err)
			// Continue to the next file even if one fails
		}
	}
	return nil
}

// moveFile handles the logic of creating the target dir and moving the file
func (o *Organizer) moveFile(fileName, targetDir string) error {
	// Create the full destination directory path (e.g., /path/to/target/Videos)
	destDir := filepath.Join(o.basePath, targetDir)

	// 1. Ensure target directory exists
	// os.MkdirAll is safe to call even if the directory already exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", targetDir, err)
	}

	// 2. Move the file
	// Old path is in the base directory (e.g., /path/to/target/my_video.mp4)
	oldPath := filepath.Join(o.basePath, fileName)
	// New path is in the destination directory (e.g., /path/to/target/Videos/my_video.mp4)
	newPath := filepath.Join(destDir, fileName)

	// Check for conflicts before moving
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists at destination: %s", newPath)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("could not move file: %w", err)
	}

	log.Printf("Moved '%s' to '%s/'", fileName, targetDir)
	return nil
}
