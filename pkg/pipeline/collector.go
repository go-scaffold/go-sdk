package pipeline

type Collector interface {
	Collect(args *Template) error
}
