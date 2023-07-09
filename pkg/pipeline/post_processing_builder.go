package pipeline

type PostProcessingBuilder interface {
	AddResultProcessor(PostProcessor) PostProcessingBuilder
	Build() (Pipeline, error)
}
