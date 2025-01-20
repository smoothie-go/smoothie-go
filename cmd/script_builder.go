package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Hzqkii/smoothie-go/cli"
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
	if rc.GetBooleanValue(recipe.FrameBlending.Enabled) {
		ScriptString += "import weighting\n"
	}
	if recipe.Miscellaneous.DedupThreshold > 0.0 {
		ScriptString += "import filldrops\n"
	}
	if rc.GetBooleanValue(recipe.Interpolation.Enabled) {
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
	ScriptString += fmt.Sprintf("    clip: vs.VideoNode = core.bs.VideoSource(source=%s, cachemode=3, cachepath=%s, showprogress=False)\n", args.InputFile, cachePath)
	ScriptString += "except vs.Error as e:\n"
	ScriptString += "    eprint(\"Failed loading video with error: \" + str(e))\n"
	ScriptString += "    raise\n\n"
	fmt.Println(ScriptString)
	return nil
}
