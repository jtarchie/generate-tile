package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/jtarchie/tile-builder/commands"
)

var command struct {
	Generate     commands.Generate     `command:"generate"`
	Preview      commands.Preview      `command:"preview"`
	ValidateTile commands.ValidateTile `command:"validate-tile"`
}

func main() {
	command.ValidateTile = commands.ValidateTile{
		Stdout: os.Stdout,
	}

	_, err := flags.Parse(&command)
	if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		os.Exit(1)
	}
}
