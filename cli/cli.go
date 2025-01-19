package cli

import (
	"github.com/Hzqkii/smoothie-go/portable"
	"log"
	"os"
)

func SetupArgs() *Arguments {
	argv := os.Args[1:]

	if len(argv) == 0 {
		log.Fatal("You must provide at least one argument")
	}

	if len(argv) == 1 {
		argument := argv[0]

		switch argument {
		case "enc", "encoding", "presets", "encpresets", "macros":
			log.Printf("Encoding presets path: %s", portable.GetEncodingPresetsPath())
			os.Exit(0)

		case "recipe", "rc", "config", "conf", "cfg":
			log.Printf("Recipe path: %s", portable.GetRecipePath())
			os.Exit(0)

		case "models", "rifemod", "mod":
			log.Printf("Models path: %s", portable.GetModelsPath())
			os.Exit(0)
		case "dir", "root", "folder":
			log.Printf("Root directory: %s", portable.GetExecutableDirectory())
		}
	}
	return ParseArgs(argv)
}
