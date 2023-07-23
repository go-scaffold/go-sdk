package pipeline

type CollectorBuilder interface {
	WithCollector(Collector) PipelineBuilder
}
