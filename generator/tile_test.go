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

			ft = tile.FormTypes[1]
			Expect(ft.Name).To(Equal("some"))
			Expect(ft.Label).To(Equal("Some"))
			Expect(ft.Description).To(Equal("Configuration settings for some"))
		})
	})
})
