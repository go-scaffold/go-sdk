package pipeline

import "io"

type Template struct {
	// Name is the unique identifier for this template, used to reference it from other templates.
	// This is optional and primarily used for common templates.
	Name string

	// Path is the file path of the template.
	Path string

	// Reader provides access to the template content.
	Reader io.ReadCloser
}
