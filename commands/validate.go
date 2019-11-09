package commands

import (
	"fmt"
	"github.com/jtarchie/tile-builder/metadata"
	"io"
	"sort"
)

type Validate struct {
	Path   string `long:"path" description:"path to the pivotal file"`
	Pivnet struct {
		Token   string `long:"token" description:"the pivnet token from your account"`
		Slug    string `long:"slug" description:"the slug of the product from Pivnet (appears in the URL)"`
		Version string `long:"version" description:"the version of the product to download¬"`
	} `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
	Stdout io.Writer
}

func (p Validate) Execute(_ []string) error {
	var (
		payload metadata.Payload
		err     error
	)

	if p.Path != "" {
		payload, err = metadata.FromTile(p.Path)
		if err != nil {
			return fmt.Errorf("could not load metadata from tile: %s", err)
		}
	} else if p.Pivnet.Token != "" {
		payload, err = metadata.FromPivnet(p.Pivnet.Token, p.Pivnet.Slug, p.Pivnet.Version)
		if err != nil {
			return fmt.Errorf("could not load metadata from pivnet: %s", err)
		}
	}

	validations, err := payload.Validate()
	if err != nil {
		return fmt.Errorf("could not determine validations on tile: %s", err)
	}

	keys := []string{}
	for field := range validations {
		keys = append(keys, field)
	}

	sort.Strings(keys)

	for _, key := range keys {
		_, _ = fmt.Fprintf(p.Stdout, "%s: %s\n", key, validations[key])
	}

	return nil
}
