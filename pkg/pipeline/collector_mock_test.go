package pipeline

import (
	"github.com/stretchr/testify/mock"
)

type collectorMock struct {
	mock.Mock
}

func (m *collectorMock) Collect(args *Template) error {
	res := m.Called(args)
	err := res.Error(0)
	return err
}
