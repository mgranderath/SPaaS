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
	Type      string
}

const dockerfileTemplate = `FROM {{.BuildName}}

WORKDIR /usr/src/app

{{if eq .Type "python"}}
COPY requirements.txt ./
RUN pip3 install --no-cache-dir -r requirements.txt
{{end}}
{{if eq .Type "nodejs"}}
COPY package*.json ./
RUN npm install
{{end}}
{{if eq .Type "ruby"}}
COPY Gemfile /usr/app/ 
COPY Gemfile.lock /usr/app/ 
RUN bundle install
{{end}}

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
	dock.Type = app.Type
	if app.Type == "python" {
		build := Buildpack{}
		if err := db.Read("buildpack", "python3", &build); err != nil {
			printErr(os.Stdout, err)
			return err
		}
		dock.BuildName = build.Name
	} else if app.Type == "nodejs" {
		build := Buildpack{}
		if err := db.Read("buildpack", "nodejs", &build); err != nil {
			printErr(os.Stdout, err)
			return err
		}
		dock.BuildName = build.Name
	} else if app.Type == "ruby" {
		build := Buildpack{}
		if err := db.Read("buildpack", "ruby", &build); err != nil {
			printErr(os.Stdout, err)
			return err
		}
		dock.BuildName = build.Name
	} else {
		return nil
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
