package metadata

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/mholt/archiver"
	"gopkg.in/yaml.v2"
)

type PropertyInput struct {
	Description string
	Label       string
	Reference   string
}

type FormType struct {
	Description    string
	Label          string          `validate:"required"`
	Name           string          `validate:"required"`
	PropertyInputs []PropertyInput `yaml:"property_inputs" validate:"required,dive"`
}

type PropertyBlueprint struct {
	Configurable       bool
	Default            interface{}         `yaml:"default,omitempty"`
	Name               string              `validate:"required"`
	Optional           bool                `yaml:",omitempty"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty" validate:"dive"`
	Type               string              `validate:"required"`
}

type Template struct {
	Consumes string `yaml:",omitempty"`
	Name     string `validate:"required"`
	Provides string `yaml:",omitempty"`
	Release  string `validate:"required"`
}

type ResourceDefinition struct {
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      interface{}
	Label        string `validate:"required"`
	Name         string `validate:"required"`
	Type         string `validate:"required"`
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
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      int
	Label        string
	Name         string
	Type         string
	ZeroIf       ZeroIf `yaml:"zero_if,omitempty"`
}

type JobType struct {
	InstanceDefinition  InstanceDefinition   `yaml:"instance_definition" validate:"required"`
	Manifest            string               `yaml:",omitempty"`
	MaxInFlight         interface{}          `yaml:"max_in_flight" validate:"required"`
	Name                string               `validate:"required"`
	ResourceDefinitions []ResourceDefinition `yaml:"resource_definitions" validate:"required,dive"`
	ResourceLabel       string               `yaml:"resource_label" validate:"required"`
	SingleAZOnly        bool                 `yaml:"single_az_only" validate:"required"`
	Templates           []Template           `validate:"required,dive"`
	UseStemcell         string               `yaml:"use_stemcell,omitempty"`
}

type StemcellCriteria struct {
	EnablePatchSecurityUpdates bool   `yaml:"enable_patch_security_updates"`
	OS                         string `validate:"required"`
	Version                    string `validate:"required"`
}

type Release struct {
	File    string `validate:"required"`
	Name    string `validate:"required"`
	Version string `validate:"required"`
}

type Payload struct {
	Description              string
	FormTypes                []FormType `yaml:"form_types" validate:"dive"`
	IconImage                string     `yaml:"icon_image" validate:"required"`
	JobTypes                 []JobType  `yaml:"job_types" validate:"dive"`
	Label                    string
	MetadataVersion          string              `yaml:"metadata_version"`
	MinimumVersionForUpgrade string              `yaml:"minimum_version_for_upgrade" validate:"required"`
	Name                     string              `validate:"required"`
	OpsmanagerSyslog         bool                `yaml:"opsmanager_syslog"`
	ProductVersion           string              `yaml:"product_version" validate:"required"`
	PropertyBlueprints       []PropertyBlueprint `yaml:"property_blueprints" validate:"dive"`
	Rank                     int
	Releases                 []Release        `validate:"required,dive"`
	StemcellCriteria         StemcellCriteria `yaml:"stemcell_criteria" validate:"required,dive"`
}

var metadataFile = regexp.MustCompile(`metadata\/.*\.yml`)

func FromTile(tilePath string) (Payload, error) {
	var (
		contents bytes.Buffer
		payload  Payload
	)

	archive := archiver.NewZip()
	err := archive.Walk(tilePath, func(f archiver.File) error {
		zfh, ok := f.Header.(zip.FileHeader)
		if ok {
			if metadataFile.MatchString(zfh.Name) {
				_, err := io.Copy(&contents, f)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return payload, fmt.Errorf("could not find metadata file in %s: %s", tilePath, err)
	}

	if contents.Len() == 0 {
		return payload, fmt.Errorf("could not find metadata file in %s", tilePath)
	}

	err = yaml.UnmarshalStrict(contents.Bytes(), &payload)
	if err != nil {
		return payload, fmt.Errorf("could not unmarshal %s: %s", tilePath, err)
	}

	return payload, nil
}
