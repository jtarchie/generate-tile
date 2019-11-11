package render

import (
	"fmt"
	"github.com/jtarchie/tile-builder/manifest"
	"github.com/jtarchie/tile-builder/metadata"
	"sort"
)

func AsBoshManifest(payload metadata.Payload) (*manifest.Payload, error) {
	deployment := &manifest.Payload{
		Name: fmt.Sprintf("%s-guid", payload.Name),
	}

	addReleases(payload, deployment)
	addStemcells(payload, deployment)

	return deployment, nil
}

func addStemcells(payload metadata.Payload, deployment *manifest.Payload) {
	stemcells := []metadata.StemcellCriteria{payload.StemcellCriteria}
	stemcells = append(stemcells, payload.AdditionalStemcellsCriteria...)
	sort.Slice(stemcells, func(i, j int) bool {
		return stemcells[i].OS < stemcells[j].OS
	})

	for _, stemcell := range stemcells {
		deployment.Stemcells = append(deployment.Stemcells, manifest.Stemcell{
			Alias:   stemcell.OS,
			OS:      stemcell.OS,
			Version: fmt.Sprintf("((%s-version))", stemcell.OS),
		})
	}
}

func addReleases(payload metadata.Payload, deployment *manifest.Payload) {
	releases := payload.Releases
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].Name < releases[j].Name
	})

	for _, release := range releases {
		deployment.Releases = append(deployment.Releases, manifest.Release{
			Name:    release.Name,
			Version: release.Version,
			SHA1:    release.SHA1,
		})
	}
}
