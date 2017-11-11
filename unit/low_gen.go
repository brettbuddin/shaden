package unit

import (
	"math"
	"math/rand"

	"buddin.us/lumen/dsp"
)

func newLowGen(name string, _ Config) (*Unit, error) {
	io := NewIO()
	g := &lowGen{
		freq:   io.NewIn("freq", dsp.Frequency(1)),
		amp:    io.NewIn("amp", dsp.Float64(1)),
		pw:     io.NewIn("pulse-width", dsp.Float64(1)),
		offset: io.NewIn("offset", dsp.Float64(0)),
		sync:   io.NewIn("sync", dsp.Float64(-1)),
	}

	io.ExposeOutProcessor(g.newSine())
	io.ExposeOutProcessor(g.newTriangle())
	io.ExposeOutProcessor(g.newPulse())
	io.ExposeOutProcessor(g.newSaw())

	u := NewUnit(io, name, nil)
	u.rate = RateControl
	return u, nil
}

type lowGen struct {
	freq, amp, pw, offset, sync *In
	algorithm                   string
	phases                      []float64
}

func (g *lowGen) nextPhase() *float64 {
	g.phases = append(g.phases, rand.Float64())
	return &(g.phases[len(g.phases)-1])
}

func (g *lowGen) newSine() *lowGenSine {
	return &lowGenSine{
		lowGen: g,
		phase:  g.nextPhase(),
		out:    NewOut("sine", newFrame()),
	}
}

func (g *lowGen) newSaw() *lowGenSaw {
	return &lowGenSaw{
		lowGen: g,
		phase:  g.nextPhase(),
		out:    NewOut("saw", newFrame()),
	}
}

func (g *lowGen) newPulse() *lowGenPulse {
	return &lowGenPulse{
		lowGen: g,
		phase:  g.nextPhase(),
		out:    NewOut("pulse", newFrame()),
	}
}

func (g *lowGen) newTriangle() *lowGenTriangle {
	return &lowGenTriangle{
		lowGen: g,
		phase:  g.nextPhase(),
		out:    NewOut("triangle", newFrame()),
	}
}

type lowGenSine struct {
	*lowGen
	phase    *float64
	out      *Out
	lastSync float64
}

func (o *lowGenSine) ExternalNeighborCount() int { return o.out.ExternalNeighborCount() }
func (o *lowGenSine) Out() *Out                  { return o.out }
func (o *lowGenSine) ProcessFrame(int)           { o.ProcessSample(0) }

func (o *lowGenSine) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
	)
	if o.lastSync < 0 && sync > 0 {
		*o.phase = 0
	}
	next := dsp.Sin(*o.phase) * amp
	*o.phase = stepPhase(freq, 0, *o.phase, 1)
	o.out.Write(i, offset+next)
}

type lowGenSaw struct {
	*lowGen
	phase    *float64
	out      *Out
	lastSync float64
}

func (o *lowGenSaw) ExternalNeighborCount() int { return o.out.ExternalNeighborCount() }
func (o *lowGenSaw) Out() *Out                  { return o.out }
func (o *lowGenSaw) ProcessFrame(int)           { o.ProcessSample(0) }

func (o *lowGenSaw) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
	)
	if o.lastSync < 0 && sync > 0 {
		*o.phase = 0
	}
	next := (2.0*(*o.phase)/twoPi - 1.0) * amp
	*o.phase = stepPhase(freq, 0, *o.phase, 1)
	o.out.Write(0, offset+next)
}

type lowGenPulse struct {
	*lowGen
	phase    *float64
	out      *Out
	lastSync float64
}

func (o *lowGenPulse) ExternalNeighborCount() int { return o.out.ExternalNeighborCount() }
func (o *lowGenPulse) Out() *Out                  { return o.out }
func (o *lowGenPulse) ProcessFrame(int)           { o.ProcessSample(0) }

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
		*o.phase = 0
	}
	if *o.phase < math.Pi*pw {
		next = 1 * amp
	} else {
		next = -1 * amp
	}

	*o.phase = stepPhase(freq, 0, *o.phase, 1)
	o.out.Write(i, offset+next)
}

type lowGenTriangle struct {
	*lowGen
	phase          *float64
	out            *Out
	last, lastSync float64
}

func (o *lowGenTriangle) ExternalNeighborCount() int { return o.out.ExternalNeighborCount() }
func (o *lowGenTriangle) Out() *Out                  { return o.out }
func (o *lowGenTriangle) ProcessFrame(int)           { o.ProcessSample(0) }

func (o *lowGenTriangle) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		offset = o.offset.Read(i)
		p      = *o.phase
		sync   = o.sync.Read(i)
		next   float64
	)
	if o.lastSync < 0 && sync > 0 {
		*o.phase = 0
	}
	if p < math.Pi {
		next = float64(-1*amp+twoDivPi*p) * amp
	} else {
		next = float64(3*amp-twoDivPi*p) * amp
	}
	*o.phase = stepPhase(freq, 0, *o.phase, 1)
	o.out.Write(i, offset+next)
}
