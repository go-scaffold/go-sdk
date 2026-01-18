package pipeline

import (
	"errors"
	"text/template"

	"github.com/go-scaffold/go-sdk/v2/pkg/templates"
)

type PipelineBuilder interface {
	Build() (Pipeline, error)
	WithCollector(p Collector) *pipelineBuilder
	WithDataPreprocessor(fn DataPreprocessor) *pipelineBuilder
	WithFunctions(functions template.FuncMap) *pipelineBuilder
	WithSharedTemplatesProvider(p TemplateProvider) *pipelineBuilder
	WithTemplateAwareFunctions(functions templates.TemplateAwareFuncMap) *pipelineBuilder
	WithTemplateProvider(p TemplateProvider) *pipelineBuilder
}

type pipelineBuilder struct {
	p *pipeline
}

func NewPipelineBuilder() PipelineBuilder {
	return &pipelineBuilder{
		p: &pipeline{},
	}
}

func (b *pipelineBuilder) Build() (Pipeline, error) {
	if len(b.p.functions) == 0 {
		return nil, errors.New("no functions specified in the context")
	}
	if b.p.templateProvider == nil {
		return nil, errors.New("no template processor specified for the pipeline")
	}
	if b.p.collector == nil {
		return nil, errors.New("no collector specified for the pipeline")
	}

	return b.p, nil
}

func (b *pipelineBuilder) WithCollector(p Collector) *pipelineBuilder {
	b.p.collector = p
	return b
}

func (b *pipelineBuilder) WithDataPreprocessor(fn DataPreprocessor) *pipelineBuilder {
	b.p.dataPreprocessor = fn
	return b
}

func (b *pipelineBuilder) WithFunctions(functions template.FuncMap) *pipelineBuilder {
	b.p.functions = functions
	return b
}

func (b *pipelineBuilder) WithSharedTemplatesProvider(p TemplateProvider) *pipelineBuilder {
	b.p.sharedTemplatesProvider = p
	return b
}

func (b *pipelineBuilder) WithTemplateAwareFunctions(functions templates.TemplateAwareFuncMap) *pipelineBuilder {
	b.p.templateAwareFns = functions
	return b
}

func (b *pipelineBuilder) WithTemplateProvider(p TemplateProvider) *pipelineBuilder {
	b.p.templateProvider = p
	return b
}
