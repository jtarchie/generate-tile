package generator_test

import (
	"github.com/jtarchie/generate-tile/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/gomega"
)

var _ = Describe("Determining the property type", func() {
	DescribeTable("determine type", func(name string, property generator.Property, expectedType string) {
		actualType, err := generator.DeterminePropertyBlueprintType(name, property)
		Expect(err).NotTo(HaveOccurred())
		Expect(actualType).To(Equal(expectedType))
	},
		Entry("defaults to string", "", generator.Property{}, "string"),

		Entry("when default is a string", "", generator.Property{Default: ""}, "string"),
		Entry("when default is a integer", "", generator.Property{Default: 1}, "integer"),
		Entry("when default is a float", "", generator.Property{Default: 1.0}, "integer"),
		Entry("when default is a boolean", "", generator.Property{Default: false}, "boolean"),
		Entry("when default is a empty array", "", generator.Property{Default: []interface{}{}}, "string_list"),

		Entry("when example is a string", "", generator.Property{Example: ""}, "string"),
		Entry("when example is a integer", "", generator.Property{Example: 1}, "integer"),
		Entry("when example is a float", "", generator.Property{Example: 1.0}, "integer"),
		Entry("when example is a boolean", "", generator.Property{Example: false}, "boolean"),
		Entry("when example is a empty array", "", generator.Property{Example: []interface{}{}}, "string_list"),

		Entry("when type is a certificate", "", generator.Property{Type: "certificate"}, "rsa_cert_credentials"),
		Entry("when type is a rsa", "", generator.Property{Type: "rsa"}, "rsa_pkey_credentials"),
		Entry("when type is a ssh", "", generator.Property{Type: "ssh"}, "rsa_pkey_credentials"),

		Entry("when the name includes _port", "server_port", generator.Property{}, "port"),
		Entry("when the name includes .port", "server.port", generator.Property{}, "port"),
		Entry("when the name includes port", "serverport", generator.Property{}, "string"),
		Entry("when the name includes _port_", "server_port_", generator.Property{}, "string"),

		Entry("when the name includes _ip", "server_ip", generator.Property{}, "ip_address"),
		Entry("when the name includes .ip", "server.ip", generator.Property{}, "ip_address"),
		Entry("when the name includes ip", "serverip", generator.Property{}, "string"),
		Entry("when the name includes _ip_", "server_ip_", generator.Property{}, "string"),

		Entry("when the example is a map", "", generator.Property{Example: map[interface{}]interface{}{
			"name": "value",
		}}, "collection"),
	)
})
