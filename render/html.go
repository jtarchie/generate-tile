package render

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/jtarchie/tile-builder/metadata"
	"html/template"
	"log"
)

func AsHTML(payload metadata.Payload) ([]byte, error) {
	box := packr.New("box", "./templates")
	htmlTemplate, err := box.FindString("form.gohtml")
	if err != nil {
		return nil, fmt.Errorf("could not load template: %s", err)
	}
	t, err := template.New("preview").Funcs(template.FuncMap{
		"getPropertyBlueprint": func(pi metadata.PropertyInput) metadata.PropertyBlueprint {
			pb, _ := payload.FindPropertyBlueprintFromPropertyInput(pi)
			return pb
		},
		"log": func(message string) string {
			log.Print(message)
			return ""
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("could not render template: %s", err)
	}

	contents := &bytes.Buffer{}
	err = t.Execute(contents, payload)
	if err != nil {
		return nil, fmt.Errorf("could not execute template: %s", err)
	}

	return contents.Bytes(), nil
}
