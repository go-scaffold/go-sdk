package pipeline

type TemplateProcessor interface {
	NextTemplate() (Template, error)
}
