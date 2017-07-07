package statedetector

type dummyXfsTracer struct {
}

func NewDummyXfsTracer() XfsTracer {
	return &binaryXfsTracer{}
}

func (b *dummyXfsTracer) Run() error {
	return nil
}
