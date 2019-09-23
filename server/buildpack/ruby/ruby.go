package buildpack_ruby

import (
	"html/template"
	"os"
	"path/filepath"
)

const dockerfileTemplate = `FROM andrius/alpine-ruby:3.4

WORKDIR /usr/src/app

RUN apk add --no-cache --virtual .build-deps ruby-dev build-base \
  libxml2-dev libxslt-dev pcre-dev libffi-dev \
  mariadb-dev postgresql-dev

COPY Gemfile Gemfile.lock ./
RUN bundle install

RUN apk del .build-deps

EXPOSE 5000:5000
COPY . . 
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

type options struct {
	Command []string
	Length  int
}

func Build(appPath string, command []string) error {
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
		Command: command,
		Length:  len(command) - 1,
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
