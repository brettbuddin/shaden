package unit

import (
	"math"
	"math/rand"

	"github.com/brettbuddin/shaden/dsp"
)

const (
	twoPi    = 2 * math.Pi
	twoDivPi = 2 / math.Pi
)

func newGen(io *IO, c Config) (*Unit, error) {
	g := &gen{
		rand:      c.Rand,
		freq:      io.NewIn("freq", dsp.Frequency(440, c.SampleRate)),
		amp:       io.NewIn("amp", dsp.Float64(1)),
		fm:        io.NewIn("freq-mod", dsp.Float64(0)),
		pw:        io.NewIn("pulse-width", dsp.Float64(1)),
		pm:        io.NewIn("phase-mod", dsp.Float64(0)),
		sync:      io.NewIn("sync", dsp.Float64(-1)),
		offset:    io.NewIn("offset", dsp.Float64(0)),
		frameSize: c.FrameSize,
	}

	io.ExposeOutputProcessor(g.newSine("sine", 1))
	io.ExposeOutputProcessor(g.newSine("sub-sine", 0.5))
	io.ExposeOutputProcessor(g.newSaw("saw", 1))
	io.ExposeOutputProcessor(g.newSaw("sub-saw", 0.5))
	io.ExposeOutputProcessor(g.newTriangle())
	io.ExposeOutputProcessor(g.newPulse("pulse", 1))
	io.ExposeOutputProcessor(g.newPulse("sub-pulse", 0.5))
	io.ExposeOutputProcessor(g.newNoise())
	io.ExposeOutputProcessor(g.newCluster())

	return NewUnit(io, nil), nil
}

type gen struct {
	rand                                *rand.Rand
	freq, amp, fm, pw, sync, pm, offset *In
	frameSize                           int
}

func (g *gen) newFrame() []float64 {
	return make([]float64, g.frameSize)
}

func (g *gen) newSine(name string, mult float64) *genSine {
	return &genSine{
		gen:   g,
		phase: g.rand.Float64() * twoPi,
		mult:  mult,
		out:   NewOut(name, g.newFrame()),
	}
}

func (g *gen) newSaw(name string, mult float64) *genSaw {
	return &genSaw{
		gen:   g,
		phase: g.rand.Float64() * twoPi,
		mult:  mult,
		out:   NewOut(name, g.newFrame()),
	}
}

func (g *gen) newPulse(name string, mult float64) *genPulse {
	return &genPulse{
		gen:   g,
		phase: g.rand.Float64() * twoPi,
		mult:  mult,
		out:   NewOut(name, g.newFrame()),
	}
}

func (g *gen) newTriangle() *genTriangle {
	return &genTriangle{
		gen:   g,
		phase: g.rand.Float64() * twoPi,
		out:   NewOut("triangle", g.newFrame()),
	}
}

func (g *gen) newNoise() *genNoise {
	return &genNoise{
		gen: g,
		out: NewOut("noise", g.newFrame()),
	}
}

func (g *gen) newCluster() *genCluster {
	return &genCluster{
		gen: g,
		out: NewOut("cluster", g.newFrame()),
	}
}

type genSine struct {
	*gen
	phase, mult, lastSync float64
	out                   *Out
}

func (o *genSine) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genSine) Out() *Out           { return o.out }

func (o *genSine) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genSine) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i) * o.mult
		amp    = o.amp.Read(i)
		fm     = o.fm.Read(i)
		pm     = o.pm.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
	)

	if o.lastSync < 0 && sync > 0 && o.phase < math.Pi/2 {
		o.phase = 0
	}

	next := dsp.Sin(o.phase + pm)
	o.phase = stepPhase(freq, fm, o.phase, o.frameSize, o.frameSize)
	o.out.Write(i, (amp*next)+offset)
	o.lastSync = sync
}

type genSaw struct {
	*gen
	phase, mult, lastSync float64
	out                   *Out
}

func (o *genSaw) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genSaw) Out() *Out           { return o.out }

func (o *genSaw) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genSaw) ProcessSample(i int) {
	var (
		amp    = o.amp.Read(i)
		fm     = o.fm.Read(i)
		freq   = o.freq.Read(i) * o.mult
		offset = o.offset.Read(i)
		pm     = o.pm.Read(i)
		sync   = o.sync.Read(i)

		next float64
	)

	if o.lastSync < 0 && sync > 0 && o.phase < math.Pi/2 {
		o.phase = 0
	}

	p := (o.phase + pm) / twoPi
	next = (2*p - 1)
	next -= blep(p, freq, fm)
	o.phase = stepPhase(freq, fm, o.phase, o.frameSize, o.frameSize)
	o.out.Write(i, (amp*next)+offset)
	o.lastSync = sync
}

type genPulse struct {
	*gen
	phase, mult, lastSync float64
	out                   *Out
}

func (o *genPulse) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genPulse) Out() *Out           { return o.out }

func (o *genPulse) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genPulse) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i) * o.mult
		amp    = o.amp.Read(i)
		fm     = o.fm.Read(i)
		pw     = math.Abs(o.pw.Read(i))
		pm     = o.pm.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)

		next float64
	)

	if o.lastSync < 0 && sync > 0 && o.phase < math.Pi/2 {
		o.phase = 0
	}

	if o.phase+pm < math.Pi*pw {
		next = 1
	} else {
		next = -1
	}
	p := (o.phase + pm) / twoPi
	next += blep(p, freq, fm)
	next -= blep(math.Mod(p+0.5, 1), freq, fm)

	o.phase = stepPhase(freq, fm, o.phase, o.frameSize, o.frameSize)
	o.out.Write(i, (amp*next)+offset)
	o.lastSync = sync
}

type genTriangle struct {
	*gen
	phase, lastSync, last float64
	out                   *Out
}

func (o *genTriangle) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genTriangle) Out() *Out           { return o.out }

func (o *genTriangle) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genTriangle) ProcessSample(i int) {
	var (
		freq   = o.freq.Read(i)
		amp    = o.amp.Read(i)
		fm     = o.fm.Read(i)
		offset = o.offset.Read(i)
		sync   = o.sync.Read(i)
		pm     = o.pm.Read(i)
		next   float64
	)

	if o.lastSync < 0 && sync > 0 && o.phase < math.Pi/2 {
		o.phase = 0
	}

	if o.phase+pm < math.Pi {
		next = 1
	} else {
		next = -1
	}
	p := (o.phase + pm) / twoPi
	next += blep(p, freq, fm)
	next -= blep(math.Mod(p+0.5, 1), freq, fm)
	next = freq*next + (1-freq)*o.last

	o.phase = stepPhase(freq, fm, o.phase, o.frameSize, o.frameSize)
	o.out.Write(i, (4*amp*next)+offset)
	o.last = next
	o.lastSync = sync
}

type genNoise struct {
	*gen
	out *Out
}

func (o *genNoise) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genNoise) Out() *Out           { return o.out }

func (o *genNoise) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genNoise) ProcessSample(i int) {
	var (
		offset = o.offset.Read(i)
		amp    = o.amp.Read(i)
	)
	o.out.Write(i, (o.rand.Float64()*2-1)*amp+offset)
}

type genCluster struct {
	*gen
	out *Out
}

func (o *genCluster) IsProcessable() bool { return o.out.ExternalNeighborCount() > 0 }
func (o *genCluster) Out() *Out           { return o.out }

func (o *genCluster) ProcessFrame(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *genCluster) ProcessSample(i int) {
	var (
		offset = o.offset.Read(i)
		amp    = o.amp.Read(i)
	)
	d := (-math.Log(o.rand.Float64()) + math.Log(o.rand.Float64())) * 0.1
	if d > 0.5 || d < -0.5 {
		o.out.Write(i, d*amp+offset)
	} else {
		o.out.Write(i, offset*amp)
	}
}

func stepPhase(freq, fm, phase float64, frameSize, n int) float64 {
	phase += (math.Abs(freq+fm) * twoPi) * (float64(frameSize) / float64(n))
	if phase >= twoPi {
		phase -= twoPi
	}
	return phase
}

func blep(p float64, freq, fm float64) float64 {
	delta := math.Abs(freq + fm)
	if p < delta {
		p /= delta
		return p + p - p*p - 1
	} else if p > 1-delta {
		p = (p - 1) / delta
		return p + p + p*p + 1
	}
	return 0
}
