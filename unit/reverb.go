package unit

import (
	"buddin.us/shaden/dsp"
)

var (
	maxPredelay = dsp.Duration(500).Float64()
	aTravelFreq = dsp.Frequency(0.5).Float64()
	bTravelFreq = dsp.Frequency(0.3).Float64()
)

func newReverb(name string, _ Config) (*Unit, error) {
	io := NewIO()
	r := &reverb{
		a:          io.NewIn("a", dsp.Float64(0)),
		b:          io.NewIn("b", dsp.Float64(0)),
		defuse:     io.NewIn("defuse", dsp.Float64(0.5)),
		mix:        io.NewIn("mix", dsp.Float64(0)),
		precutoff:  io.NewIn("cutoff-pre", dsp.Frequency(300)),
		postcutoff: io.NewIn("cutoff-post", dsp.Frequency(500)),
		decay:      io.NewIn("decay", dsp.Float64(0.5)),
		aOut:       io.NewOut("a"),
		bOut:       io.NewOut("b"),

		ap:          make([]*dsp.AllPass, 5),
		aAP:         make([]*dsp.AllPass, 2),
		bAP:         make([]*dsp.AllPass, 2),
		aFilter:     &dsp.SVFilter{Poles: 4},
		aPostFilter: &dsp.SVFilter{Poles: 4},
		bFilter:     &dsp.SVFilter{Poles: 4},
		bPostFilter: &dsp.SVFilter{Poles: 4},
	}

	r.ap[0] = dsp.NewAllPass(141)
	r.ap[1] = dsp.NewAllPass(201)
	r.ap[2] = dsp.NewAllPass(297)
	r.ap[3] = dsp.NewAllPass(351)
	r.ap[4] = dsp.NewAllPass(537)

	r.aPreDL = dsp.NewDelayLine(3541)
	r.aAP[0] = dsp.NewAllPass(2182)
	r.aAP[1] = dsp.NewAllPass(2690)
	r.aPostDL = dsp.NewDelayLine(4453)

	r.bPreDL = dsp.NewDelayLine(3541)
	r.bAP[0] = dsp.NewAllPass(2182)
	r.bAP[1] = dsp.NewAllPass(2690)
	r.bPostDL = dsp.NewDelayLine(4357)

	return NewUnit(io, name, r), nil
}

type reverb struct {
	a, b, defuse, mix, precutoff, postcutoff, decay *In
	aOut, bOut                                      *Out

	aPhase, bPhase       float64
	aFilter, aPostFilter *dsp.SVFilter
	bFilter, bPostFilter *dsp.SVFilter
	aPreDL, bPreDL       *dsp.DelayLine
	aPostDL, bPostDL     *dsp.DelayLine
	ap, aAP, bAP         []*dsp.AllPass
	aLast, bLast, phase  float64
}

func decayClamp(v float64) float64  { return dsp.Clamp(v, 0, 0.9) }
func defuseClamp(v float64) float64 { return dsp.Clamp(v, 0.4, 0.625) }

func (r *reverb) ProcessSample(i int) {
	var (
		mix        = r.mix.ReadSlow(i, ident)
		decay      = r.decay.ReadSlow(i, decayClamp)
		defuse     = r.defuse.ReadSlow(i, defuseClamp)
		precutoff  = r.precutoff.ReadSlow(i, ident)
		postcutoff = r.postcutoff.ReadSlow(i, ident)
	)

	a, b := r.a.Read(i), r.b.Read(i)

	d := r.ap[0].Tick(a+b, defuse)
	d = r.ap[1].Tick(d, defuse)
	d = r.ap[2].Tick(d, defuse)
	d = r.ap[3].Tick(d, defuse)
	d = r.ap[4].Tick(d, defuse)

	r.aFilter.Cutoff = precutoff
	r.bFilter.Cutoff = precutoff
	r.aPostFilter.Cutoff = postcutoff
	r.bPostFilter.Cutoff = postcutoff

	aTravel := dsp.Sin(r.aPhase)*0.0015 + 0.9
	r.aPhase += (aTravelFreq * twoPi)
	if r.aPhase >= twoPi {
		r.aPhase -= twoPi
	}

	aSig := d + (r.bLast * decay)
	aSig = r.aPreDL.TickRelative(aSig, aTravel)
	aSig, _, _ = r.aFilter.Tick(aSig)
	aSig = r.aAP[0].Tick(aSig, defuse)
	aSig = r.aAP[1].Tick(aSig*decay, defuse)
	aSig, _, _ = r.aPostFilter.Tick(aSig)

	aOut := r.aPostDL.Tick(aSig)
	aOut += r.aPostDL.ReadRelative(0.8)
	r.aLast = r.aPostDL.ReadRelative(0.6)

	bTravel := dsp.Sin(r.bPhase)*0.0015 + 0.9
	r.bPhase += (bTravelFreq * twoPi)
	if r.bPhase >= twoPi {
		r.bPhase -= twoPi
	}

	bSig := d + (r.aLast * decay)
	bSig = r.aPreDL.TickRelative(bSig, bTravel)
	bSig, _, _ = r.bFilter.Tick(bSig)
	bSig = r.bAP[0].Tick(bSig, defuse)
	bSig = r.bAP[1].Tick(bSig*decay, defuse)
	bSig, _, _ = r.bPostFilter.Tick(bSig)

	bOut := r.bPostDL.Tick(bSig)
	bOut += r.bPostDL.ReadRelative(0.8)
	r.bLast = r.bPostDL.ReadRelative(0.6)

	r.aOut.Write(i, dsp.Mix(mix, a, aOut))
	r.bOut.Write(i, dsp.Mix(mix, b, bOut))
}
