package commands

import (
	"fmt"
	"net/http"

	"github.com/jtarchie/tile-builder/render"
)

type Preview struct {
	Port   int      `long:"port" default:"8181" description:"port number to listen on"`
	Tile   TileArgs `group:"tile" namespace:"tile" env-namespace:"TILE"`
	Pivnet pivnet   `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
}

func (p Preview) Execute(_ []string) error {
	payload, err := loadMetadataForTile(p.Tile, p.Pivnet)
	if err != nil {
		return err
	}

	contents, err := render.AsHTML(payload)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		_, _ = response.Write(contents)
	})

	fmt.Printf("listening on http://localhost:%d\n", p.Port)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", p.Port),
	}
	return server.ListenAndServe()
}
