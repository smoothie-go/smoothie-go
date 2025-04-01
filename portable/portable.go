package portable

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/smoothie-go/smoothie-go/migrate"
)

func GetExecutableDirectory() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(exe)
}

func dropFileAtPath(path string, contents []byte) {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		log.Fatal("DEV Error: must check if stuff exists before writing")
	}

	err := os.WriteFile(path, contents, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func IsPortable() bool {
	if _, err := os.Stat(filepath.Join(GetExecutableDirectory(), "smoothie-go.portable")); err == nil {
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

	osConfigDirctory, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirectory := filepath.Join(osConfigDirctory, "smoothie-go")
	if err := os.MkdirAll(configDirectory, 0755); err != nil {
		log.Fatal(err)
	}
	return configDirectory
}

func GetLocalDirectory() string {
	if IsPortable() {
		return GetExecutableDirectory()
	}
	if runtime.GOOS == "windows" {
		return GetConfigDirectory()
	}
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
		smrsRecipe, err := GetRecipeSmrs()
		if err == nil {
			if _, err := os.Stat(smrsRecipe); err == nil {
				fmt.Printf("Would you like to migrate your recipe? (y/n) ")
				var input string
				fmt.Scanln(&input)
				if input == "y" {
					fmt.Println("Migrating recipe...")
					rcPath := smrsRecipe
					rc, err := migrate.Migrate(rcPath)
					if err != nil {
						log.Fatal(err)
					}
					dropFileAtPath(recipePath, []byte(rc))
					return recipePath
				} else {
					fmt.Println("Skipping migration...")
				}
			}
		}
		dropFileAtPath(recipePath, []byte(recipe_ini))
	}
	return recipePath
}

func GetDefaultRecipePath() string {
	defaultsPath := filepath.Join(GetConfigDirectory(), "defaults.ini")

	if _, err := os.Stat(defaultsPath); os.IsNotExist(err) {
		dropFileAtPath(defaultsPath, []byte(defaults_ini))
	}
	return defaultsPath
}

func GetEncodingPresetsPath() string {
	encodingPresetsPath := filepath.Join(GetConfigDirectory(), "encoding_presets.ini")
	if _, err := os.Stat(encodingPresetsPath); os.IsNotExist(err) {
		smrsEncPresets, err := GetEncodingPresetsSmrs()
		if err == nil {
			if _, err := os.Stat(smrsEncPresets); err == nil {
				fmt.Printf("Would you like to migrate your encoding presets? (y/n) ")
				var input string
				fmt.Scanln(&input)
				if input == "y" {
					fmt.Println("Migrating encoding presets...")
					epPath := smrsEncPresets
					ep, err := migrate.Migrate(epPath)
					if err != nil {
						log.Fatal(err)
					}
					dropFileAtPath(encodingPresetsPath, []byte(ep))
					return encodingPresetsPath
				} else {
					fmt.Println("Skipping migration...")
				}
			}
		}
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

func GetBinaryInPathOrBinPath(binary string) string {
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	execInPath, _ := exec.LookPath(binary)
	if _, err := os.Stat(filepath.Join(GetExecutableDirectory(), binary)); err == nil {
		return filepath.Join(GetExecutableDirectory(), binary)
	} else if _, err := os.Stat(execInPath); err == nil {
		return execInPath
	} else {
		return ""
	}
}

func GetDefaultModelPath() string {
	return filepath.Join(GetModelsPath(), "rife-v4.6/")
}

func GetDefaultTtaModelPath() string {
	return filepath.Join(GetModelsPath(), "rife-v3.1/")
}

func DropScriptsAtPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	writeEmbeddedFiles(scripts, path, "assets/scripts")
}

func GetMainVpyPath() string {
	var path string
	if IsPortable() {
		path = filepath.Join(GetExecutableDirectory(), "scripts", "main.vpy")
	} else {
		path = filepath.Join(GetLocalDirectory(), "scripts", "main.vpy")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		DropScriptsAtPath(filepath.Dir(path))
	}

	return path
}

func GetLogPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(GetConfigDirectory(), "smoothie-go.log")
	}
	logPath := filepath.Join(GetUserHome(), ".local", "state", "smoothie-go.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(logPath), 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return logPath
}

func GetTempPath(inputFile string) string {
	return filepath.Join(os.TempDir(), inputFile+strconv.Itoa(rand.Intn(10000000000)))
}
