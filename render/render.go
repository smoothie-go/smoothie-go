package render

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/cmd"
	"github.com/smoothie-go/smoothie-go/graph"
	"github.com/smoothie-go/smoothie-go/portable"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/temp"
	"github.com/smoothie-go/smoothie-go/vsgo"
)

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

// a closed preview window must not stop the render
type preview struct{ w io.Writer }

func (p *preview) Write(b []byte) (int, error) {
	if p.w != nil {
		if _, err := p.w.Write(b); err != nil {
			p.w = nil
		}
	}
	return len(b), nil
}

// Render executes the frame rendering pipeline.
func Render(args *cli.Arguments, rc *recipe.Recipe) {
	core, err := vsgo.NewCore()
	if err != nil {
		log.Fatal(err)
	}
	core.OnLog(func(level int, msg string) {
		if level >= vsgo.LogWarning || args.Verbose {
			log.Println(msg)
		}
	})
	clip := graph.Build(core, args, rc)

	if err := temp.InitTemp(args); err != nil {
		log.Panicln(err.Error())
	}

	hasAudioTracks, err := prepareAudio(args, rc)
	if err != nil {
		log.Panicf("Prepare audio failed: %v", err)
	}

	ffmpeg, ffplay := cmd.EncodeCommandBuilder(args, rc, hasAudioTracks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var finish []func()
	start := func(command []string) io.Writer {
		c := exec.CommandContext(ctx, command[0], command[1:]...)
		c.Stderr = os.Stderr
		stdin, err := c.StdinPipe()
		if err != nil {
			log.Panicf("Failed to setup pipeline pipes: %v", err)
		}
		if err := c.Start(); err != nil {
			log.Panicf("Failed to start %s: %v", c.Path, err)
		}
		finish = append(finish, func() {
			stdin.Close()
			c.Wait()
		})
		return stdin
	}

	out := start(ffmpeg)
	if rc.PreviewWindow.Enabled {
		out = io.MultiWriter(out, &preview{start(ffplay)})
	}

	if err := clip.WriteY4M(ctx, out, core.Threads()); err != nil {
		log.Printf("Render failed: %v", err)
		cancel()
	}
	for _, f := range finish {
		f()
	}
}
