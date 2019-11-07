package metadata_test

import (
	"github.com/jtarchie/tile-builder/metadata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/go-playground/validator.v9"
)

var _ = Describe("Validating a tile's metadata", func() {
	It("requires fields", func() {
		payload := metadata.Payload{
			FormTypes:          []metadata.FormType{{}},
			JobTypes:           []metadata.JobType{{}},
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

	It("requires blueprint types", func() {
		payload := metadata.Payload{
			PropertyBlueprints: []metadata.PropertyBlueprint{
				{
					Name: "not-a-string",
					Type: "not-a-string",
				},
			},
		}
		messages, err := payload.Validate()
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(HaveKeyWithValue(
			"Payload.PropertyBlueprints[0].Type",
			"Type must be one of [boolean ca_certificate collection disk_type_dropdown domain dropdown_select email http_url integer ip_address ip_ranges ldap_url multi_select_options network_address network_address_list port rsa_cert_credentials rsa_pkey_credentials salted_credentials secret selector service_network_az_multi_select service_network_az_single_select simple_credentials smtp_authentication stemcell_selector string_list string text uuid vm_type_dropdown wildcard_domain]",
		))
	})

	It("requires that a property input references a property blueprint", func() {
		payload := metadata.Payload{
			FormTypes: []metadata.FormType{
				{
					PropertyInputs: []metadata.PropertyInput{
						{
							Reference: ".properties.name",
						},
					},
				},
			},
			PropertyBlueprints: []metadata.PropertyBlueprint{},
		}
		messages, err := payload.Validate()
		Expect(err).NotTo(HaveOccurred())
		Expect(messages).To(HaveKeyWithValue(
			"Payload.FormTypes[0].PropertyInputs[0].Reference",
			"References a property blueprint ('.properties.name') that does not exist",
		))
	})
})