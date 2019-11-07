package commands

import (
	"fmt"
	"github.com/jtarchie/generate-tile/metadata"
)

type Validate struct {
	Path string `long:"path" required:"true" description:"path to the pivotal file"`
}

func (p Validate) Execute(_ []string) error {
	payload, err := metadata.FromTile(p.Path)
	if err != nil {
		return fmt.Errorf("could not load metadata from tile: %s", err)
	}

	validations, err := payload.Validate()
	if err != nil {
		return fmt.Errorf("could not determine validations on tile: %s", err)
	}

	for field, msg := range validations {
		fmt.Printf("%s: %s\n", field, msg)
	}

	return nil
}