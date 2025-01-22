package cmd

import (
	"encoding/json"
	"log"
	"strings"

	"path/filepath"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
	"github.com/smoothie-go/smoothie-go/recipe"
)

func VspipeCommandBuilder(args *cli.Arguments, rc *recipe.Recipe) ([]string, []string) {
	//look for Vspipe
	vspipe := portable.GetBinaryInPathOrBinPath("vspipe")
	if vspipe == "" {
		log.Panicln("Vspipe not found")
	}

	argsjson, _ := json.Marshal(args)
	rcjson, _ := json.Marshal(rc)

	if args.Verbose {
		log.Println(string(rcjson))
		log.Println(string(argsjson))
	}

	vspipeCmd := []string{
		vspipe,
		"--container", "y4m", "-",
		portable.GetMainVpyPath(),
		"--arg", "rec=" + string(rcjson),
		"--arg", "args=" + string(argsjson),
	}

	ffmpeg := portable.GetBinaryInPathOrBinPath("ffmpeg")
	if ffmpeg == "" {
		log.Panicln("FFmpeg not found")
	}

	encArgs := strings.Split(rc.Output.EncArgs, " ")
	ffArgs := strings.Split(rc.Miscellaneous.FfmpegOptions, " ")

	ffmpegCmd := []string{
		ffmpeg,
	}

	for _, arg := range ffArgs {
		ffmpegCmd = append(ffmpegCmd, arg)
	}

	for _, arg := range encArgs {
		ffmpegCmd = append(ffmpegCmd, arg)
	}

	ffmpegCmd = append(ffmpegCmd, filepath.Join(args.OutDir, args.OutputFile)+strings.ToLower(rc.Output.Container))

	return vspipeCmd, ffmpegCmd
}
