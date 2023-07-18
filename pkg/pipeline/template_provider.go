package pipeline

type TemplateProvider interface {
	NextTemplate() (Template, error)
}
