package pipeline

type PostProcessingBuilder interface {
	WithResultProcessor(PostProcessor) PostProcessingBuilder
	Build() (Pipeline, error)
}
