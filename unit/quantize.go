package unit

import (
	"fmt"
	"math"

	"buddin.us/musictheory"
	"buddin.us/shaden/dsp"
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

func newQuantize(name string, _ Config) (*Unit, error) {
	io := NewIO()

	tonic, err := dsp.ParsePitch("A4")
	if err != nil {
		return nil, err
	}

	q := &quantize{
		intervals: io.NewProp("intervals", &intervals{}, intervalSetter),
		in:        io.NewIn("in", dsp.Float64(0)),
		tonic:     io.NewIn("tonic", tonic),
		out:       io.NewOut("out"),
		ratios:    make([]float64, maxIntervals),
	}
	q.maybeUpdate()

	return NewUnit(io, name, q), nil
}

type quantize struct {
	intervals *Prop
	in, tonic *In
	out       *Out

	lastIntervalHash string
	ratios           []float64
	ratioCount       int
}

func (q *quantize) maybeUpdate() {
	intervals := q.intervals.Value().(*intervals)
	if q.lastIntervalHash != intervals.sig {
		for i, intvl := range intervals.intervals {
			q.ratios[i] = intvl.Ratio()
		}
		q.ratioCount = len(intervals.intervals)
		q.lastIntervalHash = intervals.sig
	}
}

func (q *quantize) ProcessSample(i int) {
	q.maybeUpdate()
	if q.ratioCount == 0 {
		return
	}
	var (
		tonic = q.tonic.Read(i)
		in    = dsp.Clamp(q.in.Read(i), 0, 1)
		n     = float64(q.ratioCount)
		idx   = math.Max(math.Min(math.Floor(n*in+0.5), n-1), 0)
	)

	q.out.Write(i, tonic*q.ratios[int(idx)])
}
