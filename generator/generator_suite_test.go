package generator_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generator Suite")
}

func writeFile(contents string) string {
	file, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())

	_, err = file.WriteString(contents)
	Expect(err).NotTo(HaveOccurred())

	err = file.Close()
	Expect(err).NotTo(HaveOccurred())

	return file.Name()
}

func specYAML(s string) string {
	return fmt.Sprintf(specYAMLTemplate, s)
}

const specYAMLTemplate = `
name: %s
description: >
  Example description in a spec.

templates:
  ctl.erb: bin/ctl
  start.erb: bin/start

packages:
- package1
- package2

provides:
- name: provided
  type: provided

consumes:
- name: provided
  type: provided
  optional: true

properties:
  some.property:
    description: This property is important for something.
    default: 1
    example: 100
  some.tls_property:
    description: This is a property with a type.
    type: certificate
  no_namespace:
    description: This property has no namespace (".") in it.
`
