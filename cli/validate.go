package cli

import (
	"fmt"
	"github.com/smoothie-go/smoothie-go/fruits"
	"github.com/smoothie-go/smoothie-go/portable"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ValidateArgs(args *Arguments) *Arguments {
	if args.InputFile == "" {
		log.Fatal("You must provide an input file")
	} else if _, err := os.Stat(args.InputFile); os.IsNotExist(err) {
		log.Fatal("Input file does not exist")
	}

	fps, err := GetFramerate(args.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	args.InputFps = int(fps)

	inputBaseName := filepath.Base(args.InputFile)
	inputDirectory := filepath.Dir(args.InputFile)
	extIndex := strings.LastIndex(inputBaseName, ".")

	if args.OutputFile != "" && args.OutDir == "" {
		args.OutDir = filepath.Dir(args.OutputFile)
		args.OutputFile = filepath.Base(args.OutputFile)
	} else if args.OutputFile == "" && args.OutDir != "" {
		args.OutputFile = fmt.Sprintf("%s ~ %s", inputBaseName[:extIndex], fruits.GetRandomFruit())
	} else {
		outputDirectory := inputDirectory
		args.OutputFile = fmt.Sprintf("%s ~ %s", inputBaseName[:extIndex], fruits.GetRandomFruit())
		args.OutDir = outputDirectory
	}

	if _, err := os.Stat(args.OutDir); os.IsNotExist(err) {
		log.Fatal("Output directory does not exist")
	}

	if args.RecipePath == "" {
		args.RecipePath = portable.GetRecipePath()
	} else {
		args.RecipePath = filepath.Join(portable.GetConfigDirectory(), args.RecipePath)
		if _, err := os.Stat(args.RecipePath); os.IsNotExist(err) {
			log.Fatal("Recipe file does not exist at " + args.RecipePath)
		}
	}

	if args.LogFile == "" {
		args.LogFile = portable.GetLogPath()
	}

	return args
}
