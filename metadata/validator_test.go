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
			payload := metadata.Payload{
				FormTypes: []metadata.FormType{{}},
				JobTypes:  []metadata.JobType{{}},
				PropertyBlueprints: []metadata.PropertyBlueprint{{}},
			}
			messages, err := payload.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(messages).To(Equal(validator.ValidationErrorsTranslations{
				"Payload.FormTypes[0].Label":              "Label is a required field",
				"Payload.FormTypes[0].Name":               "Name is a required field",
				"Payload.FormTypes[0].PropertyInputs":     "PropertyInputs is a required field",
				"Payload.IconImage":                       "IconImage is a required field",
				"Payload.JobTypes[0].MaxInFlight":         "MaxInFlight is a required field",
				"Payload.JobTypes[0].Name":                "Name is a required field",
				"Payload.JobTypes[0].ResourceDefinitions": "ResourceDefinitions is a required field",
				"Payload.JobTypes[0].ResourceLabel":       "ResourceLabel is a required field",
				"Payload.JobTypes[0].SingleAZOnly":        "SingleAZOnly is a required field",
				"Payload.JobTypes[0].Templates":           "Templates is a required field",
				"Payload.MinimumVersionForUpgrade":        "MinimumVersionForUpgrade is a required field",
				"Payload.Name":                            "Name is a required field",
				"Payload.ProductVersion":                  "ProductVersion is a required field",
				"Payload.PropertyBlueprints[0].Name":      "Name is a required field",
				"Payload.PropertyBlueprints[0].Type":      "Type is a required field",
				"Payload.Releases":                        "Releases is a required field",
				"Payload.StemcellCriteria.OS":             "OS is a required field",
				"Payload.StemcellCriteria.Version":        "Version is a required field",
			}))
		})
	})
})
