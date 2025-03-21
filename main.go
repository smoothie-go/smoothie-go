package main

import (
	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/render"
	"github.com/smoothie-go/smoothie-go/server"
	"github.com/smoothie-go/smoothie-go/weighting"
	"os"
)

func main() {
	if os.Getenv("SM_SERVER") == "1" {
		server.SetupRouter().Run()
	}

	args := cli.SetupArgs()

	rc := recipe.Parse(args)

	weighting.Parse(args, rc)

	render.Render(args, rc)
}
