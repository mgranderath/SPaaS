package buildpack_ruby

import (
	"html/template"
	"os"
	"path/filepath"
)

const dockerfileTemplate = `FROM gliderlabs/alpine:3.9

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
