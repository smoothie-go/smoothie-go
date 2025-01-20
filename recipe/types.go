package recipe

type Recipe struct {
	Interpolation struct {
		Enabled   string `ini:"enabled"`
		Masking   string `ini:"masking"`
		Fps       int    `ini:"fps"`
		Speed     string `ini:"speed"`
		Tuning    string `ini:"tuning"`
		Algorithm string `ini:"algorithm"`
		Gpu       string `ini:"use gpu"`
	} `ini:"interpolation"`

	FrameBlending struct {
		Enabled     string  `ini:"enabled"`
		Fps         int     `ini:"fps"`
		Intensity   float32 `ini:"intensity"`
		Weighting   string  `ini:"weighting"`
		BrightBlend string  `ini:"bright blend"`
	} `ini:"frame blending"`

	FlowBlur struct {
		Enabled    string `ini:"enabled"`
		Masking    string `ini:"masking"`
		Amount     int    `ini:"amount"`
		DoBlending string `ini:"do blending"`
	} `ini:"flowblur"`

	Output struct {
		Process    string `ini:"process"`
		EncArgs    string `ini:"enc args"`
		FileFormat string `ini:"file format"`
		Container  string `ini:"container"`
	} `ini:"output"`

	PreviewWindow struct {
		Enabled    string `ini:"enabled"`
		Process    string `ini:"process"`
		OutputArgs string `ini:"output args"`
	} `ini:"preview window"`

	ArtifactMasking struct {
		Enabled    string `ini:"enabled"`
		Feathering string `ini:"feathering"`
		FolderPath string `ini:"folder path"`
		FileName   string `ini:"file name"`
	} `ini:"artifact masking"`

	Miscellaneous struct {
		PlayDing           string `ini:"play ding"`
		AlwaysVerbose      string `ini:"always verbose"`
		DedupThreshold     int    `ini:"dedup threshold"`
		GlobalOutputFolder string `ini:"global output folder"`
		SourceIndexing     string `ini:"source indexing"`
		FfmpegOptions      string `ini:"ffmpeg options"`
		FfplayOptions      string `ini:"ffplay options"`
	} `ini:"miscellaneous"`

	Timescale struct {
		In  float32 `ini:"in"`
		Out float32 `ini:"out"`
	} `ini:"timescale"`

	ColourGrading struct {
		Enabled    string  `ini:"enabled"`
		Brightness float32 `ini:"brightness"`
		Saturation float32 `ini:"saturation"`
		Contrast   float32 `ini:"contrast"`
	} `ini:"color grading"`

	Lut struct {
		Enabled string  `ini:"enabled"`
		Path    string  `ini:"path"`
		Opacity float32 `ini:"opacity"`
	} `ini:"lut"`

	PreInterp struct {
		Enabled     string `ini:"enabled"`
		SceneChange string `ini:"scene change"`
		Tta         string `ini:"tta"`
		Uhd         string `ini:"uhd"`
		Masking     string `ini:"masking"`
		Factor      string `ini:"factor"`
		Model       string `ini:"model"` // can be auto
	} `ini:"pre-interp"` // wtf dawg, why does this use `-` when the others use ` `, gotta keep for backwards compatibility
}
