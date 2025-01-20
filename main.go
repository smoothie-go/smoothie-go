package main

import (
	"fmt"

	"encoding/json"
	"github.com/Hzqkii/smoothie-go/cli"
	"github.com/Hzqkii/smoothie-go/recipe"
)

func main() {
	args := cli.SetupArgs()

	rc := recipe.Parse(args)

	if args.Verbose {
		rcjson, _ := json.MarshalIndent(rc, "", "  ")
		argsjson, _ := json.MarshalIndent(args, "", "  ")

		fmt.Println(string(argsjson))
		fmt.Println(string(rcjson))
	}
}
