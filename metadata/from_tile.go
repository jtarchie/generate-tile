package metadata

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/mholt/archiver"
	"gopkg.in/yaml.v2"
)

func FromTile(tilePath string, strict bool) (Payload, error) {
	var (
		contents bytes.Buffer
		payload  Payload
	)

	archive := archiver.NewZip()
	err := archive.Walk(tilePath, func(f archiver.File) error {
		zfh, ok := f.Header.(zip.FileHeader)
		if ok {
			if metadataFile.MatchString(zfh.Name) {
				_, err := io.Copy(&contents, f)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return payload, fmt.Errorf("could not find metadata file in %s: %s", tilePath, err)
	}

	if strict {
		err = yaml.UnmarshalStrict(contents.Bytes(), &payload)
	} else {
		err = yaml.Unmarshal(contents.Bytes(), &payload)
	}
	if err != nil {
		return payload, fmt.Errorf("could not unmarshal %s: %s", tilePath, err)
	}

	return payload, nil
}