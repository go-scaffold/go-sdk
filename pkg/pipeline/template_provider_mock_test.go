package pipeline

import (
	"github.com/stretchr/testify/mock"
)

type templateProviderMock struct {
	mock.Mock
}

func (m *templateProviderMock) NextTemplate() (*Template, error) {
	args := m.Called()
	t := args.Get(0)
	err := args.Error(1)
	if t == nil {
		return nil, err
	}
	return t.(*Template), err
}
