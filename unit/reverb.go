package unit

import (
	"buddin.us/shaden/dsp"
)

func newReverb(io *IO, c Config) (*Unit, error) {
	var (
		aTravelFreq = dsp.Frequency(0.5, c.SampleRate).Float64()
		bTravelFreq = dsp.Frequency(0.3, c.SampleRate).Float64()
		r           = &reverb{
			a:              io.NewIn("a", dsp.Float64(0)),
			b:              io.NewIn("b", dsp.Float64(0)),
			defuse:         io.NewIn("defuse", dsp.Float64(0.625)),
			mix:            io.NewIn("mix", dsp.Float64(0)),
			precutoff:      io.NewIn("cutoff-pre", dsp.Frequency(300, c.SampleRate)),
			postcutoff:     io.NewIn("cutoff-post", dsp.Frequency(500, c.SampleRate)),
			decay:          io.NewIn("decay", dsp.Float64(0.7)),
			size:           io.NewIn("size", dsp.Float64(0.1)),
			shiftSemitones: io.NewIn("shift-semitones", dsp.Float64(0)),
			aOut:           io.NewOut("a"),
			bOut:           io.NewOut("b"),

			ap:          make([]*dsp.AllPass, 4),
			aAP:         make([]*dsp.AllPass, 2),
			bAP:         make([]*dsp.AllPass, 2),
			aFilter:     &dsp.SVFilter{Poles: 1},
			aPostFilter: &dsp.SVFilter{Poles: 2},
			bFilter:     &dsp.SVFilter{Poles: 1},
			bPostFilter: &dsp.SVFilter{Poles: 2},
			blockA:      &dsp.DCBlock{},
			blockB:      &dsp.DCBlock{},
			shiftA:      dsp.NewPitchShift(),
			shiftB:      dsp.NewPitchShift(),

			aTravelFreq: aTravelFreq,
			bTravelFreq: bTravelFreq,
		}
	)

	r.ap[0] = dsp.NewAllPass(1170)
	r.ap[1] = dsp.NewAllPass(1510)
	r.ap[2] = dsp.NewAllPass(2370)
	r.ap[3] = dsp.NewAllPass(3510)

	r.aPreDL = dsp.NewDelayLine(3541)
	r.aAP[0] = dsp.NewAllPass(21820)
	r.aAP[1] = dsp.NewAllPass(26900)
	r.aPostDL = dsp.NewDelayLine(4453)

	r.bPreDL = dsp.NewDelayLine(3541)
	r.bAP[0] = dsp.NewAllPass(21820)
	r.bAP[1] = dsp.NewAllPass(26900)
	r.bPostDL = dsp.NewDelayLine(4353)

	return NewUnit(io, r), nil
}

type reverb struct {
	a, b, defuse, mix, precutoff, postcutoff, decay, size, shiftSemitones *In
	aOut, bOut                                                            *Out

	aPhase, bPhase       float64
	aFilter, aPostFilter *dsp.SVFilter
	bFilter, bPostFilter *dsp.SVFilter
	aPreDL, bPreDL       *dsp.DelayLine
	aPostDL, bPostDL     *dsp.DelayLine
	ap, aAP, bAP         []*dsp.AllPass
	aLast, bLast         float64
	blockA, blockB       *dsp.DCBlock
	shiftA, shiftB       *dsp.PitchShift

	aTravelFreq, bTravelFreq float64
}

func decayClamp(v float64) float64      { return dsp.Clamp(v, 0, 0.99) }
func defuseClamp(v float64) float64     { return dsp.Clamp(v, 0.4, 0.625) }
func pitchShiftClamp(v float64) float64 { return dsp.Clamp(v, -12, 12) }
func sizeClamp(v float64) float64       { return dsp.Clamp(v, 0.01, 1) }

func (r *reverb) ProcessSample(i int) {
	var (
		mix            = r.mix.ReadSlow(i, ident)
		decay          = r.decay.ReadSlow(i, decayClamp)
		defuse         = r.defuse.ReadSlow(i, defuseClamp)
		precutoff      = r.precutoff.ReadSlow(i, ident)
		postcutoff     = r.postcutoff.ReadSlow(i, ident)
		size           = r.size.ReadSlow(i, sizeClamp)
		shiftSemitones = r.shiftSemitones.ReadSlow(i, pitchShiftClamp)
	)

	a, b := r.a.Read(i), r.b.Read(i)

	d := r.ap[0].TickRelative(a+b, defuse, size)
	d = r.ap[1].TickRelative(d, defuse, size)
	d = r.ap[2].TickRelative(d, defuse, size)
	d = r.ap[3].TickRelative(d, defuse, size)

	r.aFilter.Cutoff = precutoff
	r.bFilter.Cutoff = precutoff
	r.aPostFilter.Cutoff = postcutoff
	r.bPostFilter.Cutoff = postcutoff

	aSig := d + (r.shiftA.TickSemitones(r.bLast, shiftSemitones) * decay)

	aTravel := dsp.Sin(r.aPhase)*0.01 + 0.9
	advanceLFO(&r.aPhase, r.aTravelFreq)
	aSig = r.aPreDL.TickRelative(aSig, aTravel*size)
	_, aSig, _ = r.aFilter.Tick(aSig)

	aSig = r.aAP[0].TickRelative(aSig, defuse, size)
	aSig = r.aAP[1].TickRelative(aSig, defuse, size)
	_, aSig, _ = r.aPostFilter.Tick(aSig)

	aOut := r.aPostDL.TickRelative(aSig, decay)
	r.aLast = aOut

	bSig := d + (r.shiftB.TickSemitones(r.aLast, shiftSemitones) * decay)

	bTravel := dsp.Sin(r.bPhase)*0.01 + 0.9
	advanceLFO(&r.bPhase, r.bTravelFreq)
	bSig = r.bPreDL.TickRelative(bSig, bTravel*size)
	_, bSig, _ = r.bFilter.Tick(bSig)

	bSig = r.bAP[0].TickRelative(bSig, defuse, size)
	bSig = r.bAP[1].TickRelative(bSig, defuse, size)
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
