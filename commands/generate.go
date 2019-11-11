package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/imdario/mergo"
	"github.com/jtarchie/tile-builder/generator"
	"gopkg.in/yaml.v2"
)

type Generate struct {
	Path        string `long:"path" required:"true" description:"path to the bosh release (source directory or tarball)"`
	MergingFile string `long:"merge" description:"yaml file to merge results with"`
}

func (g Generate) Execute(_ []string) error {
	contents, err := parseRelease(g.Path)
	if err != nil {
		return fmt.Errorf("tile creation failed: %s", err)
	}

	if g.MergingFile != "" {
		contents, err = mergeWithContents(g.MergingFile, contents)
		if err != nil {
			return fmt.Errorf("cannot merge file: %s", err)
		}
	}

	fmt.Printf("%s", contents)

	return nil
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
	specs, err := generator.ParseRelease(releasePath)
	if err != nil {
		return nil, err
	}

	tile, err := generator.Tile(specs)
	if err != nil {
		return nil, err
	}

	contents, err := yaml.Marshal(tile)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
