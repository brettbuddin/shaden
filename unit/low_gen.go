package unit

import (
	"math"
	"math/rand"

	"buddin.us/shaden/dsp"
)

func newLowGen(io *IO, _ Config) (*Unit, error) {
	g := &lowGen{
		freq:   io.NewIn("freq", dsp.Frequency(1)),
		amp:    io.NewIn("amp", dsp.Float64(1)),
		pw:     io.NewIn("pulse-width", dsp.Float64(1)),
		offset: io.NewIn("offset", dsp.Float64(0)),
		sync:   io.NewIn("sync", dsp.Float64(-1)),
	}

	io.ExposeOutputProcessor(g.newSine())
	io.ExposeOutputProcessor(g.newTriangle())
	io.ExposeOutputProcessor(g.newPulse())
	io.ExposeOutputProcessor(g.newSaw())

	u := NewUnit(io, nil)
	u.rate = RateControl
	return u, nil
}

type lowGen struct {
	freq, amp, pw, offset, sync *In
}

func (g *lowGen) newSine() *lowGenSine {
	return &lowGenSine{
		lowGen: g,
		phase:  rand.Float64() * twoPi,
		out:    NewOut("sine", newFrame()),
	}
}

func (g *lowGen) newSaw() *lowGenSaw {
	return &lowGenSaw{
		lowGen: g,
		phase:  rand.Float64() * twoPi,
		out:    NewOut("saw", newFrame()),
	}
}

func (g *lowGen) newPulse() *lowGenPulse {
	return &lowGenPulse{
		lowGen: g,
		phase:  rand.Float64() * twoPi,
		out:    NewOut("pulse", newFrame()),
	}
}

func (g *lowGen) newTriangle() *lowGenTriangle {
	return &lowGenTriangle{
		lowGen: g,
		phase:  rand.Float64() * twoPi,
		out:    NewOut("triangle", newFrame()),
	}
}

type lowGenSine struct {
	*lowGen
	phase, lastSync float64
	out             *Out
}

func (o *lowGenSine) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *lowGenSine) Out() *Out           { return o.out }
func (o *lowGenSine) ProcessFrame(int)    { o.ProcessSample(0) }

func (o *lowGenSine) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
	)
	if o.lastSync < 0 && sync > 0 {
		o.phase = 0
	}
	next := dsp.Sin(o.phase)
	o.phase = stepPhase(freq, 0, o.phase, 1)
	o.out.Write(i, (amp*next)+offset)
}

type lowGenSaw struct {
	*lowGen
	phase, lastSync float64
	out             *Out
}

func (o *lowGenSaw) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *lowGenSaw) Out() *Out           { return o.out }
func (o *lowGenSaw) ProcessFrame(int)    { o.ProcessSample(0) }

func (o *lowGenSaw) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
	)
	if o.lastSync < 0 && sync > 0 {
		o.phase = 0
	}
	next := (2.0*(o.phase)/twoPi - 1.0)
	o.phase = stepPhase(freq, 0, o.phase, 1)
	o.out.Write(i, (amp*next)+offset)
}

type lowGenPulse struct {
	*lowGen
	phase, lastSync float64
	out             *Out
}

func (o *lowGenPulse) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *lowGenPulse) Out() *Out           { return o.out }
func (o *lowGenPulse) ProcessFrame(int)    { o.ProcessSample(0) }

func (o *lowGenPulse) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		pw     = math.Abs(o.pw.Read(i))
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
		next   float64
	)

	if o.lastSync < 0 && sync > 0 {
		o.phase = 0
	}
	if o.phase < math.Pi*pw {
		next = 1
	} else {
		next = -1
	}

	o.phase = stepPhase(freq, 0, o.phase, 1)
	o.out.Write(i, (amp*next)+offset)
}

type lowGenTriangle struct {
	*lowGen
	phase, lastSync float64
	out             *Out
}

func (o *lowGenTriangle) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *lowGenTriangle) Out() *Out           { return o.out }
func (o *lowGenTriangle) ProcessFrame(int)    { o.ProcessSample(0) }

func (o *lowGenTriangle) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		p      = o.phase
		sync   = o.sync.Read(i)
		next   float64
	)
	if o.lastSync < 0 && sync > 0 {
		o.phase = 0
	}
	if p < math.Pi {
		next = (-1 + twoDivPi*p)
	} else {
		next = (3 - twoDivPi*p)
	}
	o.phase = stepPhase(freq, 0, o.phase, 1)
	o.out.Write(i, (amp*next)+offset)
}
