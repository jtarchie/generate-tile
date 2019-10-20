package generator

import (
	"fmt"
	"log"
	"regexp"
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

type PropertyBlueprint struct {
	Name               string
	Type               string
	Optional           bool
	Configurable       bool
	Default            interface{}         `yaml:"default,omitempty"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty"`
}

type tilePayload struct {
	Description        string
	FormTypes          []formType          `yaml:"form_types"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints"`
}

func GeneratorTile(specs []SpecPayload) (tilePayload, error) {
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
	for group := range propertiesByGroup {
		groupNames = append(groupNames, group)
	}

	sort.Strings(groupNames)

	for _, group := range groupNames {
		var ft formType
		ft.Name = group
		ft.Label = strings.Title(group)
		ft.Description = fmt.Sprintf("Configuration settings for %s", group)

		propertyNames := []string{}
		for name := range propertiesByGroup[group] {
			propertyNames = append(propertyNames, name)
		}

		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := propertiesByGroup[group][name]

			var propertyInput PropertyInput
			propertyInput.Description = property.Description
			propertyInput.Label = strings.Title(breakApartPropertyName(name))
			propertyBlueprintName := fmt.Sprintf(".properties.%s", name)
			propertyInput.Reference = propertyBlueprintName

			ft.PropertyInputs = append(ft.PropertyInputs, propertyInput)

			var propertyBlueprint PropertyBlueprint
			propertyBlueprint.Name = propertyBlueprintName
			propertyBlueprint.Configurable = true
			propertyBlueprint.Optional = true
			propertyBlueprint.Default = property.Default
			propertyBlueprint.Type = DeterminePropertyBlueprintType(name, property)

			if propertyBlueprint.Type == "collection" {
				propertyBlueprint.PropertyBlueprints = []PropertyBlueprint{
					{
						Name:         "key",
						Type:         "string",
						Optional:     true,
						Configurable: true,
					},
					{
						Name:         "value",
						Type:         "string",
						Optional:     true,
						Configurable: true,
					},
				}
			}

			tile.PropertyBlueprints = append(tile.PropertyBlueprints, propertyBlueprint)
		}

		tile.FormTypes = append(tile.FormTypes, ft)
	}

	return tile, nil
}

func DeterminePropertyBlueprintType(name string, property Property) string {
	if regexp.MustCompile(`[_.]port\z`).MatchString(name) {
		return "port"
	}

	if regexp.MustCompile(`[_.]ip\z`).MatchString(name) {
		return "ip_address"
	}

	switch property.Type {
	case "certificate":
		return "rsa_cert_credentials"
	case "rsa", "ssh":
		return "rsa_pkey_credentials"
	}

	var unknown interface{}
	for _, value := range []interface{}{property.Default, property.Example} {
		if value != nil {
			unknown = value
			break
		}
	}

	switch unknown.(type) {
	case int, float32, float64:
		return "integer"
	case nil, string:
		return "string"
	case bool:
		return "boolean"
	case []interface{}:
		return "string_list"
	case map[interface{}]interface{}:
		return "collection"
	}

	log.Panicf("not able to determine type for property %s: %#v", name, property)
	return "string"
}

func breakApartPropertyName(name string) string {
	dots := strings.Split(name, ".")

	parts := []string{}
	for _, s := range dots {
		parts = append(parts, strings.Split(s, "_")...)
	}

	return strings.Join(parts, " ")
}
