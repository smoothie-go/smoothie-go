package recipe

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/fruits"
	"github.com/smoothie-go/smoothie-go/portable"
	"gopkg.in/ini.v1"
)

func GetOutputFps(args *cli.Arguments, rc *Recipe) int {
	if rc.FrameBlending.Enabled {
		return rc.FrameBlending.Fps
	}
	if rc.Interpolation.Enabled {
		return rc.Interpolation.Fps
	}
	if rc.PreInterp.Enabled {
		var factor int
		fmt.Sscanf(rc.PreInterp.Factor, "%d", &factor)
		if factor <= 0 {
			factor = 1
		}
		return args.InputFps * factor
	}
	return args.InputFps
}

func Parse(args *cli.Arguments) *Recipe {
	defaults, err := ini.Load(portable.GetDefaultRecipePath())
	if err != nil {
		log.Fatal(err)
	}

	recipe := &Recipe{}

	err = defaults.MapTo(recipe)
	if err != nil {
		log.Fatal(err)
	}

	config, err := ini.Load(args.RecipePath)
	if err != nil {
		log.Fatal(err)
	}

	err = config.MapTo(recipe)
	if err != nil {
		log.Fatal(err)
	}

	if recipe.Miscellaneous.AlwaysVerbose {
		args.Verbose = true
	}

	recipe.ini_config = config

	ParseDynamicMasking(recipe, args)

	if !args.UserSpecifiedOutDir && recipe.Miscellaneous.GlobalOutputFolder != "" {
		args.OutDir = recipe.Miscellaneous.GlobalOutputFolder
		if _, err := os.Stat(args.OutDir); os.IsNotExist(err) {
			log.Printf("WARNING: Global output folder %s does not exist. Falling back to default output directory.\n", args.OutDir)
			args.OutDir = filepath.Dir(args.InputFile)
		}
	}

	if !args.UserSpecifiedOutput && recipe.Output.FileFormat != "" {
		inputBaseName := filepath.Base(args.InputFile)
		extIndex := strings.LastIndex(inputBaseName, ".")
		filename := inputBaseName
		if extIndex != -1 {
			filename = inputBaseName[:extIndex]
		}

		formatted := recipe.Output.FileFormat
		formatted = strings.ReplaceAll(formatted, "%FILENAME%", filename)
		formatted = strings.ReplaceAll(formatted, "%FRUIT%", fruits.GetRandomFruit())

		outFps := GetOutputFps(args, recipe)
		formatted = strings.ReplaceAll(formatted, "%OUTPUT_FPS%", fmt.Sprintf("%d", outFps))

		formatted = strings.ReplaceAll(formatted, "%SPEED%", recipe.Interpolation.Speed)
		formatted = strings.ReplaceAll(formatted, "%TUNING%", recipe.Interpolation.Tuning)

		args.OutputFile = formatted
	}

	recipe = Validate(args, recipe)

	return recipe
}
