package hook

import (
	"bytes"
	"fmt"
	"html/template"
)

type templateData struct {
	Name           string
	Token          string
	CustomEndpoint string
	HTTPS          string
}

const postReceiveHook = `#!/usr/bin/env bash
./hooks/post-receive-deploy {{ .Name }}
`

// CreatePostReceive returns the code for the post-receive hook
func CreatePostReceive(name string, token string, endpoint string, HTTPS string) (string, error) {
	t := template.New("Post Receive Hook")
	t, err := t.Parse(postReceiveHook)
	if err != nil {
		return "", err
	}
	data := templateData{}
	data.Name = name
	data.Token = token
	data.CustomEndpoint = endpoint
	data.HTTPS = HTTPS
	fmt.Println(data.Name)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
