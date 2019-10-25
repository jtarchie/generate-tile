package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/jtarchie/generate-tile/commands"
	"log"
	"os"
)

var command struct {
	Generate commands.Generate `command:"generate"`
}

func main() {
	_, err := flags.Parse(&command)
	if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
		os.Exit(0)
	} else {
		log.Fatalf("could not run executable: %s", err)
	}
}

