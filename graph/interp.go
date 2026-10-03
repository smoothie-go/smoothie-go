package graph

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/smoothie-go/smoothie-go/vsgo"
)

func (g graph) preInterp(clip *vsgo.Node) *vsgo.Node {
	pre := g.rc.PreInterp
	vi := clip.Info()

	matrix, transfer, primaries, chromaloc := "170m", "601", "170m", "left"
	switch {
	case vi.Width >= 3840:
		// ITU-T H.273 (07/2021), note at the bottom of page 20
		matrix, transfer, primaries, chromaloc = "2020ncl", "st2084", "2020", "top_left"
	case vi.Width >= 1280:
		matrix, transfer, primaries = "709", "709", "709"
	case vi.Height == 576:
		matrix, transfer, primaries = "470bg", "470bg", "470bg"
	}

	clip = g.clip("resize", "Bicubic", vsgo.Args{
		"clip": clip, "format": vsgo.RGBS,
		"matrix_in_s": matrix, "transfer_in_s": transfer, "primaries_in_s": primaries,
		"range_in_s": "limited", "chromaloc_in_s": chromaloc,
	})

	factor := 1
	fmt.Sscanf(pre.Factor, "%d", &factor)
	clip = g.clip("rife", "RIFE", vsgo.Args{
		"clip":       clip,
		"factor_num": factor,
		"model_path": strings.Trim(pre.Model, `"`),
		"gpu_id":     g.rc.Interpolation.GpuId,
		"gpu_thread": 1,
		"tta":        pre.Tta,
		"uhd":        pre.Uhd,
		"sc":         pre.SceneChange,
	})

	return g.clip("resize", "Bicubic", vsgo.Args{
		"clip": clip, "format": vi.Format.ID(),
		"matrix_s": matrix, "transfer_s": transfer, "primaries_s": primaries,
		"range_s": "limited", "chromaloc_s": chromaloc,
	})
}

func (g graph) interp(clip *vsgo.Node) *vsgo.Node {
	in := g.rc.Interpolation
	// svp needs 4:2:0, so the luma plane of a 4:4:4 clip is doubled first
	scaled := clip.Info().Format.ID() == vsgo.YUV444P8
	if scaled {
		log.Println("WARNING: Slow interpolation due to YUV444P8 format")
		clip = g.scaleLuma(clip, true)
	}

	switch in.Type {
	case "svp":
		clip = g.interFrame(clip)
	case "of":
		opt := fmt.Sprintf(`{"rate":{"num":%d,"abs":true},"algo":%d,"mask":{"area":0,"area_sharp":1.2},"scene":{"blend":false,"mode":0,"limits":{"blocks":%d}}}`,
			in.Fps, in.Algorithm, in.OfBlocks)
		clip = g.clip("svp2", "SmoothFps_NVOF", vsgo.Args{"clip": clip, "opt": opt, "vec_src": clip})
	}

	if scaled {
		clip = g.scaleLuma(clip, false)
	}
	return clip
}

func (g graph) scaleLuma(clip *vsgo.Node, up bool) *vsgo.Node {
	plane := func(i int) *vsgo.Node {
		return g.clip("std", "ShufflePlanes", vsgo.Args{"clips": clip, "planes": i, "colorfamily": vsgo.Gray})
	}
	y := plane(0)
	vi := y.Info()
	width, height := vi.Width/2, vi.Height/2
	if up {
		width, height = vi.Width*2, vi.Height*2
	}
	y = g.clip("resize", "Point", vsgo.Args{"clip": y, "width": width, "height": height})
	return g.clip("std", "ShufflePlanes", vsgo.Args{
		"clips":       []*vsgo.Node{y, plane(1), plane(2)},
		"planes":      []int{0, 0, 0},
		"colorfamily": vsgo.YUV,
	})
}

// InterFrame from havsfunc
func (g graph) interFrame(clip *vsgo.Node) *vsgo.Node {
	in := g.rc.Interpolation
	preset, tuning := strings.ToLower(in.Speed), strings.ToLower(in.Tuning)
	if !slices.Contains([]string{"medium", "fast", "faster", "fastest"}, preset) {
		log.Fatalf("InterFrame: '%s' is not a valid preset", preset)
	}
	if !slices.Contains([]string{"film", "smooth", "animation", "weak"}, tuning) {
		log.Fatalf("InterFrame: '%s' is not a valid tuning", tuning)
	}
	animation, weak := tuning == "animation", tuning == "weak"

	super := "{gpu:0}"
	if in.Gpu {
		super = "{gpu:1}"
	}
	if preset != "medium" {
		super = "{pel:1," + super[1:]
	}

	block, overlap := 8, 2
	if animation || preset == "fastest" {
		block, overlap = 32, 0
	} else {
		if preset != "medium" || !in.Gpu {
			block = 16
		}
		if preset == "faster" && in.Gpu {
			overlap = 1
		}
	}
	vectors := fmt.Sprintf("{block:{w:%d,overlap:%d},main:{search:{", block, overlap)
	switch {
	case animation:
		vectors += "coarse:{type:2,distance:-6,satd:false},distance:0,"
	case preset == "faster":
		vectors += "coarse:{"
	default:
		vectors += "distance:0,coarse:{"
	}
	switch {
	case animation:
	case weak:
		vectors += "distance:-1,trymany:true,"
	default:
		vectors += "distance:-10,"
	}
	switch {
	case animation || preset == "faster" || preset == "fastest":
		vectors += "bad:{sad:2000}}}}}"
	case weak:
		vectors += "bad:{sad:2000}}}},refine:[{thsad:250,search:{distance:-1,satd:true}}]}"
	default:
		vectors += "bad:{sad:2000}}}},refine:[{thsad:250}]}"
	}

	smooth := fmt.Sprintf("{rate:{num:%d,den:1,abs:true},", in.Fps)
	if in.Gpu {
		smooth += fmt.Sprintf("gpuid:%d,", in.GpuId)
	}
	area := 0
	if tuning == "smooth" {
		area = 150
	}
	smooth += fmt.Sprintf("algo:%d,mask:{cover:80,area:%d,area_sharp:1.2},scene:{", in.Algorithm, area)
	if weak {
		smooth += "blend:true,mode:0,limits:{blocks:50}}}"
	} else {
		smooth += "blend:false,mode:0}}"
	}

	invoke := func(function string, args vsgo.Args) (*vsgo.Node, int64) {
		m, err := g.core.Invoke("svp1", function, args)
		if err != nil {
			log.Fatal(err)
		}
		defer m.Free()
		data, _ := m.Int("data")
		return m.Node("clip"), data
	}
	eightBit := g.clip("fmtc", "bitdepth", vsgo.Args{"clip": clip, "bits": 8})
	superClip, superData := invoke("Super", vsgo.Args{"clip": eightBit, "opt": super})
	vectorsClip, vectorsData := invoke("Analyse", vsgo.Args{"clip": superClip, "sdata": superData, "src": eightBit, "opt": vectors})
	return g.clip("svp2", "SmoothFps", vsgo.Args{
		"clip": clip, "super": superClip, "sdata": superData,
		"vectors": vectorsClip, "vdata": vectorsData, "opt": smooth,
	})
}
