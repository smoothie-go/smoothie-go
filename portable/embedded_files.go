package portable

import (
	"embed"
	"log"
	"os"
	"path/filepath"
)

//go:embed assets/recipe.ini
var recipe_ini string

//go:embed assets/defaults.ini
var defaults_ini string

//go:embed assets/encoding_presets.ini
var encoding_resets_ini string

//go:embed assets/models/*
var models embed.FS

// Putting this here as this is the only place where im embedding anything really
func writeEmbeddedFiles(embedFS embed.FS, targetDir string, embedPath string) error {
	entries, err := embedFS.ReadDir(embedPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		fullEmbedPath := filepath.Join(embedPath, entry.Name())
		fullTargetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			err := os.MkdirAll(fullTargetPath, os.ModePerm)
			if err != nil {
				return err
			}

			err = writeEmbeddedFiles(embedFS, fullTargetPath, fullEmbedPath)
			if err != nil {
				return err
			}
		} else {
			data, err := embedFS.ReadFile(fullEmbedPath)
			if err != nil {
				return err
			}

			err = os.WriteFile(fullTargetPath, data, 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
