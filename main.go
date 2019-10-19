package main

import (
	"flag"
	"fmt"
	"github.com/jtarchie/generate-tile/generator"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

func main() {
	dir := flag.String("dir", "", "directory to the bosh release")
	flag.Parse()

	if *dir == "" {
		log.Fatalf("please specify bosh release directory")
	}

	err := parseRelease(*dir)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
}

func parseRelease(releasePath string) error {
	specs, err := generator.ParseSpecs(releasePath)
	if err != nil {
		return err
	}

	tile, err := generator.GeneratorTile(specs)
	if err != nil {
		return err
	}

	contents, err := yaml.Marshal(tile)
	if err != nil {
		return err
	}

	fmt.Printf("%s", contents)
	return nil
}

func propertyNameToLabel(name string) string {
	var labels []string

	names := strings.Split(name, ".")
	for _, n := range names {
		labels = append(labels, strings.Split(n, "_")...)
	}
	return strings.Join(labels, " ")
}
