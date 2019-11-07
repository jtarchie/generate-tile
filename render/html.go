package render

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/jtarchie/tile-builder/metadata"
)

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
                  {{ if eq $p.Type "integer"}}
                    <div class="form-group">
                      <label for="{{$p.Name}}">{{.Label}}</label>
                      <input type="number" class="form-control" id="{{$p.Name}}" {{if eq $p.Optional false}}required{{end}}>
                      <small class="form-text text-muted">{{.Description}}</small>
                    </div>
                  {{end}}
                  {{ if eq $p.Type "string"}}
                    <div class="form-group">
                      <label for="{{$p.Name}}">{{.Label}}</label>
                      <input type="text" class="form-control" id="{{$p.Name}}" {{if eq $p.Optional false}}required{{end}}>
                      <small class="form-text text-muted">{{.Description}}</small>
                    </div>
                  {{end}}
                  {{ if eq $p.Type "boolean"}}
                    <div class="form-check">
                    <label for="{{$p.Name}}">{{.Label}}</label>
                    <input type="checkbox" class="form-control" id="{{$p.Name}}" {{if eq $p.Optional false}}required{{end}}>
                    <small class="form-text text-muted">{{.Description}}</small>
                    </div>
                  {{end}}
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

func AsHTML(payload metadata.Payload) ([]byte, error) {
	t, err := template.New("preview").Funcs(template.FuncMap{
		"getProperty": func(pi metadata.PropertyInput) metadata.PropertyBlueprint {
			pb, _ := payload.FindPropertyBlueprintFromPropertyInput(pi)
			return pb
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
