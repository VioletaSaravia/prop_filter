package main

import (
	"log"
	"os"
	"prop_filter_cli/lib"
)

func main() {
	if err := lib.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
