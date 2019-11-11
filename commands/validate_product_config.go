package commands

import (
	"fmt"
	"github.com/jtarchie/tile-builder/configuration"
	"strings"
)

type ValidateProductConfig struct {
	Config string   `long:"config" description:"config file of the product" require:"true"`
	Tile   TileArgs `group:"tile" namespace:"tile" env-namespace:"TILE"`
	Pivnet pivnet   `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
}

func (p ValidateProductConfig) Execute(_ []string) error {
	metadataPayload, err := loadMetadataForTile(p.Tile, p.Pivnet)
	if err != nil {
		return err
	}

	configPayload, err := configuration.FromFile(p.Config)
	if err != nil {
		return err
	}

	for name, _ := range configPayload.ProductProperties {
		parts := strings.Split(name, ".")
		if parts[0] == "properties" {
			for _, propertyBlueprint := range metadataPayload.PropertyBlueprints {
				if parts[1] == propertyBlueprint.Name {
					continue
				}
			}
			continue
		}

		for _, job := range metadataPayload.JobTypes {
			if parts[0] == job.Name {
				continue
			}
		}

		return fmt.Errorf("cannot determine lookup path of property '%s', expected `.properties` or `.job-name`", name)
	}

	return nil
}
