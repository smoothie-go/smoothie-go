package render

import (
	"log"
	"os"
	"os/exec"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/cmd"
	"github.com/smoothie-go/smoothie-go/recipe"
)

func Render(args *cli.Arguments, rc *recipe.Recipe) {
	vspipe, ffmpeg := cmd.VspipeCommandBuilder(args, rc)
	vspipeCmd := exec.Command(vspipe[0], vspipe[1:]...)
	ffmpegCmd := exec.Command(ffmpeg[0], ffmpeg[1:]...)

	vspipeCmd.Stderr = os.Stderr

	ffmpegCmd.Stderr = os.Stderr

	pipe, err := vspipeCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create pipe for vspipe: %v", err)
	}
	ffmpegCmd.Stdin = pipe

	if err := vspipeCmd.Start(); err != nil {
		log.Fatalf("Failed to start vspipe command: %v", err)
	}

	if err := ffmpegCmd.Start(); err != nil {
		log.Fatalf("Failed to start ffmpeg command: %v", err)
	}

	if err := ffmpegCmd.Wait(); err != nil {
		log.Fatalf("ffmpeg command finished with error: %v", err)
	}

	if err := vspipeCmd.Process.Kill(); err != nil {
		log.Fatalf("Failed to kill vspipe command: %v", err)
	}

	if err := vspipeCmd.Wait(); err != nil {
		log.Fatalf("vspipe command finished with error: %v", err)
	}
}
