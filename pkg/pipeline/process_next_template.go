package pipeline

import (
	"io/ioutil"
	"text/template"

	"github.com/go-scaffold/go-sdk/pkg/templates"
)

var _processTemplate = templates.ProcessTemplate

func processNextTemplate(templateProvider TemplateProvider, data interface{}, funcMap template.FuncMap) (*Template, error) {
	template, err := templateProvider.NextTemplate()
	if err != nil {
		return nil, err
	}

	templateReader := template.Reader
	defer templateReader.Close()

	resultReader, err := _processTemplate(templateReader, data, funcMap)
	if err != nil {
		return nil, err
	}

	return &Template{
		Path:   template.Path,
		Reader: ioutil.NopCloser(resultReader),
	}, nil
}
