package pipeline

import (
	"github.com/stretchr/testify/mock"
)

type postProcessorMock struct {
	mock.Mock
}

func (m *postProcessorMock) Process(args *ProcessData) (*ProcessData, error) {
	res := m.Called(args)
	t := res.Get(0)
	err := res.Error(1)
	if t == nil {
		return nil, err
	}
	return t.(*ProcessData), err
}
