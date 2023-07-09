package pipeline

import "io"

type Template interface {
	Reader() io.ReadCloser
	Path() string
}
