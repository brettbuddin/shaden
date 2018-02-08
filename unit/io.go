package unit

import (
	"fmt"
	"sync/atomic"

	"buddin.us/shaden/dsp"
)

var idCount uint32

// IO is the registry of inputs, outputs and properties for a Module
type IO struct {
	ID, Type string
	Prop     map[string]*Prop
	In       map[string]*In
	Out      map[string]Output
}

// NewIO returns a new IO
func NewIO(typ string) *IO {
	io := &IO{
		ID:   fmt.Sprintf("%s-%d", typ, idCount),
		Type: typ,
		Prop: map[string]*Prop{},
		In:   map[string]*In{},
		Out:  map[string]Output{},
	}
	atomic.AddUint32(&idCount, 1)
	return io
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
	in := NewIn(name, v)
	io.In[in.Name] = in
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

// ExposeOutputProcessor registers a new output that is also a Processor
func (io *IO) ExposeOutputProcessor(o OutputProcessor) {
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
