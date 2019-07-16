package buildpack_python

import (
	"github.com/mgranderath/SPaaS/config"
	"html/template"
	"os"
	"path/filepath"
)

const dockerfileTemplate = `FROM gliderlabs/alpine:3.9
WORKDIR /usr/src/app

RUN apk add --no-cache python3 && \
    if [ ! -e /usr/bin/python ]; then ln -sf python3 /usr/bin/python ; fi

RUN python3 -m ensurepip && \
    pip3 install --no-cache --upgrade pip3 setuptools wheel && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi

EXPOSE 5000:5000
COPY . . 
RUN pip3 install --no-cache-dir -r requirements.txt
CMD [{{range $index, $cmd := .Command}}"{{.}}"{{if (ne ($index) ($.Length))}},{{end}}{{end}}]
`

type options struct {
	Command []string
	Length  int
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
