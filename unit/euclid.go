package unit

import (
	"math"

	"buddin.us/shaden/dsp"
)

const maxEuclidLayers = 256

func newEuclid(name string, _ Config) (*Unit, error) {
	io := NewIO()
	e := &euclid{
		clock:       io.NewIn("clock", dsp.Float64(-1)),
		span:        io.NewIn("span", dsp.Float64(5)),
		fill:        io.NewIn("fill", dsp.Float64(2)),
		offset:      io.NewIn("offset", dsp.Float64(0)),
		counts:      make([]int, maxEuclidLayers),
		remainders:  make([]int, maxEuclidLayers),
		pattern:     make([]bool, maxEuclidLayers),
		lastTrigger: -1,
		out:         io.NewOut("out"),
	}
	return NewUnit(io, name, e), nil
}

type euclid struct {
	clock, span, fill, offset *In
	lastTrigger               float64
	lastSpan, lastFill        int

	out *Out

	pattern            []bool
	counts, remainders []int
	idx, lastIdx       int
}

func spanMin(v float64) float64 { return math.Max(v, 1) }

func (e *euclid) ProcessSample(i int) {
	var (
		span           = int(e.span.ReadSlow(i, spanMin))
		fill           = int(e.fill.ReadSlow(i, ident))
		offset         = e.offset.ReadSlowInt(i, func(v int) int { return (v + span) % span })
		trig           = e.clock.Read(i)
		out    float64 = -1
	)

	if e.lastSpan != span || e.lastFill != fill {
		for i := range e.pattern {
			e.counts[i] = 0
			e.remainders[i] = 0
			e.pattern[i] = false
		}
		euclidean(e.pattern, e.counts, e.remainders, span, fill)
	}
	e.lastSpan = span
	e.lastFill = fill

	if isTrig(e.lastTrigger, trig) {
		e.idx = (e.idx + 1) % span
	}
	idx := (e.idx + offset + span) % span
	if e.pattern[idx] && e.idx == e.lastIdx {
		out = 1
	}
	e.lastIdx = e.idx
	e.lastTrigger = trig

	e.out.Write(i, out)
}

func euclidean(pattern []bool, counts, remainders []int, n, p int) {
	if p > n {
		p = n
	}
	div := n - p
	remainders[0] = p

	var lvl int
	for {
		counts[lvl] = div / remainders[lvl]
		remainders[lvl+1] = div % remainders[lvl]
		div = remainders[lvl]
		lvl++
		if remainders[lvl] <= 1 {
			break
		}
	}

	counts[lvl] = div
	var idx int
	build(pattern, counts, remainders, lvl, &idx)
}

func build(pattern []bool, counts, remainders []int, lvl int, idx *int) {
	switch {
	case lvl > len(counts)-1:
		return
	case lvl == -1:
		pattern[*idx] = false
		*idx++
	case lvl == -2:
		pattern[*idx] = true
		*idx++
	default:
		for i := 0; i < counts[lvl]; i++ {
			build(pattern, counts, remainders, lvl-1, idx)
		}
		if remainders[lvl] != 0 {
			build(pattern, counts, remainders, lvl-2, idx)
		}
	}
}
