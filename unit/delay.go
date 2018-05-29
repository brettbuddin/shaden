package unit

import (
	"github.com/brettbuddin/shaden/dsp"
)

const maxDelayMS = 10000

func newDelay(io *IO, c Config) (*Unit, error) {
	maxDelay := dsp.Duration(maxDelayMS, c.SampleRate).Float64()

	return NewUnit(io, &delay{
		dl:       dsp.NewDelayLine(int(maxDelay)),
		in:       io.NewIn("in", dsp.Float64(0)),
		time:     io.NewIn("time", dsp.Duration(500, c.SampleRate)),
		mix:      io.NewIn("mix", dsp.Float64(0)),
		fbreturn: io.NewIn("fb-return", dsp.Float64(0)),
		fbgain:   io.NewIn("fb-gain", dsp.Float64(0)),
		out:      io.NewOut("out"),
		fbsend:   io.NewOut("fb-send"),
		block:    &dsp.DCBlock{},
		maxDelay: maxDelay,
	}), nil
}

type delay struct {
	in, time, mix, fbreturn, fbgain *In
	out, fbsend                     *Out
	dl                              *dsp.DelayLine
	maxDelay, last                  float64
	block                           *dsp.DCBlock
}

func (d *delay) ProcessSample(i int) {
	var (
		in     = d.in.Read(i)
		mix    = d.mix.Read(i)
		fbgain = d.fbgain.Read(i)
		time   = dsp.Clamp(d.time.Read(i), 0, d.maxDelay)
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
