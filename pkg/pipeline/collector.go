package pipeline

type Collector interface {
	Collect(args *Template) error
	OnPipelineCompleted() error
}
