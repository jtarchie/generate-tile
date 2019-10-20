package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/imdario/mergo"
	"github.com/jtarchie/generate-tile/generator"
	"gopkg.in/yaml.v2"
)

func main() {
	dir := flag.String("dir", "", "directory to the bosh release")
	mergeFile := flag.String("merge", "", "yaml file to merge results with")

	flag.Parse()

	if *dir == "" {
		log.Fatalf("please specify bosh release directory")
	}

	contents, err := parseRelease(*dir)
	if err != nil {
		log.Fatalf("tile creation failed: %s", err)
	}

	if *mergeFile != "" {
		contents, err = mergeWithContents(*mergeFile, contents)
		if err != nil {
			log.Fatalf("cannot merge file: %s", err)
		}
	}

	fmt.Printf("%s", contents)
}

func mergeWithContents(file string, currentTileContents []byte) ([]byte, error) {
	var currentTile map[string]interface{}

	err := yaml.Unmarshal(currentTileContents, &currentTile)
	if err != nil {
		return nil, err
	}

	var mergeFile map[string]interface{}

	mergeFileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(mergeFileContents, &mergeFile)
	if err != nil {
		return nil, err
	}

	err = mergo.Merge(&currentTile, mergeFile)
	if err != nil {
		return nil, err
	}

	contents, err := yaml.Marshal(currentTile)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func parseRelease(releasePath string) ([]byte, error) {
	specs, err := generator.ParseSpecs(releasePath)
	if err != nil {
		return nil, err
	}

	tile, err := generator.GeneratorTile(specs)
	if err != nil {
		return nil, err
	}

	contents, err := yaml.Marshal(tile)
	if err != nil {
		return nil, err
	}

	return contents, nil
}