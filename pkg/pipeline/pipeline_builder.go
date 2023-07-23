package pipeline

type PipelineBuilder interface {
	Build() (Pipeline, error)
}
