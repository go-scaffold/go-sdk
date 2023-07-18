package pipeline

import (
	"errors"
	"text/template"

	"github.com/pasdam/go-template-map-loader/pkg/tm"
)

type builder struct {
	MetadataBuilder
	DataBuilder
	FunctionsBuilder
	TemplateProviderBuilder
	PostProcessingBuilder
	lastStep *postProcessingStep

	p *pipeline
}

func NewBuilder() MetadataBuilder {
	return &builder{
		p: &pipeline{},
	}
}

func (b *builder) WithData(data map[string]interface{}) FunctionsBuilder {
	b.WithDataWithPrefix("Values", data)
	return b
}

func (b *builder) WithDataWithPrefix(prefix string, data map[string]interface{}) FunctionsBuilder {
	b.p.data = tm.MergeMaps(b.p.data, tm.WithPrefix(prefix, data))
	return b
}

func (b *builder) WithMetadata(data map[string]interface{}) DataBuilder {
	b.WithDataWithPrefix("Manifest", data)
	return b
}

func (b *builder) WithMetadataWithPrefix(prefix string, data map[string]interface{}) DataBuilder {
	b.WithDataWithPrefix(prefix, data)
	return b
}

func (b *builder) WithFunctions(functions template.FuncMap) TemplateProviderBuilder {
	b.p.functions = functions
	return b
}

func (b *builder) WithTemplateProvider(p TemplateProvider) PostProcessingBuilder {
	b.p.templateProvider = p
	return b
}

func (b *builder) AddResultProcessor(p PostProcessor) PostProcessingBuilder {
	step := &postProcessingStep{
		processor: p,
	}
	if b.p.postProcessingSteps == nil {
		b.p.postProcessingSteps = step
	} else {
		b.lastStep.nextStep = step
	}
	b.lastStep = step
	return b
}

func (b *builder) Build() (Pipeline, error) {
	dataLen := 0
	for k := range b.p.data {
		dataLen += len(b.p.data[k].(map[string]interface{}))
	}
	if dataLen == 0 {
		return nil, errors.New("no data specified in the context")
	}
	if b.p.functions == nil || len(b.p.functions) == 0 {
		return nil, errors.New("no functions specified in the context")
	}
	if b.p.templateProvider == nil {
		return nil, errors.New("no template processor specified for the pipeline")
	}
	if b.p.postProcessingSteps == nil {
		return nil, errors.New("no post processor specified for the pipeline")
	}

	return b.p, nil
}
