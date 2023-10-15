package pipeline

import (
	"errors"
	"text/template"
)

type PipelineBuilder interface {
	Build() (Pipeline, error)
	WithCollector(p Collector) *pipelineBuilder
	WithDataPrefix(prefix string) *pipelineBuilder
	WithDataPreprocessor(fn DataPreprocessor) *pipelineBuilder
	WithFunctions(functions template.FuncMap) *pipelineBuilder
	WithMetadataPrefix(prefix string) *pipelineBuilder
	WithTemplateProvider(p TemplateProvider) *pipelineBuilder
}

type pipelineBuilder struct {
	p *pipeline
}

func NewPipelineBuilder() PipelineBuilder {
	return &pipelineBuilder{
		p: &pipeline{
			prefixData:     "Values",
			prefixMetadata: "Manifest",
		},
	}
}

func (b *pipelineBuilder) Build() (Pipeline, error) {
	if b.p.functions == nil || len(b.p.functions) == 0 {
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

func (b *pipelineBuilder) WithDataPrefix(prefix string) *pipelineBuilder {
	b.p.prefixData = prefix
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

func (b *pipelineBuilder) WithMetadataPrefix(prefix string) *pipelineBuilder {
	b.p.prefixMetadata = prefix
	return b
}

func (b *pipelineBuilder) WithTemplateProvider(p TemplateProvider) *pipelineBuilder {
	b.p.templateProvider = p
	return b
}
