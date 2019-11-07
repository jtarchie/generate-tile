package generator_test

import (
	"github.com/jtarchie/tile-builder/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Determining the property default value", func() {
	DescribeTable("determine default", func(property generator.Property, expectedType types.GomegaMatcher) {
		def, err := generator.DeterminePropertyBlueprintDefault("", property)
		Expect(err).NotTo(HaveOccurred())
		Expect(def).To(expectedType)
	},
		Entry("default is a boolean", generator.Property{Default: false}, BeFalse()),
		Entry("default is a boolean", generator.Property{Default: true}, BeTrue()),
		Entry("default is a integer", generator.Property{Default: 1}, Equal(1)),
		Entry("default is a string", generator.Property{Default: "asdf"}, Equal("asdf")),
		Entry("default is a array of strings", generator.Property{Default: []interface{}{"a", "b"}}, Equal("a,b")),
		Entry("default is a map", generator.Property{Default: map[interface{}]interface{}{}}, BeNil()),
	)
})
