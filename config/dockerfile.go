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

// NodeJsDockerfile houses the information to generate the Dockerfile
type NodeJsDockerfile struct {
	Command        []string
	Length         int
	VersionDefined bool
	Version        string
}

const nodeJsTemplate = `FROM gliderlabs/alpine:3.4
WORKDIR /usr/src/app
{{ if eq .VersionDefined true }}
RUN apk add --no-cache nodejs={{.Version}}
{{ else }}
RUN apk add --no-cache nodejs
{{ end }}
COPY package*.json ./
RUN npm install
EXPOSE 5000:5000
COPY . .
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

const rubyTemplate = `FROM gliderlabs/alpine:3.4

ENV NOKOGIRI_USE_SYSTEM_LIBRARIES=1

RUN apk update \
\
&& apk add ruby \
           ruby-bigdecimal \
           ruby-bundler \
           ruby-io-console \
           ruby-irb \
           build-base \
            ruby-dev \
 && apk add --update-cache --repository http://dl-4.alpinelinux.org/alpine/edge/main/ \
            ca-certificates \
            libressl \
            libressl-dev \
 \
 && bundle config build.nokogiri --use-system-libraries \
 && bundle config git.allow_insecure true \
 && gem install --no-rdoc --no-ri \
                json \
                foreman \
 \
 && gem cleanup \
 && apk del build-base \
            ruby-dev \
            libressl-dev \
 && rm -rf /usr/lib/ruby/gems/*/cache/* \
           /var/cache/apk/* \
           /tmp/*
 
RUN apk --no-cache add --virtual build-dependencies ruby-dev build-base \
  libxml2-dev libxslt-dev pcre-dev libffi-dev \
  mariadb-dev postgresql-dev

WORKDIR /usr/src/app

COPY Gemfile Gemfile.lock ./
RUN bundle install

RUN apk del build-dependencies

EXPOSE 5000:5000
COPY . . 
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

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

// CreateNodeJsDockerfile creates a nodejs dockerfile
func CreateNodeJsDockerfile(dock NodeJsDockerfile, appPath string) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(nodeJsTemplate)
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

// CreateRubyDockerfile creates a ruby dockerfile
func CreateRubyDockerfile(dock NodeJsDockerfile, appPath string) error {
	t := template.New("Dockerfile template")
	t, err := t.Parse(rubyTemplate)
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
