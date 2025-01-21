package recipe

type Recipe struct {
	Interpolation struct {
		Enabled   bool   `ini:"enabled" json:"enabled"`
		Masking   bool   `ini:"masking" json:"masking"`
		Type      string `ini:"type" json:"type"`
		Fps       int    `ini:"fps" json:"fps"`
		Speed     string `ini:"speed" json:"speed"`
		Tuning    string `ini:"tuning" json:"tuning"`
		Algorithm string `ini:"algorithm" json:"algorithm"`
		Gpu       bool   `ini:"use gpu" json:"use_gpu"`
	} `ini:"interpolation" json:"interpolation"`

	FrameBlending struct {
		Enabled     bool    `ini:"enabled" json:"enabled"`
		Fps         int     `ini:"fps" json:"fps"`
		Intensity   float32 `ini:"intensity" json:"intensity"`
		Weighting   string  `ini:"weighting" json:"weighting"`
		BrightBlend bool    `ini:"bright blend" json:"bright_blend"`
	} `ini:"frame blending" json:"frame_blending"`

	FlowBlur struct {
		Enabled    bool   `ini:"enabled" json:"enabled"`
		Masking    bool   `ini:"masking" json:"masking"`
		Amount     int    `ini:"amount" json:"amount"`
		DoBlending string `ini:"do blending" json:"do_blending"`
	} `ini:"flowblur" json:"flow_blur"`

	Output struct {
		Process    string `ini:"process" json:"process"`
		EncArgs    string `ini:"enc args" json:"enc_args"`
		FileFormat string `ini:"file format" json:"file_format"`
		Container  string `ini:"container" json:"container"`
	} `ini:"output" json:"output"`

	PreviewWindow struct {
		Enabled    bool   `ini:"enabled" json:"enabled"`
		Process    string `ini:"process" json:"process"`
		OutputArgs string `ini:"output args" json:"output_args"`
	} `ini:"preview window" json:"preview_window"`

	ArtifactMasking struct {
		Enabled    bool   `ini:"enabled" json:"enabled"`
		Feathering bool   `ini:"feathering" json:"feathering"`
		FolderPath string `ini:"folder path" json:"folder_path"`
		FileName   string `ini:"file name" json:"file_name"`
	} `ini:"artifact masking" json:"artifact_masking"`

	Miscellaneous struct {
		PlayDing           bool   `ini:"play ding" json:"play_ding"`
		AlwaysVerbose      bool   `ini:"always verbose" json:"always_verbose"`
		DedupThreshold     int    `ini:"dedup threshold" json:"dedup_threshold"`
		GlobalOutputFolder string `ini:"global output folder" json:"global_output_folder"`
		SourceIndexing     bool   `ini:"source indexing" json:"source_indexing"`
		FfmpegOptions      string `ini:"ffmpeg options" json:"ffmpeg_options"`
		FfplayOptions      string `ini:"ffplay options" json:"ffplay_options"`
	} `ini:"miscellaneous" json:"miscellaneous"`

	Timescale struct {
		In  float32 `ini:"in" json:"in"`
		Out float32 `ini:"out" json:"out"`
	} `ini:"timescale" json:"timescale"`

	ColourGrading struct {
		Enabled    bool    `ini:"enabled" json:"enabled"`
		Brightness float32 `ini:"brightness" json:"brightness"`
		Saturation float32 `ini:"saturation" json:"saturation"`
		Contrast   float32 `ini:"contrast" json:"contrast"`
	} `ini:"color grading" json:"color_grading"`

	Lut struct {
		Enabled bool    `ini:"enabled" json:"enabled"`
		Path    string  `ini:"path" json:"path"`
		Opacity float32 `ini:"opacity" json:"opacity"`
	} `ini:"lut" json:"lut"`

	PreInterp struct {
		Enabled     bool   `ini:"enabled" json:"enabled"`
		SceneChange bool   `ini:"scene change" json:"scene_change"`
		Tta         bool   `ini:"tta" json:"tta"`
		Uhd         bool   `ini:"uhd" json:"uhd"`
		Masking     bool   `ini:"masking" json:"masking"`
		Factor      string `ini:"factor" json:"factor"`
		Model       string `ini:"model" json:"model"`
	} `ini:"pre-interp" json:"pre_interp"`
}
