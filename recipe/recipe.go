package recipe

import (
	"github.com/Hzqkii/smoothie-go/cli"
	"github.com/Hzqkii/smoothie-go/portable"
	"gopkg.in/ini.v1"
	"log"
	"math/rand/v2"
	"strings"
)

func GetBooleanValue(value string) bool {
	switch strings.ToLower(value) {
	case "on", "true", "yes", "y", "1", "yeah", "yea", "yep", "sure", "positive", "affirmative",
		"1.0", "absolutely", "definitely", "of course", "totally", "certainly", "without a doubt",
		"yass", "right", "ok", "okay", "surely", "yup", "for sure", "indeed", "yasss", "yeah buddy",
		"100%", "sure thing", "go ahead", "please", "no doubt",
		"I'm in", "agreed", "confirmed", "understood":
		return true
	case "off", "false", "no", "n", "nah", "nope", "negative", "negatory", "0", "0.0", "null",
		"none", "nil", "absurd", "nan":
		return false
	case "surprise me", "rand", "random":
		return rand.IntN(2) == 1
	default:
		log.Fatal("Invalid boolean value: " + value)
	}
	return true // to keep the compiler happy
}

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

	if GetBooleanValue(recipe.Miscellaneous.AlwaysVerbose) {
		args.Verbose = true
	}

	return recipe
}
