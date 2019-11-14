package commands

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"

	"github.com/jtarchie/tile-builder/render"
)

type Preview struct {
	Port   int      `long:"port" default:"8181" description:"port number to listen on"`
	Tile   TileArgs `group:"tile" namespace:"tile" env-namespace:"TILE"`
	Strict bool     `long:"strict" description:"use strict unmarshaling for the tile"`
	Pivnet pivnet   `group:"pivnet" namespace:"pivnet" env-namespace:"PIVNET"`
}

func (p Preview) Execute(_ []string) error {
	payload, err := loadMetadataForTile(p.Tile, p.Pivnet, p.Strict)
	if err != nil {
		return err
	}

	_, err = render.AsHTML(payload)
	if err != nil {
		return err
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		contents, _ := render.AsHTML(payload)
		return c.HTMLBlob(http.StatusOK, contents)
	})

	fmt.Printf("listening on http://localhost:%d\n", p.Port)

	return e.Start(fmt.Sprintf(":%d", p.Port))
}
