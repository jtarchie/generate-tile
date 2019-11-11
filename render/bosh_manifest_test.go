package render_test

import (
	"github.com/jtarchie/tile-builder/metadata"
	"github.com/jtarchie/tile-builder/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BoshManifest", func() {
	When("a valid product config is provided", func() {
		var (
			err      error
			manifest *render.Deployment
		)

		BeforeEach(func() {
			manifest, err = render.AsBoshManifest(metadata.Payload{
				Name: "some-tile",
				Releases: []metadata.Release{
					{
						File:    "some-release-v1.0.0.tgz",
						Name:    "some-release",
						Version: "v1.0.0",
						SHA1:    "some-sha1",
					},
					{
						File:    "another-release-v1.0.0.tgz",
						Name:    "another-release",
						Version: "v1.0.0",
						SHA1:    "another-sha1",
					},
				},
				StemcellCriteria: metadata.StemcellCriteria{
					OS:                         "ubuntu-xenial",
					Version:                    "319.70",
				},
				AdditionalStemcellsCriteria: []metadata.StemcellCriteria{
					{
						OS:                         "windows-2019",
						Version:                    "12.3",
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a deployment name", func() {
			Expect(manifest.Name).To(Equal("some-tile-guid"))
		})

		It("contains the releases", func() {
			releases := manifest.Releases
			Expect(releases).To(Equal([]render.Release{
				{
					Name:    "another-release",
					Version: "v1.0.0",
					SHA1:    "another-sha1",
				},
				{
					Name:    "some-release",
					Version: "v1.0.0",
					SHA1:    "some-sha1",
				},
			}))
		})

		It("contains the stemcells", func() {
			stemcells := manifest.Stemcells
			Expect(stemcells).To(Equal([]render.Stemcell{
				{
					Alias:   "ubuntu-xenial",
					OS:      "ubuntu-xenial",
					Version: "((ubuntu-xenial-version))",
				},
				{
					Alias:   "windows-2019",
					OS:      "windows-2019",
					Version: "((windows-2019-version))",
				},
			}))
		})
	})
})
