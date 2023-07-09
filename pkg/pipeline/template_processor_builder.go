package pipeline

type TemplateProcessorBuilder interface {
	WithTemplateProcessor(p TemplateProcessor) PostProcessingBuilder
}
