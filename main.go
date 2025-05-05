package main

import (
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

	render.Render(args, rc)

	err := temp.DeleteTempFiles()
	if err != nil {
		log.Fatal(err)
	}
}
