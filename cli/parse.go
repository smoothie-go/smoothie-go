package cli

import (
	"log"
	"os"

	"github.com/smoothie-go/smoothie-go/portable"
)

func parseArgs(args []string) *Arguments {
	var arguments Arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--verbose", "-v":
			arguments.Verbose = true
		case "--code-execution-pipeline", "-cep":
			arguments.CEP = true
		case "--input", "-i":
			if i == len(args)-1 {
				log.Fatal("You must provide an input file")
			}
			arguments.InputFile = args[i+1]
			i++
		case "--output", "-o":
			if i == len(args)-1 {
				log.Fatal("You must provide an output file")
			}
			arguments.OutputFile = args[i+1]
			arguments.UserSpecifiedOutput = true
			i++
		case "--outdir", "-od":
			if i == len(args)-1 {
				log.Fatal("You must provide an output directory")
			}
			arguments.OutDir = args[i+1]
			arguments.UserSpecifiedOutDir = true
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
			arguments.RecipePath = args[i+1] + ".ini"
			i++
		case "--dump-scripts", "-ds":
			if i == len(args)-1 {
				log.Fatal("You must provide a directory")
			}
			err := portable.DropScriptsAtPath(args[i+1])
			if err != nil {
				log.Fatal("Unable to drop scripts at " + args[i+1] + "\nError: " + err.Error())
			}
			log.Println("Ok!")
			os.Exit(0)
		case "-__dev_dump":
			arguments.DevDump = true
		}
	}
	return &arguments
}
