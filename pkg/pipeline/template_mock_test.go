package pipeline

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type templateMock struct {
	mock.Mock
}

func (m *templateMock) Reader() io.ReadCloser {
	args := m.Called()
	return args.Get(0).(io.ReadCloser)
}

func (m *templateMock) Path() string {
	args := m.Called()
	return args.String(0)
}
