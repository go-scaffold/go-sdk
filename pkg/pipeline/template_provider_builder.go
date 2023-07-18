package pipeline

type TemplateProviderBuilder interface {
	WithTemplateProvider(p TemplateProvider) PostProcessingBuilder
}
