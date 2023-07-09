package pipeline

type PostProcessor interface {
	Process(args *ProcessData) (*ProcessData, error)
}
