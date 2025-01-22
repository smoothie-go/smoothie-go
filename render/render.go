package render

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/cmd"
	"github.com/smoothie-go/smoothie-go/recipe"
)

func Render(args *cli.Arguments, rc *recipe.Recipe) {
	vspipe, ffmpeg, ffplay := cmd.VspipeCommandBuilder(args, rc)
	vspipeCmd := exec.Command(vspipe[0], vspipe[1:]...)
	ffmpegCmd := exec.Command(ffmpeg[0], ffmpeg[1:]...)
	ffplayCmd := exec.Command(ffplay[0], ffplay[1:]...)

	vspipeCmd.Stderr = os.Stderr
	ffmpegCmd.Stderr = os.Stderr
	ffplayCmd.Stderr = os.Stderr

	vspipeOut, err := vspipeCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create pipe for vspipe: %v", err)
	}

	pipeReader1, pipeWriter1 := io.Pipe()
	pipeReader2, pipeWriter2 := io.Pipe()

	go func() {
		defer pipeWriter1.Close()
		defer pipeWriter2.Close()
		multiWriter := io.MultiWriter(pipeWriter1, pipeWriter2)
		if _, err := io.Copy(multiWriter, vspipeOut); err != nil {
			log.Printf("Error while copying vspipe output: %v", err)
		}
	}()

	ffmpegCmd.Stdin = pipeReader1
	if rc.PreviewWindow.Enabled {
		ffplayCmd.Stdin = pipeReader2
	}

	if err := vspipeCmd.Start(); err != nil {
		log.Fatalf("Failed to start vspipe command: %v", err)
	}

	if err := ffmpegCmd.Start(); err != nil {
		log.Fatalf("Failed to start ffmpeg command: %v", err)
	}

	if rc.PreviewWindow.Enabled {
		if err := ffplayCmd.Start(); err != nil {
			log.Fatalf("Failed to start ffplay command: %v", err)
		}
	}

	if err := ffmpegCmd.Wait(); err != nil {
		log.Fatalf("ffmpeg command finished with error: %v", err)
	}

	if rc.PreviewWindow.Enabled {
		if err := ffplayCmd.Process.Kill(); err != nil {
			log.Printf("Failed to kill ffplay process: %v", err)
		}
	}

	if err := vspipeCmd.Wait(); err != nil {
		log.Fatalf("vspipe command finished with error: %v", err)
	}
}
