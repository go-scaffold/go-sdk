package pipeline

import (
	"text/template"
)

type FunctionsBuilder interface {
	WithFunctions(functions template.FuncMap) TemplateProviderBuilder
}
