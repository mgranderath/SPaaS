package buildpack_nodejs

import (
	"github.com/mgranderath/SPaaS/config"
	"html/template"
	"os"
	"path"
	"path/filepath"
)

const dockerfileTemplate = `FROM gliderlabs/alpine:3.9
WORKDIR /usr/src/app
RUN apk add --no-cache nodejs
{{ if .Yarn }}
RUN apk add --no-cache yarn
{{ end }}
COPY package*.json ./
{{ if .Yarn }}
RUN yarn
{{ else }}
RUN npm install
{{ end }}
EXPOSE 5000:5000
COPY . .
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

type options struct {
	Command []string
	Length  int
	Yarn    bool
}

func Build(appPath string, dockerfile config.Dockerfile) error {
	var (
		t   *template.Template
		err error
	)
	t = template.New("Dockerfile template")
	t, err = t.Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	options := options{
		Command: dockerfile.Command,
		Length:  len(dockerfile.Command) - 1,
	}
	if _, err := os.Stat(path.Join(appPath, "yarn.lock")); err == nil {
		options.Yarn = true
	} else {
		options.Yarn = false
	}
	f, err := os.Create(filepath.Join(appPath, "deploy", "Dockerfile"))
	if err != nil {
		return err
	}
	err = t.Execute(f, options)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
