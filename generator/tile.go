package generator

import (
	"fmt"
	"sort"
	"strings"
)

type propertyInput struct {
}

type formType struct {
	Name           string
	Label          string
	Description    string
	PropertyInputs []propertyInput `yaml:"property_inputs"`
}

type tilePayload struct {
	Description string
	FormTypes   []formType `yaml:"form_types"`
}

func GeneratorTile(specs []specPayload) (tilePayload, error) {
	var tile tilePayload

	propertiesByGroup := map[string][]Property{}

	for _, payload := range specs{
		for name, property := range payload.Properties {
			parts := strings.Split(name, ".")

			group := "properties"
			if len(parts) > 1 {
				group = parts[0]
			}

			propertiesByGroup[group] = append(propertiesByGroup[group], property)
		}
	}

	groupNames := []string{}
	for group, _ := range propertiesByGroup {
		groupNames = append(groupNames, group)
	}

	sort.Strings(groupNames)

	for _, group := range groupNames {
		var ft formType
		ft.Name = group
		ft.Label = strings.Title(group)
		ft.Description = fmt.Sprintf("Configuration settings for %s", group)

		tile.FormTypes = append(tile.FormTypes, ft)
	}

	return tile, nil
}
