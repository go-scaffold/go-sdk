package pipeline

import (
	"io/ioutil"
	"text/template"

	"github.com/go-scaffold/go-sdk/pkg/templates"
)

var _processTemplate = templates.ProcessTemplate

func processNextTemplate(templateProcessor TemplateProcessor, data interface{}, funcMap template.FuncMap) (*ProcessData, error) {
	template, err := templateProcessor.NextTemplate()
	if err != nil {
		return nil, err
	}

	templateReader := template.Reader()
	defer templateReader.Close()

	resultReader, err := _processTemplate(templateReader, data, funcMap)
	if err != nil {
		return nil, err
	}

	return &ProcessData{
		Path:   template.Path(),
		Reader: ioutil.NopCloser(resultReader),
	}, nil
}
