package cli

type Arguments struct {
	InputFile  string
	OutputFile string
	OutDir     string
	Vpy        string
	// overrides; later
	EncodeArgs string
	RecipePath string
	Verbose    bool
}
