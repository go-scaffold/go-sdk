package collectors

import (
	"github.com/go-scaffold/go-sdk/pkg/pipeline"
)

type baseCollector struct {
	next pipeline.Collector
}
