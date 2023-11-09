package pipeline

import (
	"io"
	"log/slog"
	"text/template"

	"github.com/go-scaffold/go-sdk/pkg/templates"
)

var _processTemplate = templates.ProcessTemplate

func processNextTemplate(templateProvider TemplateProvider, data interface{}, funcMap template.FuncMap) (*Template, error) {
	template, err := templateProvider.NextTemplate()
	if err != nil {
		return nil, err
	}

	slog.Info("Processing template file", slog.String("path", template.Path))

	templateReader := template.Reader
	defer templateReader.Close()

	resultReader, err := _processTemplate(templateReader, data, funcMap)
	if err != nil {
		return nil, err
	}

	return &Template{
		Path:   template.Path,
		Reader: io.NopCloser(resultReader),
	}, nil
}
