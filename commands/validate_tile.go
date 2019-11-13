package commands

import (
	"fmt"
	"io"
	"sort"

	"github.com/jtarchie/tile-builder/metadata"
)

type TileArgs struct {
	Path string `long:"path" description:"path to the pivotal file"`
}

type pivnet struct {
	Token   string `long:"token" description:"the pivnet token from your account"`
	Slug    string `long:"slug" description:"the slug of the product from Pivnet (appears in the URL)"`
	Version string `long:"version" description:"the version of the product to download¬"`
}

type ValidateTile struct {
	Tile   TileArgs `group:"tile" namespace:"tile" env-namespace:"TILE"`
	Pivnet pivnet   `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
	Strict bool     `long:"strict" description:"use strict unmarshaling for the tile"`
	Stdout io.Writer
}

func (p ValidateTile) Execute(_ []string) error {
	payload, err := loadMetadataForTile(p.Tile, p.Pivnet, p.Strict)
	if err != nil {
		return err
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

func loadMetadataForTile(t TileArgs, p pivnet, strict bool) (metadata.Payload, error) {
	if t.Path != "" {
		payload, err := metadata.FromTile(t.Path, false)
		if err != nil {
			return metadata.Payload{}, fmt.Errorf("could not load metadata from tile: %s", err)
		}
		return payload, nil
	} else if p.Token != "" {
		payload, err := metadata.FromPivnet(p.Token, p.Slug, p.Version, false)
		if err != nil {
			return metadata.Payload{}, fmt.Errorf("could not load metadata from pivnet: %s", err)
		}
		return payload, nil
	}

	return metadata.Payload{}, fmt.Errorf("could not determine tile or pivnet metadata")
}
