package manifest_test

import (
	"github.com/jtarchie/tile-builder/manifest"
	. "github.com/onsi/ginkgo"
	"gopkg.in/go-playground/validator.v9"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	It("requires fields", func() {
		payload := manifest.Payload{}
		messages, err := payload.Validate()
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(Equal(validator.ValidationErrorsTranslations{
			"Payload.Name":      "Name is a required field",
			"Payload.Releases":  "Releases is a required field",
			"Payload.Stemcells": "Stemcells is a required field",
		}))

		payload = manifest.Payload{
			Releases:  []manifest.Release{{}},
			Stemcells: []manifest.Stemcell{{}},
		}
		messages, err = payload.Validate()
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(Equal(validator.ValidationErrorsTranslations{
			"Payload.Name":                 "Name is a required field",
			"Payload.Releases[0].Name":     "Name is a required field",
			"Payload.Releases[0].Version":  "Version is a required field",
			"Payload.Stemcells[0].Alias":   "Alias is a required field",
			"Payload.Stemcells[0].Name":    "Key: 'Payload.Stemcells[0].Name' Error:Field validation for 'Name' failed on the 'required_without' tag",
			"Payload.Stemcells[0].OS":      "Key: 'Payload.Stemcells[0].OS' Error:Field validation for 'OS' failed on the 'required_without' tag",
			"Payload.Stemcells[0].Version": "Version is a required field",
		}))
	})
})
