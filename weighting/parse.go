package weighting

import (
	"log"
	"strconv"
	"strings"

	"github.com/smoothie-go/smoothie-go/cli"
	rc "github.com/smoothie-go/smoothie-go/recipe"
)

func Parse(args *cli.Arguments, recipe *rc.Recipe) {
	videoFps := 0

	if recipe.Interpolation.Enabled {
		videoFps = recipe.Interpolation.Fps
	} else {
		videoFps = args.InputFps
	}

	frameGap := videoFps / recipe.FrameBlending.Fps
	actualWeights := frameGap * int(recipe.FrameBlending.Intensity)

	if actualWeights > 0 {
		if actualWeights%2 == 0 {
			actualWeights++
		}
	}
	weightingStr := strings.ToLower(recipe.FrameBlending.Weighting)
	if recipe.FrameBlending.Weighting[0] == '[' && recipe.FrameBlending.Weighting[len(recipe.FrameBlending.Weighting)-1] == ']' {
		weightingStrArray := strings.Split(weightingStr[1:len(recipe.FrameBlending.Weighting)-1], ",")

		mapping := make([]float64, len(weightingStrArray))

		for i, v := range weightingStrArray {
			mapping[i], _ = strconv.ParseFloat(v, 64)
		}

		args.Weighting = Divide(actualWeights, mapping)
		return
	} else if weightingStr == "ascending" {
		args.Weighting = Ascending(actualWeights)
		return
	} else if weightingStr == "descending" {
		args.Weighting = Descending(actualWeights)
		return
	} else if weightingStr == "equal" {
		args.Weighting = Equal(actualWeights)
		return
	} else if weightingStr == "pyramid" {
		args.Weighting = Pyramid(actualWeights)
		return
	} else if weightingStr == "gaussian" {
		args.Weighting, _ = Gaussian(actualWeights, 0, 1, [2]float64{-1, 1})
		return
	} else if weightingStr == "gaussiansym" {
		args.Weighting, _ = GaussianSym(actualWeights, 1, [2]float64{-1, 1})
		return
	} else if weightingStr == "vegas" {
		args.Weighting = Vegas(videoFps, recipe.FrameBlending.Fps, float64(recipe.FrameBlending.Intensity))
		return
	} else {
		log.Println("WARNING: Unknown weighting type. Defaulting to Equal")
		args.Weighting = Equal(actualWeights)
		return
	}
}
