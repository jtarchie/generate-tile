package render

import (
	"fmt"
	"github.com/jtarchie/tile-builder/metadata"
	"sort"
)

type ReleaseStemcell struct {
	Name    string `validate:"required"`
	Version string `validate:"required"`
}

type Release struct {
	Name         string            `validate:"required"`
	Version      string            `validate:"required"`
	URL          string            `yaml:",omitempty"`
	SHA1         string            `yaml:"sha1,omitempty"`
	Stemcell     []ReleaseStemcell `yaml:",omitempty" validate:"dive"`
	ExportedFrom []ReleaseStemcell `yaml:"exported_from,omitempty" validate:"dive"`
}

type Stemcell struct {
	Alias   string `validate:"required"`
	OS      string `yaml:",omitempty" validate:"required_without=Name"`
	Version string `validate:"required"`
	Name    string `yaml:",omitempty" validate:"required_without=OS"`
}

type Features struct {
	ConvergeVariables    bool `yaml:"converge_variables"`
	RandomizeAZPlacement bool `yaml:"randomize_az_placement"`
	UseDNSAddresses      bool `yaml:"use_dns_addresses"`
	UseTmpfsConfig       bool `yaml:"use_tmpfs_config"`
	UseShortDNSAddresses bool `yaml:"use_short_dns_addresses"`
}

type Deployment struct {
	Name      string
	Features  Features   `yaml:",omitempty"`
	Releases  []Release  `validate:"required,dive"`
	Stemcells []Stemcell `validate:"required,dive"`
}

func AsBoshManifest(payload metadata.Payload) (*Deployment, error) {
	manifest := &Deployment{
		Name: fmt.Sprintf("%s-guid", payload.Name),
	}

	addReleases(payload, manifest)
	addStemcells(payload, manifest)

	return manifest, nil
}

func addStemcells(payload metadata.Payload, deployment *Deployment) {
	stemcells := []metadata.StemcellCriteria{payload.StemcellCriteria}
	stemcells = append(stemcells, payload.AdditionalStemcellsCriteria...)
	sort.Slice(stemcells, func(i, j int) bool {
		return stemcells[i].OS < stemcells[j].OS
	})

	for _, stemcell := range stemcells {
		deployment.Stemcells = append(deployment.Stemcells, Stemcell{
			Alias:   stemcell.OS,
			OS:      stemcell.OS,
			Version: fmt.Sprintf("((%s-version))", stemcell.OS),
		})
	}
}

func addReleases(payload metadata.Payload, deployment *Deployment) {
	releases := payload.Releases
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Name < releases[j].Name
	})

	for _, release := range releases {
		deployment.Releases = append(deployment.Releases, Release{
			Name:    release.Name,
			Version: release.Version,
			SHA1:    release.SHA1,
		})
	}
}
