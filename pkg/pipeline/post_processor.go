package pipeline

type PostProcessor interface {
	Process(args *Template) error
}
