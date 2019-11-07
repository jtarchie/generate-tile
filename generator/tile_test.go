package generator_test

import (
	"github.com/jtarchie/tile-builder/generator"
	tile2 "github.com/jtarchie/tile-builder/metadata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generating the tile", func() {
	When("provided specs", func() {
		It("generates a set of properties and forms", func() {
			dir := createReleaseDir()
			release, err := generator.ParseRelease(dir)
			Expect(err).NotTo(HaveOccurred())

			tile, err := generator.GeneratorTile(release)
			Expect(err).NotTo(HaveOccurred())

			Expect(tile.Description).To(Equal(""))
			ft := tile.FormTypes[0]
			Expect(ft.Name).To(Equal("properties"))
			Expect(ft.Label).To(Equal("Properties"))
			Expect(ft.Description).To(Equal("Configuration settings for Properties"))
			Expect(ft.PropertyInputs).To(Equal([]tile2.PropertyInput{
				{
					Reference:   ".properties.no_namespace",
					Label:       "No Namespace",
					Description: `This property has no namespace (".") in it.`,
				},
			}))

			ft = tile.FormTypes[1]
			Expect(ft.Name).To(Equal("some"))
			Expect(ft.Label).To(Equal("Some"))
			Expect(ft.Description).To(Equal("Configuration settings for Some"))
			Expect(ft.PropertyInputs).To(Equal([]tile2.PropertyInput{
				{
					Reference:   ".properties.some__property",
					Label:       "Some Property",
					Description: "This property is important for something.",
				},
				{
					Reference:   ".properties.some__tls_property",
					Label:       "Some Tls Property",
					Description: "This is a property with a type.",
				},
			}))

			ft = tile.FormTypes[2]
			Expect(ft.Name).To(Equal("some_long"))
			Expect(ft.Label).To(Equal("Some Long"))

			pb := tile.PropertyBlueprints
			Expect(pb[0].Name).To(Equal("no_namespace"))
			Expect(pb[0].Configurable).To(BeTrue())
			Expect(pb[0].Optional).To(BeTrue())
			Expect(pb[0].Type).To(Equal("string"))

			Expect(pb[1].Name).To(Equal("some__property"))
			Expect(pb[1].Type).To(Equal("integer"))
			Expect(pb[1].Default).To(Equal(1))
			Expect(pb[1].Optional).To(BeFalse())

			Expect(pb[2].Name).To(Equal("some__tls_property"))
			Expect(pb[2].Type).To(Equal("rsa_cert_credentials"))

			jobs := tile.JobTypes
			Expect(jobs[0].Name).To(Equal("other"))
			Expect(jobs[0].ResourceLabel).To(Equal("Other"))
			Expect(jobs[0].SingleAZOnly).To(BeFalse())
			Expect(jobs[0].UseStemcell).To(Equal(""))
			Expect(jobs[0].InstanceDefinition).To(Equal(tile2.InstanceDefinition{
				Name:         "instances",
				Label:        "Instances",
				Configurable: true,
				Constraints: tile2.Constraints{
					Min: 1,
				},
				Default: 1,
				Type:    "integer",
			}))
			Expect(jobs[0].ResourceDefinitions).To(Equal([]tile2.ResourceDefinition{
				{
					Name:         "cpu",
					Configurable: true,
					Default:      1,
					Constraints: tile2.Constraints{
						Min: 1,
					},
					Label: "CPU",
					Type:  "integer",
				},
				{
					Name:         "ram",
					Configurable: true,
					Default:      8192,
					Constraints: tile2.Constraints{
						Min: 8192,
					},
					Label: "RAM",
					Type:  "integer",
				},
				{
					Name:         "ephemeral_disk",
					Configurable: true,
					Default:      10240,
					Constraints: tile2.Constraints{
						Min: 10240,
					},
					Label: "Ephemeral Disk",
					Type:  "integer",
				},
				{
					Name:         "persistent_disk",
					Configurable: true,
					Default:      10240,
					Constraints: tile2.Constraints{
						Min: 10240,
					},
					Label: "Persistent Disk",
					Type:  "integer",
				},
			}))
			Expect(jobs[0].Manifest).To(MatchYAML(`
---
no_namespace: ((.properties.no_namespace.value))
some_long:
  property: ((.properties.some_long__property.value))
some:
  property: ((.properties.some__property.value))
  tls_property:
    certificate: ((.properties.some__tls_property.certificate))
    private_key: ((.properties.some__tls_property.private_key))
`))

			Expect(jobs[0].MaxInFlight).To(Equal(1))

			Expect(jobs[0].Templates[0].Name).To(Equal("other"))
			Expect(jobs[0].Templates[0].Release).To(Equal("my-release"))
			Expect(jobs[0].Templates[0].Consumes).To(MatchYAML("provided: {from: some-provided}"))
			Expect(jobs[0].Templates[0].Provides).To(MatchYAML("provided: {as: other-provided}"))

			Expect(jobs[1].Name).To(Equal("some"))
			Expect(jobs[2].Name).To(Equal("work"))
		})
	})

	When("there is a collection in the spec", func() {
		const specWithCollection = `
name: example

properties:
  some.collection:
    example:
      env: dev
      foo: bar
`
		It("creates the correct property", func() {
			spec, err := generator.ParseSpec(writeFile(specWithCollection))
			Expect(err).NotTo(HaveOccurred())

			release := generator.BoshReleasePayload{
				Specs: []generator.SpecPayload{spec},
			}

			tile, err := generator.GeneratorTile(release)
			Expect(err).NotTo(HaveOccurred())

			pb := tile.PropertyBlueprints
			Expect(pb[0].Name).To(Equal("some__collection"))
			Expect(pb[0].PropertyBlueprints[0].Name).To(Equal("key"))
			Expect(pb[0].PropertyBlueprints[0].Type).To(Equal("string"))
			Expect(pb[0].PropertyBlueprints[0].Configurable).To(BeTrue())
			Expect(pb[0].PropertyBlueprints[0].Optional).To(BeTrue())

			Expect(pb[0].PropertyBlueprints[1].Name).To(Equal("value"))
			Expect(pb[0].PropertyBlueprints[1].Type).To(Equal("string"))
			Expect(pb[0].PropertyBlueprints[1].Configurable).To(BeTrue())
			Expect(pb[0].PropertyBlueprints[1].Optional).To(BeTrue())
		})
	})
})
