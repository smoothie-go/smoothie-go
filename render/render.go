package render

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/cmd"
	"github.com/smoothie-go/smoothie-go/portable"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/temp"
)

type procResult struct {
	cmd *exec.Cmd
	err error
}

func hasAudioStream(input string) bool {
	ffprobe := portable.GetBinaryInPathOrBinPath("ffprobe")
	if ffprobe == "" {
		log.Panicln("FFprobe not found")
	}

	cmd := exec.Command(ffprobe, "-v", "error", "-select_streams", "a", "-show_entries",
		"stream=index", "-of", "csv=p=0", input)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to probe audio streams: %v", err)
		return false
	}

	return len(output) > 0
}

// prepareAudio checks if an audio stream exists, extracts it, and registers it.
func prepareAudio(args *cli.Arguments, rc *recipe.Recipe) (bool, error) {
	if !hasAudioStream(args.InputFile) {
		return false, nil
	}

	audioTracks, err := temp.Join("audiotracks.mka")
	if err != nil {
		return false, err
	}

	extractAudio := cmd.ExtractAudioCommandBuilder(args, rc, audioTracks)
	extractAudioCmd := exec.Command(extractAudio[0], extractAudio[1:]...)
	extractAudioCmd.Stderr = os.Stderr
	extractAudioCmd.Stdout = os.Stdout

	if err := extractAudioCmd.Run(); err != nil {
		return false, err
	}

	if err := temp.RegisterTempFile("audiotracks.mka"); err != nil {
		return false, err
	}

	return true, nil
}

// setupPipelinePipes connects stdout of vspipe to stdin of ffmpeg, optionally copying to ffplay via a MultiWriter.
func setupPipelinePipes(rc *recipe.Recipe, vspipeCmd, ffmpegCmd, ffplayCmd *exec.Cmd) error {
	vspipeOut, err := vspipeCmd.StdoutPipe()
	if err != nil {
		return err
	}

	if rc.PreviewWindow.Enabled && ffplayCmd != nil {
		pipeReader1, pipeWriter1 := io.Pipe()
		ffmpegCmd.Stdin = pipeReader1

		pipeReader2, pipeWriter2 := io.Pipe()
		ffplayCmd.Stdin = pipeReader2

		go func() {
			defer pipeWriter1.Close()
			defer pipeWriter2.Close()
			multiWriter := io.MultiWriter(pipeWriter1, pipeWriter2)
			buf := make([]byte, 1024*1024)
			if _, err := io.CopyBuffer(multiWriter, vspipeOut, buf); err != nil {
				log.Printf("Error while copying vspipe output: %v", err)
			}
		}()
	} else {
		ffmpegCmd.Stdin = vspipeOut
	}

	return nil
}

// Render executes the frame rendering pipeline.
func Render(args *cli.Arguments, rc *recipe.Recipe) {
	if err := temp.InitTemp(args); err != nil {
		log.Panicln(err.Error())
	}

	hasAudioTracks, err := prepareAudio(args, rc)
	if err != nil {
		log.Panicf("Prepare audio failed: %v", err)
	}

	vspipe, ffmpeg, ffplay := cmd.VspipeCommandBuilder(args, rc, hasAudioTracks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vspipeCmd := exec.CommandContext(ctx, vspipe[0], vspipe[1:]...)
	ffmpegCmd := exec.CommandContext(ctx, ffmpeg[0], ffmpeg[1:]...)
	vspipeCmd.Stderr = os.Stderr
	ffmpegCmd.Stderr = os.Stderr

	var ffplayCmd *exec.Cmd
	if rc.PreviewWindow.Enabled {
		ffplayCmd = exec.CommandContext(ctx, ffplay[0], ffplay[1:]...)
		ffplayCmd.Stderr = os.Stderr
	}

	if err := setupPipelinePipes(rc, vspipeCmd, ffmpegCmd, ffplayCmd); err != nil {
		log.Panicf("Failed to setup pipeline pipes: %v", err)
	}

	commands := []*exec.Cmd{vspipeCmd, ffmpegCmd}
	if rc.PreviewWindow.Enabled {
		commands = append(commands, ffplayCmd)
	}

	for _, c := range commands {
		if err := c.Start(); err != nil {
			log.Panicf("Failed to start %s: %v", c.Path, err)
		}
	}

	resCh := make(chan procResult, len(commands))
	var wg sync.WaitGroup
	wg.Add(len(commands))

	for _, c := range commands {
		go func(cmd *exec.Cmd) {
			defer wg.Done()
			err := cmd.Wait()
			resCh <- procResult{cmd: cmd, err: err}
		}(c)
	}

	// Wait for any command to finish. If the main encoder (ffmpegCmd) exits,
	// or if any command exits with an error, cancel the context to clean up the rest.
	for i := 0; i < len(commands); i++ {
		res := <-resCh
		if res.cmd == ffmpegCmd || res.err != nil {
			cancel()
		}
	}

	wg.Wait()
}
