package generator

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"sort"
)

type Property struct {
	Env         string
	Description string
	Default     interface{}
	Example     interface{}
	Type        string
	EnvFile     string `yaml:"env_file"`
	EnvFields   map[string]struct {
		EnvFile string `yaml:"env_file"`
	} `yaml:"env_fields"`
}

type SpecPayload struct {
	Name        string
	Description string
	Templates   map[string]string
	Packages    []string
	Consumes    []struct {
		Name     string
		Type     string
		Optional bool
	}
	Provides []struct {
		Name       string
		Type       string
		Properties []string
	}
	Properties map[string]Property `yaml:"properties"`
}

func ParseSpec(filename string) (SpecPayload, error) {
	var spec SpecPayload

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return spec, fmt.Errorf("could not read contentes of %s: %s", filename, err)
	}

	err = yaml.UnmarshalStrict(contents, &spec)
	if err != nil {
		return spec, fmt.Errorf("could unmarshal spec from %s: %s", filename, err)
	}

	return spec, nil
}

func ParseSpecs(releasePath string) ([]SpecPayload, error) {
	matches, err := filepath.Glob(filepath.Join(releasePath, "jobs", "*"))
	if err != nil {
		return nil, fmt.Errorf("could not find the release's jobs in %s: %s", releasePath, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no jobs found in release in %s", releasePath)
	}

	sort.Strings(matches)

	var specs []SpecPayload
	for _, jobPath := range matches {
		specPath := filepath.Join(jobPath, "spec")

		spec, err := ParseSpec(specPath)
		if err != nil {
			return nil, fmt.Errorf("could not open spec of the job %s: %s", specPath, err)
		}

		specs = append(specs, spec)
	}

	return specs, nil
}