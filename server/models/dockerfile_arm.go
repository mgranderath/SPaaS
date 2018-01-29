// +build linux,arm

package models

import (
	"html/template"
	"os"
	"path/filepath"
)

// Dockerfile stores the template value
type Dockerfile struct {
	BuildName string
	Command   []string
	Length    int
	Port      string
	Type      string
}

const dockerfileTemplate = `FROM arm32v6/alpine:3.5

WORKDIR /usr/src/app

{{if eq .Type "python"}}
RUN apk add --no-cache python3 && \
    python3 -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip3 install --upgrade pip setuptools && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
    if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
    rm -r /root/.cache
COPY requirements.txt .
RUN pip3 install --no-cache-dir -r requirements.txt
{{end}}
{{if eq .Type "nodejs"}}
RUN apk update && apk upgrade && apk add nodejs
COPY package*.json ./
RUN npm install
{{end}}
{{if eq .Type "ruby"}}
RUN apk update && \
    apk upgrade && \
    apk add ruby
COPY Gemfile /usr/src/app/Gemfile
COPY Gemfile.lock /usr/src/app/Gemfile.lock 
RUN bundle install
{{end}}

EXPOSE 5000:5000

COPY . .

CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

// CreateDockerfile creates dockerfile
func CreateDockerfile(dock Dockerfile, app Application) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(dockerfileTemplate)
	if err != nil {
		return err
	}
	dock.Type = app.Type
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
