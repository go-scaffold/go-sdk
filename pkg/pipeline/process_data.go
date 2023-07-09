package pipeline

import (
	"io"
)

type ProcessData struct {
	Path   string
	Reader io.ReadCloser
}
