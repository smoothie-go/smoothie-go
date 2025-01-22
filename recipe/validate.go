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

	// check if interp is higher than input
	if recipe.Interpolation.Enabled && recipe.Interpolation.Fps < args.InputFps {
		log.Fatal("Interpolation fps cannot be lower than input fps")
	} else if recipe.PreInterp.Enabled && recipe.Interpolation.Enabled {
		var factor int
		fmt.Sscanf(recipe.PreInterp.Factor, "%d", &factor)

		if recipe.Interpolation.Fps*factor < (args.InputFps * factor) {
			log.Fatal("Interpolation fps cannot be lower than pre-interped fps")
		}
	}

	if recipe.FrameBlending.Enabled && recipe.FrameBlending.Fps > args.InputFps {
		log.Fatal("Frame blending fps cannot be higher than input fps")
	}

	return recipe
}
