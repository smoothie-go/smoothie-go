package recipe

import (
	"log"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/portable"
	"gopkg.in/ini.v1"
)

func Parse(args *cli.Arguments) *Recipe {
	defaults, err := ini.Load(portable.GetDefaultRecipePath())
	if err != nil {
		log.Fatal(err)
	}

	recipe := &Recipe{}

	err = defaults.MapTo(recipe)
	if err != nil {
		log.Fatal(err)
	}

	config, err := ini.Load(args.RecipePath)
	if err != nil {
		log.Fatal(err)
	}

	err = config.MapTo(recipe)
	if err != nil {
		log.Fatal(err)
	}

	if recipe.Miscellaneous.AlwaysVerbose {
		args.Verbose = true
	}

	recipe.ini_config = config

	ParseDynamicMasking(recipe, args)

	recipe = Validate(args, recipe)

	return recipe
}
