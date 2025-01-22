package cli

type Arguments struct {
	InputFile  string `json:"input_file"`
	InputFps   int    `json:"input_fps"`
	OutputFile string `json:"output_file"`
	OutDir     string `json:"out_dir"`
	// overrides; later
	EncodeArgs string    `json:"encode_args"`
	RecipePath string    `json:"recipe_path"`
	LogFile    string    `json:"log_file"`
	Verbose    bool      `json:"verbose"`
	Weighting  []float64 `json:"weighting"`
}
