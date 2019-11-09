package metadata

import (
	"archive/zip"
	"fmt"
	"github.com/pivotal-cf/go-pivnet/v2"
	"github.com/pivotal-cf/go-pivnet/v2/logshim"
	"gopkg.in/yaml.v2"
	"howett.net/ranger"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type httpClient struct {
	pivnet.Client
}

func (h httpClient) Do(request *http.Request) (*http.Response, error) {
	authRequest, err := h.Client.CreateRequest(request.Method, request.URL.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", authRequest.Header.Get("Authorization"))

	response, err := h.Client.HTTP.Do(request)
	if err != nil {
		return nil, err
	}

	response.Header.Add("Content-Type", "application/multipart")
	return response, nil
}

func (h httpClient) Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(req)
}

func (h httpClient) Head(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", endpoint, nil)
	if err != nil {
		return nil, err
	}
	return h.Do(req)
}

var _ ranger.HTTPClient = httpClient{}

func FromPivnet(token, slug, version string) (Payload, error) {
	var (
		payload Payload
	)

	config := pivnet.ClientConfig{
		Host:              pivnet.DefaultHost,
		UserAgent:         "tile-builder",
		SkipSSLValidation: false,
	}

	client := pivnet.NewClient(
		pivnet.NewAccessTokenOrLegacyToken(
			token,
			pivnet.DefaultHost,
			false,
		),
		config,
		logshim.NewLogShim(
			log.New(os.Stderr, "", log.LstdFlags),
			log.New(os.Stderr, "", log.LstdFlags),
			false,
		),
	)

	releases, err := client.Releases.List(slug)
	if err != nil {
		return payload, fmt.Errorf("could not get releases for product with %s: %s", slug, err)
	}

	for _, release := range releases {
		if release.Version == version {
			err := client.EULA.Accept(slug, release.ID)
			if err != nil {
				return payload, fmt.Errorf("could not match EULA for release %d: %s", release.ID, err)
			}

			productFiles, err := client.ProductFiles.ListForRelease(slug, release.ID)
			if err != nil {
				return payload, fmt.Errorf("could not get productFiles for release %d: %s", release.ID, err)
			}

			for _, productFile := range productFiles {
				matched, err := filepath.Match("*.pivotal", filepath.Base(productFile.AWSObjectKey))
				if err != nil {
					return payload, fmt.Errorf("could not match productFile %s: %s", productFile.AWSObjectKey, err)
				}
				if matched {
					link, err := productFile.DownloadLink()
					if err != nil {
						return payload, fmt.Errorf("could not get download link for productFile: %s", err)
					}

					parsedURL, _ := url.Parse(link)

					httpClient := &ranger.HTTPRanger{
						URL: parsedURL,
						Client: httpClient{
							Client: client,
						},
					}

					reader, err := ranger.NewReader(httpClient)
					if err != nil {
						return payload, fmt.Errorf("can not create a range client: %s", err)
					}

					length, err := reader.Length()
					if err != nil {
						return payload, fmt.Errorf("can not find length of productFile: %s", err)
					}

					zipReader, err := zip.NewReader(reader, length)
					if err != nil {
						return payload, fmt.Errorf("can not create a zip client: %s", err)
					}

					for _, zipFile := range zipReader.File {
						if metadataFile.MatchString(zipFile.Name) {
							reader, err := zipFile.Open()
							if err != nil {
								return payload, fmt.Errorf("can not open zip file %s: %s", zipFile.Name, err)
							}
							contents, err := ioutil.ReadAll(reader)
							if err != nil {
								return payload, fmt.Errorf("can not read zip file %s: %s", zipFile.Name, err)
							}

							err = yaml.UnmarshalStrict(contents, &payload)
							if err != nil {
								return payload, fmt.Errorf("could not unmarshal %s: %s", zipFile.Name, err)
							}

							return payload, nil
						}
					}
				}
			}
		}
	}

	return payload, fmt.Errorf("could not find release with version %s for %s", version, slug)
}
