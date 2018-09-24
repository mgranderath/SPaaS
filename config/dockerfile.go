package config

import (
	"html/template"
	"os"
	"path/filepath"
)

// Dockerfile stores the template value
type Dockerfile struct {
	BuildName string
	Command   []string
	Port      string
	Type      string
	Length    int
}

const dockerfileTemplate = `FROM gliderlabs/alpine:3.4
WORKDIR /usr/src/app
{{if eq .Type "nodejs"}}
RUN apk add --no-cache nodejs
COPY package*.json ./
RUN npm install
{{end}}
EXPOSE 5000:5000
COPY . .
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

// CreateDockerfile creates dockerfile
func CreateDockerfile(dock Dockerfile, appPath string) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	dock.Length = len(dock.Command) - 1
	f, err := os.Create(filepath.Join(appPath, "deploy", "Dockerfile"))
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
