package unit

import "buddin.us/shaden/dsp"

// IO is the registry of inputs, outputs and properties for a Module
type IO struct {
	Prop map[string]*Prop
	In   map[string]*In
	Out  map[string]Output
}

// NewIO returns a new IO
func NewIO() *IO {
	return &IO{
		Prop: map[string]*Prop{},
		In:   map[string]*In{},
		Out:  map[string]Output{},
	}
}

// NewProp registers a new property
func (io *IO) NewProp(name string, v interface{}, setter func(*Prop, interface{}) error) *Prop {
	switch assert := v.(type) {
	case int:
		v = float64(assert)
	}
	p := &Prop{
		name:   name,
		value:  v,
		setter: setter,
	}
	io.Prop[p.name] = p
	return p
}

// NewIn registers a new input
func (io *IO) NewIn(name string, v dsp.Valuer) *In {
	f := newFrame()
	in := &In{
		Name:          name,
		frame:         f,
		constantFrame: f,
	}
	io.In[in.Name] = in
	in.setNormal(v)
	return in
}

// NewOut registers a new output
func (io *IO) NewOut(name string) *Out {
	return io.newOut(name, newFrame())
}

// NewOutWithFrame registers a new output that has a specific frame
func (io *IO) NewOutWithFrame(name string, f []float64) *Out {
	return io.newOut(name, f)
}

// ExposeOutProcessor registers a new output that is also a Processor
func (io *IO) ExposeOutProcessor(o OutputProcessor) {
	io.Out[o.Out().Name] = o
}

func (io *IO) newOut(name string, f []float64) *Out {
	o := &Out{
		Name:  name,
		frame: f,
	}
	io.Out[name] = o
	return o
}
