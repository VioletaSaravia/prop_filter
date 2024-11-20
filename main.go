package main

import (
	"log"
	"os"
	"prop_filter_cli/lib"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Args:     true,
		Flags:    lib.Flags,
		HelpName: "Property Filter",
		Name:     "propfilter",
		Usage:    "Filter large sets of real estate properties based on their particular attributes.",
		Action:   lib.Run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
