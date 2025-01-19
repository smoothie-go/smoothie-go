package main

import (
	"fmt"
	"github.com/Hzqkii/smoothie-go/cli"
)

func main() {
	args := cli.SetupArgs()

	fmt.Println(args)
}
