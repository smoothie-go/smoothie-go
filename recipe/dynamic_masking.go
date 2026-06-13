package recipe

import (
	"log"
	"regexp"

	"github.com/smoothie-go/smoothie-go/cli"
)

func ParseDynamicMasking(recipe *Recipe, args *cli.Arguments) error {
	dynMasking, err := recipe.ini_config.GetSection("dynamic masking")
	if err != nil {
		return err
	}
	if !dynMasking.Key("enabled").MustBool(true) {
		return nil
	}
	matched := false
	var matchedMask string
	for _, key := range dynMasking.KeyStrings() {
		if key == "enabled" {
			continue
		}
		regex, err := regexp.Compile(".*" + key + ".*")
		if err != nil {
			return err
		}
		if regex.Match([]byte(args.InputFile)) {
			matchedMask = dynMasking.Key(key).String()
			recipe.ArtifactMasking.FileName = matchedMask
			matched = true
		}
	}
	if matched && recipe.ArtifactMasking.Enabled {
		log.Printf("Using dynamic mask: %s\n", matchedMask)
	}
	return nil
}
