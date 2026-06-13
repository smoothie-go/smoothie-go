package recipe

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
)

func Validate(args *cli.Arguments, recipe *Recipe) *Recipe {
	// enc args
	outEncArgs, err := ParseEncodingArgs(portable.GetEncodingPresetsPath(), recipe.Output.EncArgs)
	if err != nil {
		log.Fatal(err)
	}

	recipe.Output.EncArgs = outEncArgs

	// model
	if recipe.PreInterp.Model == "" ||
		recipe.PreInterp.Model == "default" ||
		recipe.PreInterp.Model == "auto" {
		if recipe.PreInterp.Tta {
			recipe.PreInterp.Model = portable.GetDefaultTtaModelPath()
		} else {
			recipe.PreInterp.Model = portable.GetDefaultModelPath()
		}
	}

	if recipe.Timescale.In <= 0 {
		recipe.Timescale.In = 1.0
	}
	if recipe.Timescale.Out <= 0 {
		recipe.Timescale.Out = 1.0
	}

	if recipe.FrameBlending.Enabled && recipe.FrameBlending.Fps <= 0 {
		log.Fatal("Frame blending FPS must be greater than zero")
	}
	if recipe.Interpolation.Enabled && recipe.Interpolation.Fps <= 0 {
		log.Fatal("Interpolation FPS must be greater than zero")
	}

	if recipe.Interpolation.Enabled && recipe.Interpolation.Gpu {
		if !isOpenCLAvailable() {
			log.Println("WARNING: GPU interpolation is enabled, but OpenCL was not detected on your system. GPU acceleration may fail or fallback to CPU.")
		}
	}

	interpEnabled := recipe.Interpolation.Enabled
	frameBlendingEnabled := recipe.FrameBlending.Enabled
	preinterpEnabled := recipe.PreInterp.Enabled

	var preinterpFactor int
	n, err := fmt.Sscanf(recipe.PreInterp.Factor, "%d", &preinterpFactor)
	if err != nil || n < 1 {
		preinterpFactor = 1
	}
	interpFps := recipe.Interpolation.Fps
	frameBlendingFps := recipe.FrameBlending.Fps
	inputFps := args.InputFps

	if interpEnabled && frameBlendingEnabled {
		if interpFps <= inputFps {
			log.Fatal("Interpolation FPS must be higher than input FPS")
		}
		if frameBlendingFps >= interpFps {
			log.Fatal("Frame blending FPS must be lower than interpolation FPS")
		}
	} else if preinterpEnabled && frameBlendingEnabled {
		if (inputFps * preinterpFactor) <= inputFps {
			log.Fatal("Interpolation FPS must be higher than input FPS")
		}
	} else if frameBlendingEnabled && frameBlendingFps > inputFps {
		log.Fatal("Frame blending FPS cannot be higher than input FPS")
	}

	return recipe
}

func isOpenCLAvailable() bool {
	if runtime.GOOS == "windows" {
		handle, err := syscall.LoadLibrary("OpenCL.dll")
		if err == nil {
			syscall.FreeLibrary(handle)
			return true
		}
		return false
	}
	if runtime.GOOS == "linux" {
		if files, err := os.ReadDir("/etc/OpenCL/vendors"); err == nil && len(files) > 0 {
			return true
		}
		commonPaths := []string{
			"/usr/lib/libOpenCL.so",
			"/usr/lib/libOpenCL.so.1",
			"/usr/lib64/libOpenCL.so",
			"/usr/lib64/libOpenCL.so.1",
			"/usr/lib/x86_64-linux-gnu/libOpenCL.so",
			"/usr/lib/x86_64-linux-gnu/libOpenCL.so.1",
		}
		for _, path := range commonPaths {
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
		return false
	}
	return true
}
