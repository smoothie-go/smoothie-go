package cli

import (
	"log"
)

func parseArgs(args []string) *Arguments {
	var arguments Arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--verbose", "-v":
			arguments.Verbose = true
			break
		case "--input", "-i":
			if i == len(args)-1 {
				log.Fatal("You must provide an input file")
			}
			arguments.InputFile = args[i+1]
			i++
			break
		case "--output", "-o":
			if i == len(args)-1 {
				log.Fatal("You must provide an output file")
			}
			arguments.OutputFile = args[i+1]
			i++
		case "--outdir", "-od":
			if i == len(args)-1 {
				log.Fatal("You must provide an output directory")
			}
			arguments.OutDir = args[i+1]
			i++
		case "--vpy", "-vp":
			if i == len(args)-1 {
				log.Fatal("You must provide a vpy file")
			}
			arguments.Vpy = args[i+1]
			i++
		case "--encargs", "-e":
			if i == len(args)-1 {
				log.Fatal("You must provide encoding arguments")
			}
			arguments.EncodeArgs = args[i+1]
			i++
		case "--recipe", "-r":
			if i == len(args)-1 {
				log.Fatal("You must provide a recipe name")
			}
			arguments.RecipePath = args[i+1]
			i++
		}
	}
	return &arguments
}
