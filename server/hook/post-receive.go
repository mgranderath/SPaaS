package hook

import (
	"bytes"
	"html/template"
)

type templateData struct {
	Name           string
	Token          string
	CustomEndpoint string
	HTTPS          string
}

const postReceiveHook = `#!/usr/bin/env python
#!/usr/bin/python
# -*- coding: utf-8 -*-
import json
import urllib2
import urllib
import socket, ssl
socket._fileobject.default_bufsize = 0

INFO_START = "\33[33m"
SUCCESS_START = "\33[32m"
ERROR_START = "\33[91m"
END = "\033[0m"

url = '{{ .HTTPS }}{{ .CustomEndpoint }}:80/api/app/{{ .Name }}/deploy'
headers = {'Authorization': 'Bearer {{ .Token }}'}
values = {}
data = urllib.urlencode(values)
req = urllib2.Request(url, data, headers=headers)
response = urllib2.urlopen(req, context=ssl.SSLContext(ssl.PROTOCOL_TLSv1))
for line in response:
    if line == '\n':
        pass
    obj = json.loads(line)
    output = ''
    if obj['type'] == 'info':
        output += INFO_START + 'INFO:'.ljust(10)
        output += obj['message'] + END + '\r'
    elif obj['type'] == 'success':
        output += SUCCESS_START + 'SUCCESS:'.ljust(10)
        output += obj['message'] + END
    else:
        output += ERROR_START + 'ERROR:'.ljust(10)
        output += obj['message'] + END
    print output
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
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
