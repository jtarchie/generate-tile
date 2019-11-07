package render_test

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/jtarchie/tile-builder/metadata"
	"github.com/jtarchie/tile-builder/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AsHTML", func() {
	It("returns tabs for each form group", func() {
		doc := renderMetadata()

		Expect(doc.Find(".nav .nav-link").First().Text()).To(Equal("Some Properties"))
	})

	It("generates string type as field", func() {
		doc := renderMetadata()
		Expect(doc.Find(`.form-group input#string[type="text"]`).Length()).To(Equal(1))
	})
})

func renderMetadata() *goquery.Document {
	contents, err := render.AsHTML(metadata.Payload{
		FormTypes: []metadata.FormType{
			{
				Label: "Some Properties",
				PropertyInputs: []metadata.PropertyInput{
					{
						Reference: ".properties.string",
					},
				},
			},
		},
		PropertyBlueprints: []metadata.PropertyBlueprint{
			{
				Name: "string",
				Type: "string",
			},
		},
	})
	Expect(err).NotTo(HaveOccurred())
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	Expect(err).NotTo(HaveOccurred())

	return doc
}
