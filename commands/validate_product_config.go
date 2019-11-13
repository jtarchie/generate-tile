package commands

import (
	"fmt"
	"github.com/jtarchie/tile-builder/configuration"
)

type ValidateProductConfig struct {
	Config string   `long:"config" description:"config file of the product" require:"true"`
	Tile   TileArgs `group:"tile" namespace:"tile" env-namespace:"TILE"`
	Pivnet pivnet   `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
	Strict bool     `long:"strict" description:"use strict unmarshaling for the tile"`
}

func (p ValidateProductConfig) Execute(_ []string) error {
	metadataPayload, err := loadMetadataForTile(p.Tile, p.Pivnet, p.Strict)
	if err != nil {
		return err
	}

	configPayload, err := configuration.FromFile(p.Config)
	if err != nil {
		return err
	}

	for reference, pi := range configPayload.ProductProperties {
		pb, found := metadataPayload.FindPropertyBlueprintFromPropertyInput(reference)

		if !found {
			return fmt.Errorf("cannot determine lookup path of property '%s', expected `.properties` or `.job-name`", reference)
		}

		err = pb.ValidateValue(pi.Value)
		if err != nil {
			return fmt.Errorf("property '%s' value is incorrect: %s", reference, err)
		}
	}

	return nil
}
