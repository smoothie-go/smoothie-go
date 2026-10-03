package graph

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path"
	"path/filepath"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/vsgo"
)

type graph struct {
	core *vsgo.Core
	rc   *recipe.Recipe
}

func (g graph) filter(namespace, function string, args vsgo.Args) *vsgo.Node {
	n, err := g.core.Clip(namespace, function, args)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func Build(core *vsgo.Core, args *cli.Arguments, rc *recipe.Recipe) *vsgo.Node {
	g := graph{core, rc}
	input := filepath.ToSlash(args.InputFile)
	clip := g.filter("bs", "VideoSource", vsgo.Args{
		"source":       input,
		"cachemode":    3, // absolute cache path
		"cachepath":    filepath.ToSlash(args.OutDir) + "/" + path.Base(input) + "-bsindex",
		"showprogress": false,
	})
	original := clip

	if rc.Miscellaneous.DedupThreshold > 0 {
		log.Fatal("Dedup is not available")
	}
	if rc.ColourGrading.Enabled {
		clip = g.tweak(clip)
	}
	if rc.Lut.Enabled {
		log.Fatal("LUT is not available")
	}
	if rc.PreInterp.Enabled {
		clip = g.preInterp(clip)
	}
	if rc.Interpolation.Enabled {
		clip = g.interp(clip)
		if vi := clip.Info(); vi.FPSNum != int64(rc.Interpolation.Fps)*vi.FPSDen {
			log.Fatalf("Interpolation failed. Interpolation FPS: %d/%d, Recipe FPS: %d", vi.FPSNum, vi.FPSDen, rc.Interpolation.Fps)
		}
	}
	if rc.FrameBlending.Enabled {
		clip = g.frameBlend(clip, args.Weighting)
	}
	if rc.ArtifactMasking.Enabled {
		clip = g.mask(clip, original)
	}
	return clip
}

func (g graph) tweak(clip *vsgo.Node) *vsgo.Node {
	cg := g.rc.ColourGrading
	f := clip.Info().Format
	if f.ColorFamily == vsgo.RGB || f.SampleType != vsgo.Integer {
		log.Fatal("Color grading needs an integer YUV or Gray clip")
	}
	shift := f.BitsPerSample - 8
	peak := 1<<f.BitsPerSample - 1
	lumaMin, lumaMax := 0, peak
	chromaMin, chromaMax := 0, peak
	if cg.Coring != 0 {
		// limited range
		lumaMin, lumaMax = 16<<shift, 235<<shift
		chromaMin, chromaMax = 16<<shift, 240<<shift
	}

	if f.ColorFamily == vsgo.YUV {
		hue := cg.Hue * math.Pi / 180
		cos := math.Cos(hue) * cg.Saturation
		sin := math.Sin(hue) * cg.Saturation
		gray := 128 << shift
		u := g.filter("std", "ShufflePlanes", vsgo.Args{"clips": clip, "planes": 1, "colorfamily": vsgo.Gray})
		v := g.filter("std", "ShufflePlanes", vsgo.Args{"clips": clip, "planes": 2, "colorfamily": vsgo.Gray})
		expr := func(format string) *vsgo.Node {
			return g.filter("std", "Expr", vsgo.Args{
				"clips": []*vsgo.Node{u, v},
				"expr":  fmt.Sprintf(format, gray, cos, gray, sin, gray, chromaMin, chromaMax),
			})
		}
		clip = g.filter("std", "ShufflePlanes", vsgo.Args{
			"clips": []*vsgo.Node{
				clip,
				expr("x %d - %g * y %d - %g * + %d + %d max %d min"),
				expr("y %d - %g * x %d - %g * - %d + %d max %d min"),
			},
			"planes":      []int{0, 0, 0},
			"colorfamily": vsgo.YUV,
		})
	}

	lut := make([]int, peak+1)
	for i := range lut {
		luma := float64(i-lumaMin)*cg.Contrast + cg.Brightness + float64(lumaMin)
		lut[i] = min(max(int(luma+0.5), lumaMin), lumaMax)
	}
	return g.filter("std", "Lut", vsgo.Args{"clip": clip, "planes": 0, "lut": lut})
}

// output frame n is input frame n*in/out
func (g graph) changeFPS(clip *vsgo.Node, fpsNum, fpsDen int64) *vsgo.Node {
	vi := clip.Info()
	ratio := big.NewRat(vi.FPSNum*fpsDen, vi.FPSDen*fpsNum)
	in, out := ratio.Num().Int64(), ratio.Denom().Int64()
	if in == 1 {
		// SelectEvery needs a cycle above 1
		in, out = 2, 2*out
	}
	offsets := make([]int64, out)
	for k := range offsets {
		offsets[k] = int64(k) * in / out
	}
	clip = g.filter("std", "SelectEvery", vsgo.Args{"clip": clip, "cycle": in, "offsets": offsets})
	return g.filter("std", "Trim", vsgo.Args{"clip": clip, "length": int64(vi.NumFrames) * out / in})
}

func (g graph) frameBlend(clip *vsgo.Node, weights []float64) *vsgo.Node {
	fb := g.rc.FrameBlending
	format := clip.Info().Format.ID()
	matrix := int64(1) // bt.709
	if fb.BrightBlend {
		frame, err := clip.Frame(0)
		if err != nil {
			log.Fatal(err)
		}
		if m, ok := frame.PropInt("_Matrix"); ok {
			matrix = m
		}
		frame.Free()
		clip = g.filter("resize", "Bicubic", vsgo.Args{
			"clip": clip, "format": vsgo.RGB48,
			"transfer_in_s": "709", "transfer_s": "linear", "matrix_in_s": "709",
		})
	}
	clip = g.filter("frameblender", "FrameBlend", vsgo.Args{"clip": clip, "weights": weights, "log": true})
	if fb.BrightBlend {
		clip = g.filter("resize", "Bicubic", vsgo.Args{
			"clip": clip, "format": format, "matrix": matrix,
			"transfer_in_s": "linear", "transfer_s": "709", "matrix_s": "709",
		})
	}
	return g.changeFPS(clip, int64(fb.Fps), 1)
}

func (g graph) mask(clip, original *vsgo.Node) *vsgo.Node {
	am := g.rc.ArtifactMasking
	maskPath := filepath.Join(am.FolderPath, am.FileName)
	if _, err := os.Stat(maskPath); err != nil {
		log.Fatalf("Mask file not found at: %s", maskPath)
	}
	vi := clip.Info()

	mask := g.filter("bs", "VideoSource", vsgo.Args{"source": maskPath, "cachemode": 0})
	resize := vsgo.Args{
		"clip": mask, "width": vi.Width, "height": vi.Height,
		"format": vsgo.VideoFormat{
			ColorFamily:   vsgo.Gray,
			SampleType:    vi.Format.SampleType,
			BitsPerSample: vi.Format.BitsPerSample,
		}.ID(),
	}
	if mask.Info().Format.ColorFamily == vsgo.RGB {
		resize["matrix_s"] = "709"
	}
	mask = g.filter("resize", "Bicubic", resize)
	if am.Feathering {
		mask = g.filter("std", "Minimum", vsgo.Args{"clip": mask})
		mask = g.filter("std", "BoxBlur", vsgo.Args{"clip": mask, "hradius": 6, "hpasses": 2, "vradius": 6, "vpasses": 2})
	}
	mask = g.filter("std", "Loop", vsgo.Args{"clip": mask, "times": vi.NumFrames})

	if original.Info().NumFrames != vi.NumFrames {
		original = g.changeFPS(original, vi.FPSNum, vi.FPSDen)
		if n := original.Info().NumFrames; n < vi.NumFrames {
			last := g.filter("std", "Trim", vsgo.Args{"clip": original, "first": n - 1})
			last = g.filter("std", "Loop", vsgo.Args{"clip": last, "times": vi.NumFrames - n})
			original = g.filter("std", "Splice", vsgo.Args{"clips": []*vsgo.Node{original, last}})
		} else if n > vi.NumFrames {
			original = g.filter("std", "Trim", vsgo.Args{"clip": original, "length": vi.NumFrames})
		}
	}
	return g.filter("std", "MaskedMerge", vsgo.Args{"clipa": original, "clipb": clip, "mask": mask, "first_plane": true})
}
