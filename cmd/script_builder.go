package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Hzqkii/smoothie-go/cli"
	"github.com/Hzqkii/smoothie-go/portable"
	rc "github.com/Hzqkii/smoothie-go/recipe"
)

func getPyBool(value bool) string {
	if value {
		return "True"
	}
	return "False"
}

func BuildScriptAtDirectory(path string, args *cli.Arguments, recipe *rc.Recipe) error {

	ScriptString := "import vapoursynth as vs\n\nfrom vapoursynth import core\n\n"
	/*
	 * Scripts taken from https://github.com/couleur-tweak-tips/smoothie-rs/tree/main/target/scripts
	 */
	if rc.GetBooleanValue(recipe.ColourGrading.Enabled) {
		ScriptString += "import adjust\n"
	}
	if rc.GetBooleanValue(recipe.PreInterp.Enabled) {
		ScriptString += "import helpers\n"
	}
	if rc.GetBooleanValue(recipe.FrameBlending.Enabled) {
		ScriptString += "import weighting\n"
	}
	if recipe.Miscellaneous.DedupThreshold > 0.0 {
		ScriptString += "import filldrops\n"
	}
	if rc.GetBooleanValue(recipe.Interpolation.Enabled) || rc.GetBooleanValue(recipe.FrameBlending.Enabled) {
		ScriptString += "import havsfunc\n\n"
	}
	ScriptString += "verbose = " + getPyBool(args.Verbose) + "\n\n"

	ScriptString += "def verbose(msg: str):\n"
	ScriptString += "    if verbose:\n"
	ScriptString += "        print(\"VERB: \" + msg, file=sys.stderr)\n\n"

	ScriptString += "def eprint(msg: str):\n"
	ScriptString += "    print(\"ERR: \" + msg, file=sys.stderr)\n\n"

	// we are already in the temp directory created for this video
	cachePath := filepath.Base(args.InputFile) + "-bs_cache"

	ScriptString += "try:\n"
	ScriptString += fmt.Sprintf("    clip: vs.VideoNode = core.bs.VideoSource(source=\"%s\", cachemode=3, cachepath=\"%s\", showprogress=False)\n", args.InputFile, cachePath)
	ScriptString += "except vs.Error as e:\n"
	ScriptString += "    eprint(\"Failed loading video with error: \" + str(e))\n"
	ScriptString += "    raise\n\n"

	if rc.GetBooleanValue(recipe.PreInterp.Enabled) {
		modelPath := recipe.PreInterp.Model
		if modelPath == "auto" || modelPath == "default" {
			if rc.GetBooleanValue(recipe.PreInterp.Tta) {
				modelPath = portable.GetDefaultTtaModelPath()
			} else {
				modelPath = portable.GetDefaultModelPath()
			}
		}
		// i know absolute jack about any of this, so this is just taken from smrs again
		ScriptString += "heuristic = yuv_heuristic(clip.width, clip.height)\n" +
			"not_in_heuristic = {}\n" +
			"for key, value in heuristic.items():\n" +
			"    not_in_heuristic[key.replace(\"_in\", \"\")] = value\n\n" +
			"clip = core.resize.Bicubic(clip, format=vs.RGB, **heuristic)\n"

		//masking is for later

		var preinterpFactor int
		fmt.Sscanf(recipe.PreInterp.Factor, "%d", &preinterpFactor)

		ScriptString += fmt.Sprintf("clip = core.rife.RIFE(clip=clip, factor_num=%d, model_path=\"%s\", gpu_id=0, gpu_thread=1, tta=%s, uhd=%s, sc=False)\n\n",
			preinterpFactor,
			modelPath,
			getPyBool(rc.GetBooleanValue(recipe.PreInterp.Tta)),
			getPyBool(rc.GetBooleanValue(recipe.PreInterp.Uhd)))

		ScriptString += "clip = core.resize.Bicubic(clip=clip, format=og_format, **not_in_heuristic)\n"
	}

	fmt.Println(ScriptString)
	return nil
}
