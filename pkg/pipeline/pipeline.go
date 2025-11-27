package pipeline

import (
	"errors"
	"io"
	"text/template"

	"github.com/go-scaffold/go-sdk/v2/pkg/templates"
)

var _processNextTemplate = processNextTemplate

type Pipeline interface {
	Process(processData map[string]interface{}) error
}

type pipeline struct {
	dataPreprocessor DataPreprocessor
	functions        template.FuncMap
	templateAwareFns templates.TemplateAwareFuncMap
	collector        Collector
	templateProvider TemplateProvider
}

func (p *pipeline) Process(processData map[string]interface{}) error {
	var err error

	if p.dataPreprocessor != nil {
		processData, err = p.dataPreprocessor(processData)
		if err != nil {
			return err
		}
	}

	for err == nil {
		err = p.processNext(processData)
	}
	if errors.Is(err, io.EOF) {
		return p.collector.OnPipelineCompleted()
	}
	return err
}

func (p *pipeline) processNext(data map[string]interface{}) error {
	result, err := _processNextTemplate(p.templateProvider, data, p.functions, p.templateAwareFns)
	if err != nil {
		return err
	}

	return p.collector.Collect(result)
}
