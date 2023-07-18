package pipeline

import "io"

type Template struct {
	Reader io.ReadCloser
	Path   string
}
