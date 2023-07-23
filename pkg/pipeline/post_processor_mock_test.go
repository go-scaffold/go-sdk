package pipeline

import (
	"github.com/stretchr/testify/mock"
)

type postProcessorMock struct {
	mock.Mock
}

func (m *postProcessorMock) Process(args *Template) error {
	res := m.Called(args)
	err := res.Error(0)
	return err
}
