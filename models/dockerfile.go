package models

import (
	"html/template"
	"os"
	"path/filepath"
)

// Dockerfile : stores the template value
type Dockerfile struct {
	Buildpack string
}

const dockerfileTemplate = `FROM {{.Buildpack}}

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip3 install --no-cache-dir -r requirements.txt

EXPOSE 5000:5000

COPY . .

CMD [ "python3", "app.py" ]
`

// CreateDockerfile : create dockerfile
func CreateDockerfile(app Application) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	dock := Dockerfile{}
	if app.Type == "python" {
		dock.Buildpack = "arm32v6/python:alpine3.6"
	} else {
		return err
	}
	f, err := os.Create(filepath.Join(app.Path, "deploy", "Dockerfile"))
	if err != nil {
		return err
	}
	err = t.Execute(f, dock)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
