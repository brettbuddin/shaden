package unit

import (
	"fmt"
	"math"

	"buddin.us/lumen/dsp"
	"buddin.us/musictheory"
)

const (
	maxIntervals = 32
)

type intervals struct {
	intervals []musictheory.Interval
	sig       string
}

func intervalSetter(p *Prop, v interface{}) error {
	slice, ok := v.([]interface{})
	if !ok {
		return InvalidPropValueError{Prop: p, Value: v}
	}
	if l := len(slice); l > maxIntervals {
		return fmt.Errorf("number of intervals %v exceeds maximum allowed %v", l, maxIntervals)
	}

	value := &intervals{}
	for _, e := range slice {
		intvl := e.(musictheory.Interval)
		value.intervals = append(value.intervals, intvl)
		value.sig += intvl.String()
	}
	p.value = value
	return nil
}

func tonicSetter(p *Prop, v interface{}) error {
	if _, ok := v.(*musictheory.Pitch); !ok {
		return InvalidPropValueError{Prop: p, Value: v}
	}
	p.value = v
	return nil
}

func newQuantize(name string, _ Config) (*Unit, error) {
	io := NewIO()

	defaultTonic, err := musictheory.ParsePitch("A4")
	if err != nil {
		return nil, err
	}

	q := &quantize{
		intervals: io.NewProp("intervals", &intervals{}, intervalSetter),
		tonic:     io.NewProp("tonic", defaultTonic, tonicSetter),
		in:        io.NewIn("in", dsp.Float64(0)),
		out:       io.NewOut("out"),
		pitches:   make([]dsp.Hz, maxIntervals),
	}
	q.maybeUpdate()

	return NewUnit(io, name, q), nil
}

type quantize struct {
	intervals, tonic *Prop
	in               *In
	out              *Out

	lastTonic        *musictheory.Pitch
	lastIntervalHash string
	pitches          []dsp.Hz
	pitchCount       int
}

func (q *quantize) maybeUpdate() {
	var (
		intervals = q.intervals.Value().(*intervals)
		tonic     = q.tonic.Value().(*musictheory.Pitch)
	)
	if (q.lastTonic == nil || !q.lastTonic.Eq(*tonic)) || q.lastIntervalHash != intervals.sig {
		for i, intvl := range intervals.intervals {
			p := tonic.Transpose(intvl).(musictheory.Pitch)
			q.pitches[i] = dsp.Frequency(p.Freq())
		}
		q.pitchCount = len(intervals.intervals)
		q.lastTonic = tonic
		q.lastIntervalHash = intervals.sig
	}
}

func (q *quantize) ProcessSample(i int) {
	q.maybeUpdate()
	if q.pitchCount == 0 {
		return
	}
	in := dsp.Clamp(q.in.Read(i), 0, 1)

	n := float64(q.pitchCount)
	idx := math.Floor(n*in + 0.5)
	idx = math.Min(idx, n-1)
	idx = math.Max(idx, 0)

	q.out.Write(i, q.pitches[int(idx)].Float64())
}
