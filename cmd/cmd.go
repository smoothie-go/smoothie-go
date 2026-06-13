package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"path/filepath"

	"github.com/fatih/color"
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
		tempoStr, "-map", "0:a?", "-c:a", "flac",
		"-compression_level", "0", outfile,
	}

	return extAudioCmd
}

func VspipeCommandBuilder(args *cli.Arguments, rc *recipe.Recipe, hasAudioTracks bool) ([]string, []string, []string) {
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
		printVerboseConfig(args, rc)
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
	ffmpegCmd := []string{ffmpeg}
	if hasAudioTracks {
		audioTracks, err := temp.Join("audiotracks.mka")
		if err != nil {
			log.Panicf("Unable to get Audio Tracks: %v\n", err)
		}
		ffmpegCmd = append(ffmpegCmd, "-i", audioTracks)
	}

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

	ffplayCmd := []string{
		ffplay,
	}

	for _, arg := range ffplayArgs {
		ffplayCmd = append(ffplayCmd, arg)
	}

	ffmpegCmd = append(ffmpegCmd, filepath.Join(args.OutDir, args.OutputFile)+strings.ToLower(rc.Output.Container))

	if args.CEP {
		dim := color.New(color.FgHiBlack).SprintFunc()
		cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
		greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		white := color.New(color.FgWhite).SprintFunc()

		fmt.Printf("%s %s %s\n", dim("┌───"), cyanBold("Command Execution Pipeline"), dim(" ──────────────────────────────────────"))
		fmt.Printf("%s %s %s %s\n", dim("│"), greenBold("Vspipe:"), cyan(vspipeCmd[0]), white(strings.Join(vspipeCmd[1:], " ")))
		fmt.Printf("%s %s %s %s\n", dim("│"), greenBold("FFmpeg:"), cyan(ffmpegCmd[0]), white(strings.Join(ffmpegCmd[1:], " ")))
		if rc.PreviewWindow.Enabled {
			fmt.Printf("%s %s %s %s\n", dim("│"), greenBold("FFplay:"), cyan(ffplayCmd[0]), white(strings.Join(ffplayCmd[1:], " ")))
		}
		fmt.Printf("%s\n", dim("└─────────────────────────────────────────────────────────────────────"))
	}

	return vspipeCmd, ffmpegCmd, ffplayCmd
}

func printVerboseConfig(args *cli.Arguments, rc *recipe.Recipe) {
	cyanBold := color.New(color.FgCyan, color.Bold).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	greenBold := color.New(color.FgGreen, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	purple := color.New(color.FgMagenta).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s %s %s\n", dim("┌───"), cyanBold("smoothie-go Config"), dim(" ───────────────────────────────────────────────"))
	fmt.Printf("%s %s %s %s%s%s\n", dim("│"), bold("Input File  :"), blue(args.InputFile), dim("("), yellow(fmt.Sprintf("%d FPS", args.InputFps)), dim(")"))
	fmt.Printf("%s %s %s %s%s%s\n", dim("│"), bold("Output File :"), blue(args.OutputFile), dim("(Dir: "), yellow(args.OutDir), dim(")"))

	if rc.Interpolation.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Interpolation:"), greenBold("ENABLED"))
		fmt.Printf("%s Type     : %s\n", dim("│   ├──"), cyan(rc.Interpolation.Type))
		fmt.Printf("%s FPS      : %s\n", dim("│   ├──"), cyan(rc.Interpolation.Fps))
		fmt.Printf("%s Speed    : %s\n", dim("│   ├──"), cyan(rc.Interpolation.Speed))
		fmt.Printf("%s Tuning   : %s\n", dim("│   ├──"), cyan(rc.Interpolation.Tuning))
		fmt.Printf("%s GPU      : %s\n", dim("│   └──"), purple(rc.Interpolation.Gpu))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Interpolation:"), dim("DISABLED"))
	}

	if rc.FrameBlending.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Frame Blending:"), greenBold("ENABLED"))
		fmt.Printf("%s FPS      : %s\n", dim("│   ├──"), cyan(rc.FrameBlending.Fps))
		fmt.Printf("%s Intensity: %s\n", dim("│   ├──"), cyan(fmt.Sprintf("%.2f", rc.FrameBlending.Intensity)))
		fmt.Printf("%s Weighting: %s\n", dim("│   ├──"), cyan(rc.FrameBlending.Weighting))
		fmt.Printf("%s Bright   : %s\n", dim("│   └──"), purple(rc.FrameBlending.BrightBlend))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Frame Blending:"), dim("DISABLED"))
	}

	if rc.FlowBlur.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Flow Blur:"), greenBold("ENABLED"))
		fmt.Printf("%s Masking  : %s\n", dim("│   ├──"), purple(rc.FlowBlur.Masking))
		fmt.Printf("%s Amount   : %s\n", dim("│   ├──"), cyan(rc.FlowBlur.Amount))
		fmt.Printf("%s Blending : %s\n", dim("│   └──"), cyan(rc.FlowBlur.DoBlending))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Flow Blur:"), dim("DISABLED"))
	}

	if rc.PreInterp.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Pre-Interp:"), greenBold("ENABLED"))
		fmt.Printf("%s Model    : %s\n", dim("│   ├──"), cyan(rc.PreInterp.Model))
		fmt.Printf("%s Factor   : %s\n", dim("│   ├──"), cyan(rc.PreInterp.Factor))
		fmt.Printf("%s Scene Chg: %s\n", dim("│   ├──"), purple(rc.PreInterp.SceneChange))
		fmt.Printf("%s UHD/TTA  : UHD=%s, TTA=%s\n", dim("│   └──"), purple(rc.PreInterp.Uhd), purple(rc.PreInterp.Tta))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Pre-Interp:"), dim("DISABLED"))
	}

	if rc.ColourGrading.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Color Grading:"), greenBold("ENABLED"))
		fmt.Printf("%s Brightness: %s\n", dim("│   ├──"), cyan(fmt.Sprintf("%.2f", rc.ColourGrading.Brightness)))
		fmt.Printf("%s Saturation: %s\n", dim("│   ├──"), cyan(fmt.Sprintf("%.2f", rc.ColourGrading.Saturation)))
		fmt.Printf("%s Contrast  : %s\n", dim("│   ├──"), cyan(fmt.Sprintf("%.2f", rc.ColourGrading.Contrast)))
		fmt.Printf("%s Hue/Coring: Hue=%s, Coring=%s\n", dim("│   └──"), cyan(fmt.Sprintf("%.2f", rc.ColourGrading.Hue)), cyan(fmt.Sprintf("%.2f", rc.ColourGrading.Coring)))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Color Grading:"), dim("DISABLED"))
	}

	if rc.Lut.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("LUT:"), greenBold("ENABLED"))
		fmt.Printf("%s Path     : %s\n", dim("│   ├──"), blue(rc.Lut.Path))
		fmt.Printf("%s Opacity  : %s\n", dim("│   └──"), cyan(fmt.Sprintf("%.2f", rc.Lut.Opacity)))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("LUT:"), dim("DISABLED"))
	}

	if rc.ArtifactMasking.Enabled {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Artifact Masking:"), greenBold("ENABLED"))
		fmt.Printf("%s Folder   : %s\n", dim("│   ├──"), blue(rc.ArtifactMasking.FolderPath))
		fmt.Printf("%s File     : %s\n", dim("│   └──"), blue(rc.ArtifactMasking.FileName))
	} else {
		fmt.Printf("%s %s %s\n", dim("├──"), bold("Artifact Masking:"), dim("DISABLED"))
	}

	fmt.Printf("%s %s %s %s%s%s\n", dim("├──"), bold("Output:"), green(rc.Output.Process), dim("("), yellow(rc.Output.Container), dim(")"))
	fmt.Printf("%s Encoder Args: %s\n", dim("│   └──"), cyan(rc.Output.EncArgs))
	fmt.Printf("%s %s Enabled=%s %s%s%s%s%s\n", dim("└──"), bold("Preview Window:"), purple(rc.PreviewWindow.Enabled), dim("(Process="), cyan(rc.PreviewWindow.Process), dim(", Args="), cyan(rc.PreviewWindow.OutputArgs), dim(")"))
	fmt.Printf("%s\n", dim("└─────────────────────────────────────────────────────────────────────"))
}
