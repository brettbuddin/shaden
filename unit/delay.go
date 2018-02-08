package unit

import (
	"buddin.us/shaden/dsp"
)

const maxDelayMS = 10000

var maxDelayValue = dsp.Duration(maxDelayMS).Float64()

func newDelay(io *IO, _ Config) (*Unit, error) {
	d := &delay{
		dl:       dsp.NewDelayLine(int(maxDelayValue)),
		in:       io.NewIn("in", dsp.Float64(0)),
		time:     io.NewIn("time", dsp.Duration(500)),
		mix:      io.NewIn("mix", dsp.Float64(0)),
		fbreturn: io.NewIn("fb-return", dsp.Float64(0)),
		fbgain:   io.NewIn("fb-gain", dsp.Float64(0)),
		out:      io.NewOut("out"),
		fbsend:   io.NewOut("fb-send"),
		block:    &dsp.DCBlock{},
	}
	return NewUnit(io, d), nil
}

type delay struct {
	in, time, mix, fbreturn, fbgain *In
	out, fbsend                     *Out
	dl                              *dsp.DelayLine
	last                            float64
	block                           *dsp.DCBlock
}

func (d *delay) ProcessSample(i int) {
	var (
		in     = d.in.Read(i)
		mix    = d.mix.Read(i)
		fbgain = d.fbgain.Read(i)
		time   = dsp.Clamp(d.time.Read(i), 0, maxDelayValue)
	)

	wet := d.dl.TickAbsolute(in+d.last*fbgain, time)

	if d.fbsend.DestinationCount() > 0 {
		d.fbsend.Write(i, wet)
		d.last = d.fbreturn.Read(i)
	} else {
		d.last = wet
	}
	d.out.Write(i, d.block.Tick(dsp.Mix(mix, in, d.last)))
}
