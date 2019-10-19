package generator

import (
	"fmt"
	"sort"
	"strings"
)

type PropertyInput struct {
	Reference   string
	Label       string
	Description string
}

type formType struct {
	Name           string
	Label          string
	Description    string
	PropertyInputs []PropertyInput `yaml:"property_inputs"`
}

type tilePayload struct {
	Description string
	FormTypes   []formType `yaml:"form_types"`
}

func GeneratorTile(specs []specPayload) (tilePayload, error) {
	var tile tilePayload

	propertiesByGroup := map[string]map[string]Property{}

	for _, payload := range specs {
		for name, property := range payload.Properties {
			parts := strings.Split(name, ".")

			group := "properties"
			if len(parts) > 1 {
				group = parts[0]
			}

			if propertiesByGroup[group] == nil {
				propertiesByGroup[group] = map[string]Property{}
			}

			propertiesByGroup[group][name] = property
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

		propertyNames := []string{}
		for name, _ := range propertiesByGroup[group] {
			propertyNames = append(propertyNames, name)
		}

		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			var propertyInput PropertyInput

			property := propertiesByGroup[group][name]

			propertyInput.Description = property.Description
			propertyInput.Label = strings.Title(breakApartPropertyName(name))
			propertyInput.Reference = fmt.Sprintf(".properties.%s", name)

			ft.PropertyInputs = append(ft.PropertyInputs, propertyInput)
		}

		tile.FormTypes = append(tile.FormTypes, ft)
	}

	return tile, nil
}

func breakApartPropertyName(name string) string {
	dots := strings.Split(name, ".")

	parts := []string{}
	for _, s := range dots {
		parts = append(parts, strings.Split(s, "_")...)
	}

	return strings.Join(parts, " ")
}