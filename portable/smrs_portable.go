package portable

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func findUsingPathSmrs() string {
	var smoothiePath string
	var err error
	if runtime.GOOS == "windows" { // it's always portable on Windows
		smoothiePath, err = exec.LookPath("smoothie-rs.exe")
		if err != nil {
			log.Printf("Error finding smoothie-rs.exe: %v", err)
		}
	} else {
		smoothiePath, err = exec.LookPath("smoothie-rs")
		if err != nil {
			log.Printf("Error finding smoothie-rs: %v", err)
		}
		if smoothiePath == "/usr/bin/smoothie-rs" { // Adjust path for Arch package
			smoothiePath = "/opt/smoothie-rs/smoothie-rs"
		}
	}
	return smoothiePath
}

func getTargetSmrs() string {
	smrs := findUsingPathSmrs()
	if smrs == "" {
		return ""
	}
	return filepath.Dir(filepath.Dir(smrs))
}

func isPortableSmrs() bool {
	target := getTargetSmrs()
	if target == "" {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	if _, err := os.Stat(filepath.Join(target, "linux-portable-enable")); err == nil {
		return true
	}
	return false
}

func getConfigDirectorySmrs() string {
	if isPortableSmrs() {
		return filepath.Join(getTargetSmrs(), "smoothie-rs")
	}

	confDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config directory: %v", err)
	}
	return filepath.Join(confDir, "smoothie-rs")
}

func GetRecipeSmrs() (string, error) {
	if findUsingPathSmrs() == "" {
		return "", errors.New("smoothie-rs not found")
	}
	if isPortableSmrs() {
		return filepath.Join(getTargetSmrs(), "recipe.ini"), nil
	}
	return filepath.Join(getConfigDirectorySmrs(), "recipe.ini"), nil
}

func GetEncodingPresetsSmrs() (string, error) {
	if findUsingPathSmrs() == "" {
		return "", errors.New("smoothie-rs not found")
	}

	if isPortableSmrs() {
		return filepath.Join(getTargetSmrs(), "encoding_presets.ini"), nil
	}
	return filepath.Join(getConfigDirectorySmrs(), "encoding_presets.ini"), nil
}
