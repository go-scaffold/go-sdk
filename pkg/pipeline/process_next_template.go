package pipeline

import (
	"io"
	"log/slog"
	"text/template"

	"github.com/go-scaffold/go-sdk/v2/pkg/templates"
)

var _processTemplate = templates.ProcessTemplateWithBaseTemplate

func processNextTemplate(templateProvider TemplateProvider, data interface{}, funcMap template.FuncMap, templateAwareFnGen templates.TemplateAwareFuncMap, baseTemplate *template.Template) (*Template, error) {
	template, err := templateProvider.NextTemplate()
	if err != nil {
		return nil, err
	}

	slog.Info("Processing template file", slog.String("path", template.Path))

	templateReader := template.Reader
	defer templateReader.Close()

	resultReader, err := _processTemplate(templateReader, data, funcMap, templateAwareFnGen, baseTemplate)
	if err != nil {
		return nil, err
	}

	return &Template{
		Path:   template.Path,
		Reader: io.NopCloser(resultReader),
	}, nil
}
