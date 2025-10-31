package collectors

import (
	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
)

type baseCollector struct {
	next pipeline.Collector
}
