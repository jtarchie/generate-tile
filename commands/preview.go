package commands

import (
	"bytes"
	"fmt"
	"github.com/jtarchie/generate-tile/metadata"
	"html/template"
	"net/http"
)

type Preview struct {
	Path string `long:"path" required:"true" description:"path to the pivotal file"`
	Port int    `long:"port" default:"8181" description:"port number to listen on"`
}

var htmlTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
  </head>
<body>
<div class="container-fluid">
  <div class="row">
    <h1>{{.Label}} @ v{{.ProductVersion}}</h1>
  </div>
  <div class="row">
    <div class="col-4">
      <nav class="nav nav-pills flex-column" id="form-types" role="tablist">
        {{ range $index, $ft := .FormTypes }}
          <a class="nav-link {{if eq $index 0}}active{{end}}" href="#{{$ft.Name}}" id="{{$ft.Name}}-tab" data-toggle="tab" href="#{{$ft.Name}}" role="tab" aria-controls="{{$ft.Name}}" aria-selected="false">{{$ft.Label}}</a>
        {{end}}
      </nav>
    </div>
    <div class="col-8">
      <div class="tab-content">
        {{range $index, $ft := .FormTypes}}
          <div class="tab-pane {{if eq $index 0}}active{{end}}" id="{{$ft.Name}}" role="tabpanel" aria-labelledby="{{$ft.Name}}-tab">
            <p>{{$ft.Description}}</p>
            <form id="form-{{$ft.Name}}">
              {{range .PropertyInputs}}
{{ $p := getProperty . }} 
<div class="form-group">
  {{ if eq $p.Type "string"}}
    <label for="{{$p.Name}}">{{.Label}}</label>
    <input type="text" class="form-control" id="{{$p.Name}}" {{if eq $p.Optional false}}required{{end}}>
    <small class="form-text text-muted">{{.Description}}</small>
  {{end}}
</div>
              {{end}}
              <button type="submit" class="btn btn-primary">Submit</button>
            </form>
          </div>
        {{end}}
      </div>
    </div>
  </div>
</div>
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
</body>
</html>
`

func (p Preview) Execute(_ []string) error {
	payload, err := metadata.FromTile(p.Path)
	if err != nil {
		return fmt.Errorf("could not load metadata from tile: %s", err)
	}

	t, err := template.New("preview").Funcs(template.FuncMap{
		"getProperty": func(pi metadata.PropertyInput) metadata.PropertyBlueprint {
			for _, pb := range payload.PropertyBlueprints {
				if fmt.Sprintf(".properties.%s",pb.Name) == pi.Reference {
					return pb
				}
			}
			return metadata.PropertyBlueprint{}
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("could not render template: %s", err)
	}

	contents := &bytes.Buffer{}
	err = t.Execute(contents, payload)
	if err != nil {
		return fmt.Errorf("could not execute template: %s", err)
	}

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		_ = t.Execute(response, payload)
	})

	fmt.Printf("listening on http://localhost:%d\n", p.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d",p.Port), nil)
}
