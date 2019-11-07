package metadata_test

import (
	"github.com/jtarchie/generate-tile/metadata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/go-playground/validator.v9"
)

var _ = Describe("Validating a tile's metadata", func() {
	When("give an invalid tile", func() {
		It("requires fields", func() {
			payload := metadata.Payload{}
			messages, err := payload.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(messages).To(Equal(validator.ValidationErrorsTranslations{
				"Payload.Name": "Name is a required field",
			}))
		})
	})
})
