package generator_test

import (
	"github.com/jtarchie/generate-tile/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Parsing spec files", func() {
	It("parses a bosh spec into struct", func() {
		spec, err := generator.ParseSpec(writeFile(specYAML("example")))
		Expect(err).NotTo(HaveOccurred())

		Expect(spec.Name).To(Equal("example"))
		Expect(spec.Description).To(Equal("Example description in a spec.\n"))
		Expect(spec.Templates).To(HaveKeyWithValue("ctl.erb", "bin/ctl"))
		Expect(spec.Templates).To(HaveKeyWithValue("start.erb", "bin/start"))
		Expect(spec.Packages).To(Equal([]string{"package1", "package2"}))
		Expect(spec.Properties).To(HaveKeyWithValue("some.property", generator.Property{
			Description: "This property is important for something.",
			Default:     "INFO",
			Example:     "INFO | ERROR",
		}))
		Expect(spec.Properties).To(HaveKeyWithValue("some.tls_property", generator.Property{
			Description: "This is a property with a type.",
			Type:        "certificate",
		}))
	})

	It("parses a directory of spec files", func() {
		dir := createReleaseDir()

		specs, err := generator.ParseSpecs(dir)
		Expect(err).NotTo(HaveOccurred())

		Expect(specs).To(HaveLen(3))
		Expect(specs[0].Name).To(Equal("other"))
		Expect(specs[1].Name).To(Equal("some"))
		Expect(specs[2].Name).To(Equal("work"))
	})
})

func createReleaseDir() string {
	dir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	jobNames := []string{"some", "other", "work"}
	for _, jobName := range jobNames {
		path := filepath.Join(dir, "jobs", jobName)
		err := os.MkdirAll(path, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(path, "spec"), []byte(specYAML(jobName)), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())
	}

	return dir
}
