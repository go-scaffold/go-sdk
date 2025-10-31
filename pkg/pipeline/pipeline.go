package pipeline

import (
	"errors"
	"io"
	"text/template"
)

var _processNextTemplate = processNextTemplate

type Pipeline interface {
	Process(processData map[string]interface{}) error
}

type pipeline struct {
	dataPreprocessor DataPreprocessor
	functions        template.FuncMap
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
		return nil
	}
	return err
}

func (p *pipeline) processNext(data map[string]interface{}) error {
	result, err := _processNextTemplate(p.templateProvider, data, p.functions)
	if err != nil {
		return err
	}

	err = p.collector.Collect(result)

	return err
}
