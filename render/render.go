package render

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/cmd"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/temp"
)

type procResult struct {
	cmd *exec.Cmd
	err error
}

func Render(args *cli.Arguments, rc *recipe.Recipe) {
	if err := temp.InitTemp(args); err != nil {
		log.Panicln(err.Error())
	}

	audioTracks, err := temp.Join("audiotracks.mka")
	if err != nil {
		log.Panicln(err.Error())
	}
	extractAudio := cmd.ExtractAudioCommandBuilder(args, rc, audioTracks)

	extractAudioCmd := exec.Command(extractAudio[0], extractAudio[1:]...)
	extractAudioCmd.Stderr = os.Stderr
	extractAudioCmd.Stdout = os.Stdout
	if err := extractAudioCmd.Run(); err != nil {
		log.Panicf("Extract audio failed: %v", err)
	}

	if err := temp.RegisterTempFile("audiotracks.mka"); err != nil {
		log.Panicln(err.Error())
	}

	vspipe, ffmpeg, ffplay := cmd.VspipeCommandBuilder(args, rc)
	vspipeCmd := exec.Command(vspipe[0], vspipe[1:]...)
	ffmpegCmd := exec.Command(ffmpeg[0], ffmpeg[1:]...)
	vspipeCmd.Stderr = os.Stderr
	ffmpegCmd.Stderr = os.Stderr

	var ffplayCmd *exec.Cmd
	var pipeReader2 *io.PipeReader
	var pipeWriter2 *io.PipeWriter

	if rc.PreviewWindow.Enabled {
		ffplayCmd = exec.Command(ffplay[0], ffplay[1:]...)
		ffplayCmd.Stderr = os.Stderr
		pipeReader2, pipeWriter2 = io.Pipe()
		ffplayCmd.Stdin = pipeReader2
	}

	vspipeOut, err := vspipeCmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create pipe for vspipe: %v", err)
	}

	pipeReader1, pipeWriter1 := io.Pipe()
	go func() {
		defer pipeWriter1.Close()
		if rc.PreviewWindow.Enabled {
			defer pipeWriter2.Close()
			multiWriter := io.MultiWriter(pipeWriter1, pipeWriter2)
			if _, err := io.Copy(multiWriter, vspipeOut); err != nil {
				log.Printf("Error while copying vspipe output: %v", err)
			}
		} else {
			if _, err := io.Copy(pipeWriter1, vspipeOut); err != nil {
				log.Printf("Error while copying vspipe output: %v", err)
			}
		}
	}()

	ffmpegCmd.Stdin = pipeReader1

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

	killEverything := false
	var procErr procResult
	for i := 0; i < len(commands); i++ {
		procErr = <-resCh
		if procErr.cmd == ffmpegCmd || procErr.err != nil {
			killEverything = true
			break
		}
	}

	// kill everything if ffmpeg is dead or if theres error
	if killEverything {
		for _, c := range commands {
			if c.Process != nil {
				c.Process.Kill()
			}
		}
	}

	wg.Wait()
}
