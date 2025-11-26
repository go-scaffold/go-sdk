package collectors

import (
	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/stretchr/testify/mock"
)

type mockCollector struct {
	mock.Mock
}

func (m *mockCollector) Collect(tpl *pipeline.Template) error {
	args := m.Called(tpl)
	return args.Error(0)
}

func (m *mockCollector) OnPipelineCompleted() error {
	args := m.Called()
	return args.Error(0)
}
