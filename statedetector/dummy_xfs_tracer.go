package statedetector

type dummyXfsTracer struct {
}

func NewDummyXfsTracer() XfsTracer {
	return &binaryXfsTracer{}
}

func (b *dummyXfsTracer) Run() error {
	return nil
}

func (b *dummyXfsTracer) Start() error {
	return nil
}

func (b *dummyXfsTracer) Stop() error {
	return nil
}
