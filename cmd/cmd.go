package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"path/filepath"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/temp"
)

func ExtractAudioCommandBuilder(args *cli.Arguments, rc *recipe.Recipe, outfile string) []string {
	ffmpeg := portable.GetBinaryInPathOrBinPath("ffmpeg")
	if ffmpeg == "" {
		log.Panicln("FFmpeg not found")
	}

	tempo := rc.Timescale.Out / rc.Timescale.In
	tempoStr := fmt.Sprintf(`[0:a]atempo=%f`, tempo)
	extAudioCmd := []string{
		ffmpeg, "-loglevel", "error", "-i", args.InputFile, "-filter_complex",
		tempoStr, "-map", "0:a", "-c:a", "flac",
		"-compression_level", "0", outfile,
	}

	return extAudioCmd
}

func VspipeCommandBuilder(args *cli.Arguments, rc *recipe.Recipe) ([]string, []string, []string) {
	//look for Vspipe
	vspipe := portable.GetBinaryInPathOrBinPath("vspipe")
	if vspipe == "" {
		log.Panicln("Vspipe not found")
	}

	ffmpeg := portable.GetBinaryInPathOrBinPath("ffmpeg")
	if ffmpeg == "" {
		log.Panicln("FFmpeg not found")
	}

	ffplay := portable.GetBinaryInPathOrBinPath("ffplay")
	if ffplay == "" {
		log.Panicln("FFplay not found")
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

	encArgs := strings.Split(rc.Output.EncArgs, " ")
	ffArgs := strings.Split(rc.Miscellaneous.FfmpegOptions, " ")
	ffplayArgs := strings.Split(rc.Miscellaneous.FfplayOptions, " ")

	audioTracks, err := temp.Join("audiotracks.mka")
	if err != nil {
		log.Panicf("Unable to get Audio Tracks: %v\n", err)
	}

	ffmpegCmd := []string{ffmpeg, "-i", audioTracks}

	tScale := math.Pow(float64(rc.Timescale.Out/rc.Timescale.In), -1)
	tScaleFilter := fmt.Sprintf("setpts=%f*PTS", tScale)

	var userFilter string
	var filteredEncArgs []string

	for i := 0; i < len(encArgs); i++ {
		if (encArgs[i] == "-filter:v") || (encArgs[i] == "-vf") {
			if i+1 < len(encArgs) {
				userFilter = encArgs[i+1]
				i++
			}
		} else {
			filteredEncArgs = append(filteredEncArgs, encArgs[i])
		}
	}

	combinedFilter := tScaleFilter
	if userFilter != "" {
		combinedFilter = combinedFilter + "," + userFilter
	}

	ffmpegCmd = append(ffmpegCmd, ffArgs...)

	ffmpegCmd = append(ffmpegCmd, "-filter:v", combinedFilter)

	ffmpegCmd = append(ffmpegCmd, filteredEncArgs...)

	fmt.Println(ffmpegCmd)

	ffplayCmd := []string{
		ffplay,
	}

	for _, arg := range ffplayArgs {
		ffplayCmd = append(ffplayCmd, arg)
	}

	ffmpegCmd = append(ffmpegCmd, filepath.Join(args.OutDir, args.OutputFile)+strings.ToLower(rc.Output.Container))

	return vspipeCmd, ffmpegCmd, ffplayCmd
}
