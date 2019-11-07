package commands

import (
	"fmt"
	"net/http"

	"github.com/jtarchie/tile-builder/render"

	"github.com/jtarchie/tile-builder/metadata"
)

type Preview struct {
	Path   string `long:"path" required:"true" description:"path to the pivotal file"`
	Port   int    `long:"port" default:"8181" description:"port number to listen on"`
	Server *http.Server
}

func (p Preview) Execute(_ []string) error {
	payload, err := metadata.FromTile(p.Path)
	if err != nil {
		return fmt.Errorf("could not load metadata from tile: %s", err)
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
