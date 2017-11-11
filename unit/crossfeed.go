package unit

import "buddin.us/lumen/dsp"

func newCrossfeed(name string, _ Config) (*Unit, error) {
	io := NewIO()
	return NewUnit(io, name, &crossfeed{
		a:      io.NewIn("a", dsp.Float64(0)),
		b:      io.NewIn("b", dsp.Float64(0)),
		amount: io.NewIn("amount", dsp.Float64(0)),
		aOut:   io.NewOut("a"),
		bOut:   io.NewOut("b"),
	}), nil
}

type crossfeed struct {
	a, b, amount *In
	aOut, bOut   *Out
}

func (c *crossfeed) ProcessSample(i int) {
	amt := dsp.Clamp(c.amount.Read(i), 0, 1)
	a, b := c.a.Read(i), c.b.Read(i)
	c.aOut.Write(i, a+(amt*b))
	c.bOut.Write(i, b+(amt*a))
}
