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

	interpEnabled := recipe.Interpolation.Enabled
	frameBlendingEnabled := recipe.FrameBlending.Enabled
	interpFps := recipe.Interpolation.Fps
	frameBlendingFps := recipe.FrameBlending.Fps
	inputFps := args.InputFps

	if interpEnabled && frameBlendingEnabled {
		if interpFps <= inputFps {
			log.Fatal("Interpolation FPS must be higher than input FPS")
		}
		if frameBlendingFps >= interpFps {
			log.Fatal("Frame blending FPS must be lower than interpolation FPS")
		}
	} else if frameBlendingEnabled && frameBlendingFps > inputFps {
		log.Fatal("Frame blending FPS cannot be higher than input FPS")
	}

	return recipe
}
