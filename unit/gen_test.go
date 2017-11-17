package unit

type genOutput interface {
	CondProcessor
	FrameProcessor
	SampleProcessor
	Output
}

// Ensure gen and low-gen outputs conform to the interfaces necessary for processing.
var _ = []genOutput{
	&genSine{},
	&genSaw{},
	&genPulse{},
	&genTriangle{},
	&genCluster{},
	&genNoise{},

	&lowGenSine{},
	&lowGenSaw{},
	&lowGenPulse{},
	&lowGenTriangle{},
}
