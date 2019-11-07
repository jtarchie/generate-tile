package commands

import (
	"fmt"
	"sort"

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

	keys := []string{}
	for field, _ := range validations {
		keys = append(keys, field)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s: %s\n", key, validations[key])
	}

	return nil
}