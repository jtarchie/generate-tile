package metadata

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/mholt/archiver"
	"gopkg.in/yaml.v2"
	"io"
	"regexp"
)

type PropertyInput struct {
	Description string
	Label       string
	Reference   string
}

type FormType struct {
	Description    string
	Label          string
	Name           string
	PropertyInputs []PropertyInput `yaml:"property_inputs"`
}

type PropertyBlueprint struct {
	Configurable       bool
	Default            interface{} `yaml:"default,omitempty"`
	Name               string
	Optional           bool                `yaml:",omitempty"`
	PropertyBlueprints []PropertyBlueprint `yaml:"property_blueprints,omitempty"`
	Type               string
}

type Template struct {
	Consumes string `yaml:",omitempty"`
	Name     string
	Provides string `yaml:",omitempty"`
	Release  string
}

type ResourceDefinition struct {
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      interface{}
	Label        string
	Name         string
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
	Configurable bool
	Constraints  Constraints `yaml:",omitempty"`
	Default      int
	Label        string
	Name         string
	Type         string
	ZeroIf       ZeroIf `yaml:"zero_if,omitempty"`
}

type JobType struct {
	InstanceDefinition  InstanceDefinition `yaml:"instance_definition"`
	Manifest            string             `yaml:",omitempty"`
	MaxInFlight         interface{}        `yaml:"max_in_flight"`
	Name                string
	ResourceDefinitions []ResourceDefinition `yaml:"resource_definitions"`
	ResourceLabel       string               `yaml:"resource_label"`
	SingleAZOnly        bool                 `yaml:"single_az_only"`
	Templates           []Template
	UseStemcell         string `yaml:"use_stemcell,omitempty"`
}

type StemcellCriteria struct {
	EnablePatchSecurityUpdates bool `yaml:"enable_patch_security_updates"`
	OS                         string
	Version                    string
}

type Release struct {
	File    string
	Name    string
	Version string
}

type Payload struct {
	Description              string
	FormTypes                []FormType `yaml:"form_types" validate:"dive"`
	IconImage                string     `yaml:"icon_image"`
	JobTypes                 []JobType  `yaml:"job_types"`
	Label                    string
	MetadataVersion          string              `yaml:"metadata_version"`
	MinimumVersionForUpgrade string              `yaml:"minimum_version_for_upgrade"`
	Name                     string              `validate:"required"`
	OpsmanagerSyslog         bool                `yaml:"opsmanager_syslog"`
	ProductVersion           string              `yaml:"product_version"`
	PropertyBlueprints       []PropertyBlueprint `yaml:"property_blueprints"`
	Rank                     int
	Releases                 []Release
	StemcellCriteria         StemcellCriteria `yaml:"stemcell_criteria"`
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
