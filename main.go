package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/jtarchie/generate-tile/commands"
)

var command struct {
	Generate commands.Generate `command:"generate"`
	Preview  commands.Preview  `command:"preview"`
	Validate commands.Validate `command:"validate"`
}

func main() {
	_, err := flags.Parse(&command)
	if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		os.Exit(1)
	}
}
