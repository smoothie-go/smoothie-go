package recipe

import (
	"fmt"
	"log"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
)

func Validate(args *cli.Arguments, recipe *Recipe) *Recipe {
	// enc args
	outEncArgs, err := ParseEncodingArgs(portable.GetEncodingPresetsPath(), recipe.Output.EncArgs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(outEncArgs)

	recipe.Output.EncArgs = outEncArgs

	// model
	if recipe.PreInterp.Model == "" ||
		recipe.PreInterp.Model == "default" ||
		recipe.PreInterp.Model == "auto" {
		if recipe.PreInterp.Tta {
			recipe.PreInterp.Model = portable.GetDefaultTtaModelPath()
		} else {
			recipe.PreInterp.Model = portable.GetDefaultModelPath()
		}
	}
	return recipe
}
