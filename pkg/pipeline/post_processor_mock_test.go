package pipeline

import (
	"github.com/stretchr/testify/mock"
)

type postProcessorMock struct {
	mock.Mock
}

func (m *postProcessorMock) Process(args *Template) (*Template, error) {
	res := m.Called(args)
	t := res.Get(0)
	err := res.Error(1)
	if t == nil {
		return nil, err
	}
	return t.(*Template), err
}
