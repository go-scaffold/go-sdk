package collectors

import (
	"github.com/go-scaffold/go-sdk/pkg/filters"
	"github.com/go-scaffold/go-sdk/pkg/pipeline"
)

type filterCollector struct {
	baseCollector

	filter filters.Filter
}

func NewFilterCollector(filter filters.Filter, nextCollector pipeline.Collector) pipeline.Collector {
	return &filterCollector{
		filter: filter,
		baseCollector: baseCollector{
			next: nextCollector,
		},
	}
}

func (p *filterCollector) Collect(args *pipeline.Template) error {
	if p.filter.Accept(args.Path) {
		return p.next.Collect(args)
	}
	return nil
}
