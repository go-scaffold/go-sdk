package pipeline

type postProcessingStep struct {
	PostProcessor

	nextStep  *postProcessingStep
	processor PostProcessor
}

func (p *postProcessingStep) Process(args *ProcessData) (*ProcessData, error) {
	out, err := p.processor.Process(args)
	defer args.Reader.Close()
	if err != nil {
		return nil, err
	}

	if p.nextStep != nil {
		return p.nextStep.Process(out)
	}
	return out, err
}
