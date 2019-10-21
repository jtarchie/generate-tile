package generator_test

import (
	"fmt"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jtarchie/generate-tile/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
			Default:     1,
			Example:     100,
		}))
		Expect(spec.Properties).To(HaveKeyWithValue("some.tls_property", generator.Property{
			Description: "This is a property with a type.",
			Type:        "certificate",
		}))
	})

	It("parses a directory of spec files", func() {
		dir := createReleaseDir()

		release, err := generator.ParseRelease(dir)
		Expect(err).NotTo(HaveOccurred())

		specs := release.Specs
		Expect(specs).To(HaveLen(3))
		Expect(specs[0].Name).To(Equal("other"))
		Expect(specs[1].Name).To(Equal("some"))
		Expect(specs[2].Name).To(Equal("work"))

		Expect(release.Name).To(Equal("my-release"))
		Expect(release.LatestVersion).To(Equal("1.0.0"))
	})

	It("parses a bosh release", func() {
		dir := createReleaseTarball()

		release, err := generator.ParseRelease(dir)
		Expect(err).NotTo(HaveOccurred())

		Expect(release.Name).To(Equal("my-release"))
		Expect(release.LatestVersion).To(Equal("1.0.0"))

		specs := release.Specs
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

	releasesPath := filepath.Join(dir, "releases")
	err = os.MkdirAll(releasesPath, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(filepath.Join(releasesPath, "example-0.0.1.yml"), []byte(`{name: my-release, version: 0.0.1}`), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	err = ioutil.WriteFile(filepath.Join(releasesPath, "example-1.0.0.yml"), []byte(`{name: my-release, version: 1.0.0}`), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	return dir
}

func createReleaseTarball() string {
	buildDir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	jobNames := []string{"some", "other", "work"}
	for _, jobName := range jobNames {
		path := filepath.Join(buildDir, "jobs", jobName)
		err := os.MkdirAll(path, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		jobMF := filepath.Join(path, "job.MF")
		err = ioutil.WriteFile(jobMF, []byte(specYAML(jobName)), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = archiver.Archive(
			[]string{path},
			filepath.Join(buildDir, "jobs", fmt.Sprintf("%s.tgz", jobName)),
		)
		Expect(err).NotTo(HaveOccurred())

		err = os.RemoveAll(path)
		Expect(err).NotTo(HaveOccurred())
	}

	err = ioutil.WriteFile(filepath.Join(buildDir, "release.MF"), []byte(`{name: my-release, version: 1.0.0}`), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	releaseDir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	releasePath := filepath.Join(releaseDir, "release.tgz")
	err = archiver.Archive(
		[]string{buildDir},
		releasePath,
	)
	Expect(err).NotTo(HaveOccurred())

	fmt.Println(releasePath)
	return releasePath
}
