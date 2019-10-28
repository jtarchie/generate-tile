package generator_test

import (
	"github.com/jtarchie/generate-tile/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Determining the property default value", func() {
	DescribeTable("determine default", func(property generator.Property, expectedType interface{}) {
		Expect(generator.DeterminePropertyBlueprintDefault(property)).To(Equal(expectedType))
	},
		Entry("default is a boolean", generator.Property{Default: false}, false),
		Entry("default is a boolean", generator.Property{Default: true}, true),
		Entry("default is a integer", generator.Property{Default: 1}, 1),
		Entry("default is a string", generator.Property{Default: "asdf"}, "asdf"),
	)
})
