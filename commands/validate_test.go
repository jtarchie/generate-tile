package commands_test

import (
	"github.com/jtarchie/generate-tile/commands"
	"github.com/mholt/archiver"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
)

var _ = Describe("Validate", func() {
	It("writes validation errors to stdout", func() {
		stdout := gbytes.NewBuffer()
		productPath := createProductFile()
		command := commands.Validate{
			Path: productPath,
			Stdout: stdout,
		}
		err := command.Execute(nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(stdout).To(gbytes.Say("Payload.Name: Name is a required field"))
	})
})

func createProductFile() string {
	dir, err := ioutil.TempDir("", "")
	Expect(err).NotTo(HaveOccurred())

	productPath := filepath.Join(dir, "product.pivotal")

	metadataPath := filepath.Join(dir, "metadata")
	err = os.MkdirAll(metadataPath, os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	metadataFile, err := os.Create(filepath.Join(metadataPath, "metadata.yml"))
	Expect(err).NotTo(HaveOccurred())

	_, err = metadataFile.Write([]byte(``))
	Expect(err).NotTo(HaveOccurred())
	err = metadataFile.Close()
	Expect(err).NotTo(HaveOccurred())

	productZip := filepath.Join(dir, "product.zip")

	zip := archiver.NewZip()
	err = zip.Archive([]string{
		metadataPath,
	}, productZip)
	Expect(err).NotTo(HaveOccurred())

	err = os.Rename(productZip, productPath)
	Expect(err).NotTo(HaveOccurred())

	return productPath
}
