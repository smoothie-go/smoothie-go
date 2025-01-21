package main

import (
	"fmt"

	"encoding/json"

	"github.com/smoothie-go/smoothie-go/cli"
	"github.com/smoothie-go/smoothie-go/recipe"
	"github.com/smoothie-go/smoothie-go/weighting"
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
	//fmt.Println("Ascending:", weighting.Ascending(5))
	//fmt.Println("Descending:", weighting.Descending(5))
	//fmt.Println("Equal:", weighting.Equal(5))
	//fmt.Println("Pyramid:", weighting.Pyramid(5))
	//gaussian, err := weighting.Gaussian(5, 2, 1, [2]float64{0, 2})
	//if err != nil {
	//	fmt.Println("Error:", err)
	//} else {
	//	fmt.Println("Gaussian:", gaussian)
	//}

	fmt.Println(args.InputFps)
	weighting.Parse(args, rc)

	fmt.Println(args.Weighting)
}
