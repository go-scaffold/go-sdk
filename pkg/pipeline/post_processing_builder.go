package pipeline

type CollectorBuilder interface {
	WithCollector(Collector) CollectorBuilder
	Build() (Pipeline, error)
}
