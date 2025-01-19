package portable

import (
	"log"
	"os"
	"path/filepath"
)

func GetExecutableDirectory() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exe)
}

// Used for dropping all the scripts, configs and models on request
func dropFileAtPath(path string, contents []byte) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("DEV Error: must check if stuff exists before writing")
	}

	err := os.WriteFile(path, contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func IsPortable() bool {
	if _, err := os.Stat(filepath.Join(GetExecutableDirectory(), "portable")); err == nil {
		return true
	}
	return false
}

func GetUserHome() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDirectory
}

func GetConfigDirectory() string {
	if IsPortable() {
		return GetExecutableDirectory()
	}

	configDirectory := filepath.Join(GetUserHome(), ".config", "smoothie-go")
	if err := os.MkdirAll(configDirectory, 0755); err != nil {
		log.Fatal(err)
	}
	return configDirectory
}

func GetLocalDirectory() string {
	localDirectory := filepath.Join(GetUserHome(), ".local", "share", "smoothie-go")
	if err := os.MkdirAll(localDirectory, 0755); err != nil {
		log.Fatal(err)
	}
	return localDirectory
}

func GetRecipePathCustom(name string) string {
	return filepath.Join(GetConfigDirectory(), name)
}

func GetRecipePath() string {
	recipePath := GetRecipePathCustom("recipe.ini")
	if _, err := os.Stat(recipePath); os.IsNotExist(err) {
		dropFileAtPath(recipePath, []byte(recipe_ini))
	}
	return recipePath
}

func GetEncodingPresetsPath() string {
	encodingPresetsPath := filepath.Join(GetConfigDirectory(), "encoding_presets.ini")
	if _, err := os.Stat(encodingPresetsPath); os.IsNotExist(err) {
		dropFileAtPath(encodingPresetsPath, []byte(encoding_resets_ini))
	}
	return encodingPresetsPath
}

func GetModelsPath() string {
	modelsPath := filepath.Join(GetLocalDirectory(), "models")
	if _, err := os.Stat(modelsPath); os.IsNotExist(err) {
		err := os.MkdirAll(modelsPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
		writeEmbeddedFiles(models, modelsPath, "assets/models")
	}
	return modelsPath
}
