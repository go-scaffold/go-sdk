package pipeline

import (
	"errors"
	"io"
	"log/slog"
	"text/template"

	"github.com/go-scaffold/go-sdk/v2/pkg/templates"
)

var _processNextTemplate = processNextTemplate

type Pipeline interface {
	Process(processData map[string]interface{}) error
}

type pipeline struct {
	dataPreprocessor       DataPreprocessor
	functions              template.FuncMap
	templateAwareFns       templates.TemplateAwareFuncMap
	collector              Collector
	templateProvider       TemplateProvider
	namedTemplatesProvider TemplateProvider
}

// loadCommonTemplates loads all common templates into a base template that can be
// reused across all main templates in the pipeline.
func (p *pipeline) loadCommonTemplates() (*template.Template, error) {
	if p.namedTemplatesProvider == nil {
		return nil, nil
	}

	baseTemplate := template.New("").Funcs(p.functions)

	for {
		commonTemplate, err := p.namedTemplatesProvider.NextTemplate()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		slog.Info("Loading common template", slog.String("name", commonTemplate.Name))

		content, err := io.ReadAll(commonTemplate.Reader)
		commonTemplate.Reader.Close()
		if err != nil {
			return nil, err
		}

		_, err = baseTemplate.New(commonTemplate.Name).Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	return baseTemplate, nil
}

func (p *pipeline) Process(processData map[string]interface{}) error {
	var err error

	if p.dataPreprocessor != nil {
		processData, err = p.dataPreprocessor(processData)
		if err != nil {
			return err
		}
	}

	// Load common templates once before processing main templates
	baseTemplate, err := p.loadCommonTemplates()
	if err != nil {
		return err
	}

	for err == nil {
		err = p.processNext(processData, baseTemplate)
	}
	if errors.Is(err, io.EOF) {
		return p.collector.OnPipelineCompleted()
	}
	return err
}

func (p *pipeline) processNext(data map[string]interface{}, baseTemplate *template.Template) error {
	result, err := _processNextTemplate(p.templateProvider, data, p.functions, p.templateAwareFns, baseTemplate)
	if err != nil {
		return err
	}

	return p.collector.Collect(result)
}
