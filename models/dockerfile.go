package models

import (
	"html/template"
	"os"
	"path/filepath"
)

// Dockerfile : stores the template value
type Dockerfile struct {
	BuildName string
	Command   []string
	Length    int
	Port      string
}

const dockerfileTemplate = `FROM {{.BuildName}}

WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip3 install --no-cache-dir -r requirements.txt

EXPOSE 5000:5000

COPY . .

CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

// CreateDockerfile : create dockerfile
func CreateDockerfile(dock Dockerfile, app Application) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	if app.Type == "python" {
		build := Buildpack{}
		if err := db.Read("buildpack", "python3", &build); err != nil {
			printErr(os.Stdout, err)
			return err
		}
		dock.BuildName = build.Name
	} else {
		return err
	}
	dock.Length = len(dock.Command) - 1
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
