package main

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Templater struct {
	tmpl *template.Template
}

func NewTemplater(tmplStr string) (*Templater, error) {
	tmpl, err := template.New("webhook").
		Funcs(sprig.FuncMap()).
		Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	return &Templater{tmpl: tmpl}, nil
}

func (t *Templater) template(data any) (string, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
