package unit

import (
	"buddin.us/shaden/dsp"
)

var (
	aTravelFreq = dsp.Frequency(0.5).Float64()
	bTravelFreq = dsp.Frequency(0.3).Float64()
)

func newReverb(name string, _ Config) (*Unit, error) {
	io := NewIO()
	r := &reverb{
		a:          io.NewIn("a", dsp.Float64(0)),
		b:          io.NewIn("b", dsp.Float64(0)),
		defuse:     io.NewIn("defuse", dsp.Float64(0.625)),
		mix:        io.NewIn("mix", dsp.Float64(0)),
		precutoff:  io.NewIn("cutoff-pre", dsp.Frequency(300)),
		postcutoff: io.NewIn("cutoff-post", dsp.Frequency(500)),
		decay:      io.NewIn("decay", dsp.Float64(0.5)),
		aOut:       io.NewOut("a"),
		bOut:       io.NewOut("b"),

		ap:          make([]*dsp.AllPass, 4),
		aAP:         make([]*dsp.AllPass, 2),
		bAP:         make([]*dsp.AllPass, 2),
		aFilter:     &dsp.SVFilter{Poles: 1},
		aPostFilter: &dsp.SVFilter{Poles: 2},
		bFilter:     &dsp.SVFilter{Poles: 1},
		bPostFilter: &dsp.SVFilter{Poles: 2},
		blockA:      &dsp.DCBlock{},
		blockB:      &dsp.DCBlock{},
	}

	r.ap[0] = dsp.NewAllPass(117)
	r.ap[1] = dsp.NewAllPass(151)
	r.ap[2] = dsp.NewAllPass(237)
	r.ap[3] = dsp.NewAllPass(351)

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
	aLast, bLast         float64
	blockA, blockB       *dsp.DCBlock
}

func decayClamp(v float64) float64  { return dsp.Clamp(v, 0, 0.99) }
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

	r.aFilter.Cutoff = precutoff
	r.bFilter.Cutoff = precutoff
	r.aPostFilter.Cutoff = postcutoff
	r.bPostFilter.Cutoff = postcutoff

	aSig := d + (r.bLast * decay)

	aTravel := dsp.Sin(r.aPhase)*0.007 + 0.9
	advanceLFO(&r.aPhase, aTravelFreq)
	aSig = r.aPreDL.TickRelative(aSig, aTravel)
	_, aSig, _ = r.aFilter.Tick(aSig)

	aSig = r.aAP[0].Tick(aSig, defuse)
	aSig = r.aAP[1].Tick(aSig, defuse)
	_, aSig, _ = r.aPostFilter.Tick(aSig)

	aOut := r.aPostDL.TickRelative(aSig, decay)
	r.aLast = aOut

	bSig := d + (r.aLast * decay)

	bTravel := dsp.Sin(r.bPhase)*0.007 + 0.9
	advanceLFO(&r.bPhase, bTravelFreq)
	bSig = r.bPreDL.TickRelative(bSig, bTravel)
	_, bSig, _ = r.bFilter.Tick(bSig)

	bSig = r.bAP[0].Tick(bSig, defuse)
	bSig = r.bAP[1].Tick(bSig, defuse)
	_, bSig, _ = r.bPostFilter.Tick(bSig)

	bOut := r.bPostDL.TickRelative(bSig, decay)
	r.bLast = bOut

	r.aOut.Write(i, r.blockA.Tick(dsp.Mix(mix, a, aOut)))
	r.bOut.Write(i, r.blockB.Tick(dsp.Mix(mix, b, bOut)))
}

func advanceLFO(phase *float64, freq float64) {
	*phase += (freq * twoPi)
	if *phase >= twoPi {
		*phase -= twoPi
	}
}
