package hook

import (
	"bytes"
	"fmt"
	"html/template"
)

type templateData struct {
	Name string
}

const postReceiveHook = `#!/usr/bin/env bash
./hooks/post-receive-deploy {{ .Name }}
`

// GetPostReceiveHookHelperString returns the code for the post-receive hook
func GetPostReceiveHookHelperString(name string) (string, error) {
	t := template.New("Post Receive Hook")
	t, err := t.Parse(postReceiveHook)
	if err != nil {
		return "", err
	}
	data := templateData{}
	data.Name = name
	fmt.Println(data.Name)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
