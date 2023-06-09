package pipeline

import (
	"errors"
	"io"
	"text/template"
)

var _processNextTemplate = processNextTemplate

type Pipeline interface {
	Process() error
}

type pipeline struct {
	data                map[string]interface{}
	functions           template.FuncMap
	postProcessingSteps *postProcessingStep
	templateProcessor   TemplateProcessor
}

func (p *pipeline) Process() error {
	var err error
	for err == nil {
		err = p.processNext()
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (p *pipeline) processNext() error {
	data, err := _processNextTemplate(p.templateProcessor, p.data, p.functions)
	if err != nil {
		return err
	}

	_, err = p.postProcessingSteps.Process(data)

	return err
}
