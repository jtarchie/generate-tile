package render_test

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/jtarchie/tile-builder/metadata"
	"github.com/jtarchie/tile-builder/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = PDescribe("AsHTML", func() {
	It("returns tabs for each form group", func() {
		doc := renderMetadata()

		Expect(doc.Find(".nav .nav-link").First().Text()).To(Equal("boolean"))
	})

	It("generates fields based on type", func() {
		doc := renderMetadata()
		Expect(doc.Find(`.form-check input.form-check-input#boolean[type="checkbox"]`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group input#integer[type="number"]`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group input#ip_address[type="text"][pattern]`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group input#port[type="number"][min="0"]`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group textarea#rsa_cert_credentials_certificate`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group textarea#rsa_cert_credentials_private_key`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group textarea#rsa_pkey_credentials_private_key`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group textarea#rsa_pkey_credentials_public_key`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group input#string[type="text"]`).Length()).To(Equal(1))
		Expect(doc.Find(`.form-group input#string_list[type="text"]`).Length()).To(Equal(1))
	})
})

func renderMetadata() *goquery.Document {
	payload := metadata.Payload{}
	for _, t := range []string{
		"boolean",
		"integer",
		"ip_address",
		"port",
		"rsa_cert_credentials",
		"rsa_pkey_credentials",
		"string",
		"string_list",
	} {
		payload.FormTypes = append(payload.FormTypes, metadata.FormType{
			Label: t,
			PropertyInputs: []metadata.PropertyInput{
				{
					Reference: fmt.Sprintf(".properties.%s", t),
				},
			},
		})
		payload.PropertyBlueprints = append(payload.PropertyBlueprints, metadata.PropertyBlueprint{
			Name: t,
			Type: t,
		})
	}

	contents, err := render.AsHTML(payload)
	Expect(err).NotTo(HaveOccurred())

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(contents))
	Expect(err).NotTo(HaveOccurred())

	return doc
}
