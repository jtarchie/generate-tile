package generator

import (
	"fmt"
	"github.com/mholt/archiver"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar"

	"gopkg.in/yaml.v2"
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

type consumePayload struct {
	Name     string
	Type     string
	Optional bool
}

type providerPayload struct {
	Name       string
	Type       string
	Properties []string
}

type SpecPayload struct {
	Name        string
	Description string
	Templates   map[string]string
	Packages    []string
	Consumes    []consumePayload
	Provides    []providerPayload
	Properties  map[string]Property `yaml:"properties"`
}

type BoshReleasePayload struct {
	Specs         []SpecPayload
	Name          string
	LatestVersion string
}

func ParseSpec(filename string) (SpecPayload, error) {
	var (
		spec     SpecPayload
		contents []byte
		err      error
	)

	contents = []byte(filename)

	if _, err := os.Stat(filename); err == nil {
		contents, err = ioutil.ReadFile(filename)
		if err != nil {
			return spec, fmt.Errorf("could not read contents of %s: %s", filename, err)
		}
	}

	err = yaml.UnmarshalStrict(contents, &spec)
	if err != nil {
		return spec, fmt.Errorf("could unmarshal spec from %s: %s", filename, err)
	}

	return spec, nil
}

type ReleasePayload struct {
	Name               string
	Version            string
	CommitHash         string `yaml:"commit_hash"`
	UncommittedChanges bool   `yaml:"uncommitted_changes"`
	Jobs               []struct {
		Name        string
		Version     string
		Fingerprint string
		Sha1        string `yaml:"sha1"`
	}
	Packages []struct {
		Name         string
		Version      string
		Fingerprint  string
		Sha1         string `yaml:"sha1"`
		Dependencies interface{}
	}
	License struct {
		Version     string
		Fingerprint string
		Sha1        string `yaml:"sha1"`
	}
}

func ParseRelease(releasePath string) (BoshReleasePayload, error) {
	info, err := os.Stat(releasePath)
	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not state release %s: %s", releasePath, err)
	}

	var boshRelease BoshReleasePayload
	if info.IsDir() {
		boshRelease, err = parseReleaseDir(releasePath)
		if err != nil {
			return BoshReleasePayload{}, fmt.Errorf("could not parse directory: %s", err)
		}
	} else {
		boshRelease, err = parseReleaseTarball(releasePath)
		if err != nil {
			return BoshReleasePayload{}, fmt.Errorf("could not parse bosh release: %s", err)
		}
	}

	return boshRelease, nil
}

func parseReleaseTarball(releasePath string) (BoshReleasePayload, error) {
	var boshRelease BoshReleasePayload

	err := archiver.Walk(releasePath, func(f archiver.File) error {
		if strings.Contains(f.Name(), "release.MF") {
			contents, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("could not read contents of release %s: %s", f.Name(), err)
			}

			var release ReleasePayload

			err = yaml.UnmarshalStrict(contents, &release)
			if err != nil {
				return fmt.Errorf("could not unmarshal release %s: %s", f.Name(), err)
			}

			boshRelease.Name = release.Name
			boshRelease.LatestVersion = release.Version
		}

		if strings.Contains(f.Name(), "job.MF") {
			contents, err := ioutil.ReadAll(f)
			if err != nil {
				return fmt.Errorf("could not read contents of spec %s: %s", f.Name(), err)
			}

			spec, err := ParseSpec(string(contents))
			if err != nil {
				return fmt.Errorf("could not open spec of the job %s: %s", f.Name(), err)
			}

			boshRelease.Specs = append(boshRelease.Specs, spec)
		}
		return nil
	})

	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not unarchive %s: %s", releasePath, err)
	}

	return boshRelease, nil
}

func parseReleaseDir(releasePath string) (BoshReleasePayload, error) {
	matches, err := filepath.Glob(filepath.Join(releasePath, "jobs", "*"))
	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not find the release's jobs in %s: %s", releasePath, err)
	}
	if len(matches) == 0 {
		return BoshReleasePayload{}, fmt.Errorf("no jobs found in release in %s", releasePath)
	}

	sort.Strings(matches)

	var specs []SpecPayload
	for _, jobPath := range matches {
		specPath := filepath.Join(jobPath, "spec")

		spec, err := ParseSpec(specPath)
		if err != nil {
			return BoshReleasePayload{}, fmt.Errorf("could not open spec of the job %s: %s", specPath, err)
		}

		specs = append(specs, spec)
	}

	var boshRelease BoshReleasePayload
	boshRelease.Specs = specs

	matches, err = doublestar.Glob(filepath.Join(releasePath, "releases", "**", "*-*.yml"))
	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not find the release's releases in %s: %s", releasePath, err)
	}
	if len(matches) == 0 {
		return BoshReleasePayload{}, fmt.Errorf("no releases found in release in %s", releasePath)
	}

	sort.Strings(matches)

	latestRelease := matches[len(matches)-1]
	contents, err := ioutil.ReadFile(latestRelease)
	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not open release %s: %s", latestRelease, err)
	}

	var release ReleasePayload

	err = yaml.UnmarshalStrict(contents, &release)
	if err != nil {
		return BoshReleasePayload{}, fmt.Errorf("could not unmarshal release %s: %s", latestRelease, err)
	}

	boshRelease.Name = release.Name
	boshRelease.LatestVersion = release.Version

	return boshRelease, nil
}
