package generator

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
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
	Optional           bool `yaml:",omitempty"`
	Configurable       bool
	Default            interface{}         `yaml:"default,omitempty"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty"`
}

type Template struct {
	Name     string
	Release  string
	Consumes string `yaml:",omitempty"`
	Provides string `yaml:",omitempty"`
}

type ResourceDefinition struct {
	Name         string
	Configurable bool
	Default      interface{}
	Constraints  Constraints `yaml:",omitempty"`
	Label        string
	Type         string
}

type Constraints struct {
	Max int `yaml:",omitempty"`
	Min int `yaml:",omitempty"`
}

type ZeroIf struct {
	PropertyReference string   `yaml:"property_reference"`
	PropertyValues    []string `yaml:"property_values"`
}

type InstanceDefinition struct {
	Name         string
	Label        string
	Configurable bool
	Default      int
	Constraints  Constraints `yaml:",omitempty"`
	Type         string
	ZeroIf       ZeroIf `yaml:"zero_if,omitempty"`
}

type JobType struct {
	Name                string
	ResourceLabel       string `yaml:"resource_label"`
	Templates           []Template
	SingleAZOnly        bool                 `yaml:"single_az_only"`
	MaxInFlight         interface{}          `yaml:"max_in_flight"`
	UseStemcell         string               `yaml:"use_stemcell,omitempty"`
	ResourceDefinitions []ResourceDefinition `yaml:"resource_definitions"`
	InstanceDefinition  InstanceDefinition   `yaml:"instance_definition"`
	Manifest            string               `yaml:",omitempty"`
}

type tilePayload struct {
	Description        string
	FormTypes          []formType          `yaml:"form_types"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints"`
	JobTypes           []JobType           `yaml:"job_types"`
}

func GeneratorTile(release BoshReleasePayload) (tilePayload, error) {
	specs := release.Specs

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
		ft.Label = strings.Title(breakApartName(group))
		ft.Description = fmt.Sprintf("Configuration settings for %s", group)

		propertyNames := []string{}
		for name := range propertiesByGroup[group] {
			propertyNames = append(propertyNames, name)
		}

		sort.Strings(propertyNames)

		for _, name := range propertyNames {
			property := propertiesByGroup[group][name]

			createPropertyInput(property, name, &ft)

			propertyBlueprint, err := createPropertyBlueprint(property, name)
			if err != nil {
				return tilePayload{}, err
			}

			tile.PropertyBlueprints = append(tile.PropertyBlueprints, propertyBlueprint)
		}

		tile.FormTypes = append(tile.FormTypes, ft)
	}

	for _, spec := range specs {
		jobType, err := createJob(spec, release)
		if err != nil {
			return tilePayload{}, err
		}
		tile.JobTypes = append(tile.JobTypes, jobType)
	}

	return tile, nil
}

func createJob(spec SpecPayload, release BoshReleasePayload) (JobType, error) {
	var jobType JobType

	jobType.Name = spec.Name
	jobType.ResourceLabel = strings.Title(breakApartName(spec.Name))
	templates, err := generateTemplateForSpec(release, spec)
	if err != nil {
		return JobType{}, err
	}

	jobType.Templates = templates
	jobType.MaxInFlight = 1

	attachInstanceDefinition(&jobType)
	attachResourceDefinitions(&jobType)

	manifest := map[string]interface{}{}
	for name, property := range spec.Properties {
		parts := strings.Split(name, ".")

		root := manifest
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			if _, ok := root[part]; !ok {
				root[part] = map[string]interface{}{}
			}
			root = root[part].(map[string]interface{})
		}
		option, err := CreateManifestFromProperty(name, property)
		if err != nil {
			return JobType{}, fmt.Errorf("could not create manifest for property %s: %s", name, err)
		}
		root[parts[len(parts)-1]] = option
	}

	manifestYAML, err := yaml.Marshal(manifest)
	if err != nil {
		return JobType{}, err
	}

	jobType.Manifest = string(manifestYAML)
	return jobType, nil
}

func createPropertyInput(property Property, name string, ft *formType) {
	var propertyInput PropertyInput
	propertyInput.Description = property.Description
	propertyInput.Label = strings.Title(breakApartName(name))
	propertyInput.Reference = fmt.Sprintf(".properties.%s", propertyBlueprintNameFromPropertyName(name))
	ft.PropertyInputs = append(ft.PropertyInputs, propertyInput)
}

func createPropertyBlueprint(property Property, name string) (PropertyBlueprint, error) {
	var propertyBlueprint PropertyBlueprint
	propertyBlueprint.Name = propertyBlueprintNameFromPropertyName(name)
	propertyBlueprint.Configurable = true

	def, err := DeterminePropertyBlueprintDefault(name, property)
	if err != nil {
		return PropertyBlueprint{}, err
	}
	propertyBlueprint.Default = def
	if propertyBlueprint.Default == nil {
		propertyBlueprint.Optional = true
	}

	pbType, err := DeterminePropertyBlueprintType(name, property)
	if err != nil {
		return PropertyBlueprint{}, err
	}

	propertyBlueprint.Type = pbType

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

	return propertyBlueprint, nil
}

func DeterminePropertyBlueprintDefault(name string, property Property) (interface{}, error) {
	switch v := property.Default.(type) {
	case int, float32, float64:
		return v, nil
	case nil, string:
		return v, nil
	case bool:
		return v, nil
	case []interface{}:
		list := []string{}
		for _, item := range v {
			list = append(list, fmt.Sprintf("%s", item))
		}

		return strings.Join(list, ","), nil
	case map[interface{}]interface{}:
		return nil, nil
	}
	return nil, fmt.Errorf("could not determine default for %s of %t", name, property.Default)
}

func propertyBlueprintNameFromPropertyName(name string) string {
	return strings.Join(strings.Split(name, "."), "__")
}

func CreateManifestFromProperty(name string, property Property) (interface{}, error) {
	pbType, err := DeterminePropertyBlueprintType(name, property)
	if err != nil {
		return nil, err
	}

	propertyBlueprintName := propertyBlueprintNameFromPropertyName(name)

	switch pbType {
	case "rsa_cert_credentials":
		return map[string]string{
			"certificate": fmt.Sprintf("((.properties.%s.certificate))", propertyBlueprintName),
			"private_key": fmt.Sprintf("((.properties.%s.private_key))", propertyBlueprintName),
		}, nil
	}
	return fmt.Sprintf("((.properties.%s.value))", propertyBlueprintName), nil
}

func attachResourceDefinitions(jobType *JobType) {
	jobType.ResourceDefinitions = []ResourceDefinition{
		{
			Name:         "cpu",
			Configurable: true,
			Default:      1,
			Constraints: Constraints{
				Min: 1,
			},
			Label: "CPU",
			Type:  "integer",
		},
		{
			Name:         "ram",
			Configurable: true,
			Default:      8192,
			Constraints: Constraints{
				Min: 8192,
			},
			Label: "RAM",
			Type:  "integer",
		},
		{
			Name:         "ephemeral_disk",
			Configurable: true,
			Default:      10240,
			Constraints: Constraints{
				Min: 10240,
			},
			Label: "Ephemeral Disk",
			Type:  "integer",
		},
		{
			Name:         "persistent_disk",
			Configurable: true,
			Default:      10240,
			Constraints: Constraints{
				Min: 10240,
			},
			Label: "Persistent Disk",
			Type:  "integer",
		},
	}
}

func attachInstanceDefinition(jobType *JobType) {
	jobType.InstanceDefinition = InstanceDefinition{
		Name:         "instances",
		Label:        "Instances",
		Configurable: true,
		Default:      1,
		Constraints: Constraints{
			Min: 1,
		},
		Type: "integer",
	}
}

type consumer struct {
	From string
}

type provider struct {
	As string
}

func generateTemplateForSpec(release BoshReleasePayload, spec SpecPayload) ([]Template, error) {
	var template Template

	template.Name = spec.Name
	template.Release = release.Name

	consuming := map[string]consumer{}
	for _, consume := range spec.Consumes {
		generateConsumer(release, spec, consume, consuming)
	}

	contents, err := yaml.Marshal(consuming)
	if err != nil {
		return nil, fmt.Errorf("could not marshal consumers: %s", err)
	}
	template.Consumes = string(contents)

	providing := map[string]provider{}
	for _, provide := range spec.Provides {
		providing[provide.Name] = provider{
			As: fmt.Sprintf("%s-%s", spec.Name, provide.Name),
		}
	}

	contents, err = yaml.Marshal(providing)
	if err != nil {
		return nil, fmt.Errorf("could not marshal providers: %s", err)
	}
	template.Provides = string(contents)

	return []Template{template}, nil
}

func generateConsumer(release BoshReleasePayload, payload SpecPayload, consume consumePayload, consuming map[string]consumer) {
	for _, spec := range release.Specs {
		if spec.Name == payload.Name {
			continue
		}

		for _, provide := range spec.Provides {
			if provide.Name == consume.Name && provide.Type == consume.Type {
				consuming[consume.Name] = consumer{
					From: fmt.Sprintf("%s-%s", spec.Name, consume.Name),
				}
				return
			}
		}
	}
}

func DeterminePropertyBlueprintType(name string, property Property) (string, error) {
	if regexp.MustCompile(`[_.]port\z`).MatchString(name) {
		return "port", nil
	}

	if regexp.MustCompile(`[_.]ip\z`).MatchString(name) {
		return "ip_address", nil
	}

	switch property.Type {
	case "certificate":
		return "rsa_cert_credentials", nil
	case "rsa", "ssh":
		return "rsa_pkey_credentials", nil
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
		return "integer", nil
	case nil, string:
		return "string", nil
	case bool:
		return "boolean", nil
	case []interface{}:
		return "string_list", nil
	case map[interface{}]interface{}:
		return "collection", nil
	}

	return "", fmt.Errorf("not able to determine type for property %s: %#v", name, property)
}

func breakApartName(name string) string {
	dots := strings.Split(name, ".")

	parts := []string{}
	for _, s := range dots {
		parts = append(parts, strings.Split(s, "_")...)
	}

	return strings.Join(parts, " ")
}
