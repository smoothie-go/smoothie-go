package portable

import (
	"bytes"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/smoothie-go/smoothie-go/migrate"
	"gopkg.in/ini.v1"
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

func GetGPUs() []string {
	var gpus []string
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-CimInstance Win32_VideoController | Select-Object -ExpandProperty Name")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					gpus = append(gpus, line)
				}
			}
		}
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "lspci | grep -E -i 'vga|3d|display'")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					parts := strings.SplitN(line, ": ", 2)
					if len(parts) > 1 {
						gpus = append(gpus, parts[1])
					} else {
						gpus = append(gpus, line)
					}
				}
			}
		}
	}
	return gpus
}

func promptAndConfigureGPU(iniBytes []byte) []byte {
	cfg, err := ini.Load(iniBytes)
	if err != nil {
		log.Println("Error parsing recipe: ", err)
		return iniBytes
	}

	sec := cfg.Section("interpolation")

	fmt.Println("\n=== Smoothie-go: GPU Configuration ===")
	if IsOpenCLAvailable() {
		fmt.Println("OpenCL was detected on your system. GPU acceleration is highly recommended!")
	} else {
		fmt.Println("Warning: OpenCL was not detected. GPU acceleration may fallback to CPU.")
	}

	fmt.Printf("Would you like to enable GPU acceleration? (y/n) [y]: ")
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	useGpu := true
	if input == "n" || input == "no" {
		useGpu = false
	}

	gpuId := 0
	if useGpu {
		gpus := GetGPUs()
		if len(gpus) > 0 {
			fmt.Println("\nDetected GPU(s):")
			for i, gpu := range gpus {
				fmt.Printf("  [%d] %s\n", i, gpu)
			}
			fmt.Println()
		} else {
			fmt.Println("\nNo GPUs automatically detected. Defaulting device ID to 0.")
		}

		fmt.Printf("Enter GPU ID to use [0]: ")
		var idInput string
		fmt.Scanln(&idInput)
		idInput = strings.TrimSpace(idInput)
		if idInput != "" {
			if parsedId, err := strconv.Atoi(idInput); err == nil {
				gpuId = parsedId
			}
		}
	}

	sec.Key("use gpu").SetValue(strconv.FormatBool(useGpu))
	sec.Key("gpu id").SetValue(strconv.Itoa(gpuId))

	var buf bytes.Buffer
	if _, err := cfg.WriteTo(&buf); err != nil {
		log.Println("Error writing recipe: ", err)
		return iniBytes
	}
	return buf.Bytes()
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
					rcBytes := promptAndConfigureGPU([]byte(rc))
					dropFileAtPath(recipePath, rcBytes)
					return recipePath
				} else {
					fmt.Println("Skipping migration...")
				}
			}
		}
		recipeBytes := promptAndConfigureGPU([]byte(recipe_ini))
		dropFileAtPath(recipePath, recipeBytes)
	} else {
		cfg, err := ini.Load(recipePath)
		if err == nil {
			sec := cfg.Section("interpolation")
			if !sec.HasKey("gpu id") {
				data, err := os.ReadFile(recipePath)
				if err == nil {
					updatedData := promptAndConfigureGPU(data)
					dropFileAtPath(recipePath, updatedData)
				}
			}
		}
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

func DropScriptsAtPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	err := writeEmbeddedFiles(scripts, path, "assets/scripts")
	return err
}

func GetMainVpyPath() string {
	var path string
	path = filepath.Join(GetExecutableDirectory(), "scripts", "main.vpy")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := DropScriptsAtPath(filepath.Dir(path))
		if err != nil {
			log.Fatal("Unable to drop vapoursynth scripts at " + filepath.Dir(path))
		}
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
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", base, rand.Uint64()))
}
