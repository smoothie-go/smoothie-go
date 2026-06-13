package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/render"
	"github.com/smoothie-go/smoothie-go/temp"
	"github.com/smoothie-go/smoothie-go/weighting"
)

func main() {
	args := cli.SetupArgs()

	rc := recipe.Parse(args)

	weighting.Parse(args, rc)

	if args.DevDump {
		jargs, _ := json.MarshalIndent(args, "", "\t")
		jrc, _ := json.MarshalIndent(rc, "", "\t")

		fmt.Printf(`
=============== args ===============
%s


============== recipe ==============
%s

		`, string(jargs), string(jrc))

		return
	}

	render.Render(args, rc)

	err := temp.DeleteTempFiles()
	if err != nil {
		log.Fatal(err)
	}
}
