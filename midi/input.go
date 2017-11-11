package midi

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rakyll/portmidi"

	"buddin.us/lumen/dsp"
	"buddin.us/lumen/unit"
	"buddin.us/musictheory"
)

var pitches = map[int]float64{}

func init() {
	p := musictheory.NewPitch(musictheory.C, musictheory.Natural, 0)
	for i := 12; i < 127; i++ {
		pitches[i] = dsp.Frequency(p.Freq()).Float64()
		p = p.Transpose(musictheory.Minor(2)).(musictheory.Pitch)
	}
}

func newInput(c unit.Config) (*unit.Unit, error) {
	var config struct {
		Device   int
		Channels []int
	}
	if err := mapstructure.Decode(c, &config); err != nil {
		return nil, err
	}

	stream, err := portmidi.NewInputStream(portmidi.DeviceID(config.Device), int64(dsp.FrameSize))
	if err != nil {
		return nil, err
	}

	if len(config.Channels) == 0 {
		config.Channels = []int{1}
	}

	stop := make(chan struct{})
	io := unit.NewIO()

	ctrl := &input{
		stopEvent: stop,
		stream:    stream,
		eventChan: eventStream(stream, stop),
		events:    make([]portmidi.Event, dsp.FrameSize),
	}

	for _, ch := range config.Channels {
		io.ExposeOutProcessor(ctrl.newPitch(ch))
		io.ExposeOutProcessor(ctrl.newPitchRaw(ch))
		io.ExposeOutProcessor(ctrl.newGate(ch))
		io.ExposeOutProcessor(ctrl.newBend(ch))
		for i := 1; i < 128; i++ {
			io.ExposeOutProcessor(ctrl.newCC(ch, i))
		}
	}

	return unit.NewUnit(io, "midi-input", ctrl), nil
}

type input struct {
	stream    *portmidi.Stream
	eventChan <-chan portmidi.Event
	events    []portmidi.Event
	stopEvent chan struct{}
}

func (in *input) newPitch(ch int) *pitch {
	return &pitch{
		input: in,
		ch:    ch,
		out:   unit.NewOut(fmt.Sprintf("%d/pitch", ch), make([]float64, dsp.FrameSize)),
	}
}

func (in *input) newPitchRaw(ch int) *pitchRaw {
	return &pitchRaw{
		input: in,
		ch:    ch,
		out:   unit.NewOut(fmt.Sprintf("%d/pitchraw", ch), make([]float64, dsp.FrameSize)),
	}
}

func (in *input) newGate(ch int) *gate {
	return &gate{
		input:   in,
		stateFn: gateUp,
		state:   &gateState{which: -1, chOffset: int64(ch) - 1},
		out:     unit.NewOut(fmt.Sprintf("%d/gate", ch), make([]float64, dsp.FrameSize)),
	}
}

func (in *input) newCC(ch, num int) *cc {
	return &cc{
		input: in,
		ch:    int64(ch),
		num:   int64(num),
		out:   unit.NewOut(fmt.Sprintf("%d/cc/%d", ch, num), make([]float64, dsp.FrameSize)),
	}
}

func (in *input) newBend(ch int) *bend {
	return &bend{
		input: in,
		ch:    int64(ch),
		out:   unit.NewOut(fmt.Sprintf("%d/bend", ch), make([]float64, dsp.FrameSize)),
	}
}

func (in *input) IsProcessable() bool {
	return true
}

func (in *input) ProcessSample(i int) {
	if in.stream == nil {
		return
	}
	select {
	case e := <-in.eventChan:
		in.events[i] = e
	default:
		in.events[i] = portmidi.Event{}
	}
}

func (in *input) Close() error {
	if in.stream != nil {
		if err := in.stream.Close(); err != nil {
			return err
		}
		in.stream = nil
		go func() { in.stopEvent <- struct{}{} }()
	}
	return nil
}

const (
	statusNoteOn     = 144
	statusNoteOff    = 128
	statusCC         = 176
	statusPitchWheel = 224
)

type pitch struct {
	*input
	ch   int
	freq float64
	out  *unit.Out
}

func (o *pitch) IsProcessable() bool {
	return o.out.DestinationCount() > 0
}

func (o *pitch) ProcessSample(i int) {
	if e := o.events[i]; e.Status == int64(statusNoteOn+o.ch-1) {
		if v, ok := pitches[int(e.Data1)]; ok && e.Data2 > 0 {
			o.freq = v
		}
	}
	o.out.Write(i, o.freq)
}

func (o *pitch) Out() *unit.Out {
	return o.out
}

type pitchRaw struct {
	*input
	ch   int
	note float64
	out  *unit.Out
}

func (o *pitchRaw) IsProcessable() bool {
	return o.out.DestinationCount() > 0
}

func (o *pitchRaw) ProcessSample(i int) {
	if e := o.events[i]; e.Status == int64(statusNoteOn+o.ch-1) {
		o.note = float64(e.Data1)
	}
	o.out.Write(i, o.note)
}

func (o *pitchRaw) Out() *unit.Out {
	return o.out
}

type gate struct {
	*input
	state   *gateState
	stateFn gateStateFunc
	out     *unit.Out
}

type gateState struct {
	event           portmidi.Event
	which, chOffset int64
	value           float64
}

type gateStateFunc func(*gateState) gateStateFunc

func gateRolling(s *gateState) gateStateFunc {
	s.value = -1
	return gateDown
}

func gateDown(s *gateState) gateStateFunc {
	s.value = 1
	if s.event.Status == 0 && s.event.Timestamp == 0 {
		return gateDown
	}

	which := s.event.Data1

	switch s.event.Status {
	case statusNoteOn + s.chOffset:
		if s.event.Data2 > 0 {
			if which != s.which {
				s.which = which
				return gateRolling
			}
			s.which = -1
			return gateUp
		}
		if which == s.which {
			s.which = -1
			return gateUp
		}
	case statusNoteOff + s.chOffset:
		if which == s.which {
			s.which = -1
			return gateUp
		}
	}
	return gateDown
}

func gateUp(s *gateState) gateStateFunc {
	s.value = -1
	if s.event.Status == 0 && s.event.Timestamp == 0 {
		return gateUp
	}
	if s.event.Status == statusNoteOn+s.chOffset && s.event.Data2 > 0 {
		s.which = s.event.Data1
		return gateDown
	}
	return gateUp
}

func (o *gate) IsProcessable() bool {
	return o.out.DestinationCount() > 0
}

func (o *gate) ProcessSample(i int) {
	o.state.event = o.events[i]
	o.stateFn = o.stateFn(o.state)
	o.out.Write(i, o.state.value)
}

func (o *gate) Out() *unit.Out {
	return o.out
}

type cc struct {
	*input
	ch, num int64
	value   float64
	out     *unit.Out
}

func (o *cc) IsProcessable() bool {
	return o.out.DestinationCount() > 0
}

func (o *cc) Process(n int) {
	for i := 0; i < n; i++ {
		o.ProcessSample(i)
	}
}

func (o *cc) ProcessSample(i int) {
	if e := o.events[i]; e.Status == statusCC+o.ch-1 && e.Data1 == o.num {
		o.value = float64(e.Data2) / 127
	}
	o.out.Write(i, o.value)
}

func (o *cc) Out() *unit.Out {
	return o.out
}

type bend struct {
	*input
	ch    int64
	value float64
	out   *unit.Out
}

func (o *bend) IsProcessable() bool {
	return o.out.DestinationCount() > 0
}

func (o *bend) ProcessSample(i int) {
	if e := o.events[i]; e.Status == statusPitchWheel+o.ch-1 && e.Data1 == 0 {
		o.value = float64(e.Data2) / 127
	}
	o.out.Write(i, o.value)
}

func (o *bend) Out() *unit.Out {
	return o.out
}
