package generator_test

import (
	"github.com/jtarchie/generate-tile/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generating the tile", func() {
	When("provided specs", func() {
		It("generates a set of properties and forms", func() {
			dir := createReleaseDir()
			specs, err := generator.ParseSpecs(dir)
			Expect(err).NotTo(HaveOccurred())

			tile, err := generator.GeneratorTile(specs)
			Expect(err).NotTo(HaveOccurred())

			Expect(tile.Description).To(Equal(""))
			ft := tile.FormTypes[0]
			Expect(ft.Name).To(Equal("properties"))
			Expect(ft.Label).To(Equal("Properties"))
			Expect(ft.Description).To(Equal("Configuration settings for properties"))
			Expect(ft.PropertyInputs).To(Equal([]generator.PropertyInput{
				{
					Reference:   ".properties.no_namespace",
					Label:       "No Namespace",
					Description: `This property has no namespace (".") in it.`,
				},
			}))

			ft = tile.FormTypes[1]
			Expect(ft.Name).To(Equal("some"))
			Expect(ft.Label).To(Equal("Some"))
			Expect(ft.Description).To(Equal("Configuration settings for some"))
			Expect(ft.PropertyInputs).To(Equal([]generator.PropertyInput{
				{
					Reference:   ".properties.some.property",
					Label:       "Some Property",
					Description: "This property is important for something.",
				},
				{
					Reference:   ".properties.some.tls_property",
					Label:       "Some Tls Property",
					Description: "This is a property with a type.",
				},
			}))

			pb := tile.PropertyBlueprints
			Expect(pb[0].Name).To(Equal(".properties.no_namespace"))
			Expect(pb[0].Configurable).To(BeTrue())
			Expect(pb[0].Optional).To(BeTrue())
			Expect(pb[0].Type).To(Equal("string"))

			Expect(pb[1].Name).To(Equal(".properties.some.property"))
			Expect(pb[1].Type).To(Equal("integer"))
			Expect(pb[1].Default).To(Equal(1))

			Expect(pb[2].Name).To(Equal(".properties.some.tls_property"))
			Expect(pb[2].Type).To(Equal("rsa_cert_credentials"))
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

			specs := []generator.SpecPayload{spec}
			tile, err := generator.GeneratorTile(specs)
			Expect(err).NotTo(HaveOccurred())

			pb := tile.PropertyBlueprints
			Expect(pb[0].Name).To(Equal(".properties.some.collection"))
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
