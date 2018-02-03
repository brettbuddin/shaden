package unit

type noopSampleProc struct{}

func (p noopSampleProc) ProcessSample(i int) {}

type frameProcessor struct {
	frame  func(int)
	sample func(int)
}

func (p frameProcessor) ProcessFrame(n int) {
	if p.frame != nil {
		p.frame(n)
	}
}

func (p frameProcessor) ProcessSample(i int) {
	if p.sample != nil {
		p.sample(i)
	}
}

type sampleProcessor struct {
	fn func(int)
}

func (p sampleProcessor) ProcessSample(i int) {
	if p.fn != nil {
		p.fn(i)
	}
}

type closer struct {
	fn func() error
}

func (c closer) Close() error {
	if c.fn != nil {
		return c.fn()
	}
	return nil
}

type outProcessor struct {
	fn  func(int)
	out *Out
}

func (p outProcessor) Out() *Out { return p.out }

func (p outProcessor) ProcessSample(i int) {
	if p.fn != nil {
		p.fn(i)
	}
}

func (p outProcessor) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		p.ProcessSample(i)
	}
}

type sampleProcessorCloser struct {
	sampleProcessor
	closer
}

type outProcessorCloser struct {
	closer
	outProcessor
}
