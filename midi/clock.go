package midi

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rakyll/portmidi"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/unit"
)

func newClock(c unit.Config) (*unit.Unit, error) {
	var config struct {
		Device    int
		FrameRate int
	}
	if err := mapstructure.Decode(c, &config); err != nil {
		return nil, err
	}

	stream, err := portmidi.NewInputStream(portmidi.DeviceID(config.Device), int64(dsp.FrameSize))
	if err != nil {
		return nil, err
	}

	if config.FrameRate == 0 {
		config.FrameRate = 24
	}

	stop := make(chan struct{})

	io := unit.NewIO()
	clk := &clock{
		stopEvent: stop,
		stream:    stream,
		events:    eventStream(stream, stop),
		frameRate: config.FrameRate,
		out:       io.NewOut("out"),
		reset:     io.NewOut("reset"),
		start:     io.NewOut("start"),
		stop:      io.NewOut("stop"),
		spp:       io.NewOut("spp"),
	}
	return unit.NewUnit(io, "midi-clock", clk), nil
}

type clock struct {
	out, reset, start, stop, spp *unit.Out
	stream                       *portmidi.Stream
	events                       <-chan portmidi.Event
	stopEvent                    chan struct{}
	frameRate, count             int
}

const (
	midiClockTick  = 248
	midiClockReset = 250
	midiClockStart = 251
	midiClockStop  = 252
	midiClockSPP   = 242
)

func (c *clock) ProcessSample(i int) {
	if c.stream == nil {
		return
	}
	var (
		spp   float64 = -1
		stop  float64 = -1
		start float64 = -1
		reset float64 = -1
	)
	select {
	case e := <-c.events:
		if e.Status == midiClockTick || e.Status == midiClockReset {
			c.count++
		}

		switch e.Status {
		case midiClockReset:
			reset = 1
			c.count = 0
		case midiClockStart:
			start = 1
		case midiClockStop:
			stop = 1
		case midiClockSPP:
			spp = float64(e.Data1 + (e.Data2 * 127))
		}
	default:
	}

	c.start.Write(i, start)
	c.stop.Write(i, stop)
	c.reset.Write(i, reset)
	c.spp.Write(i, spp)

	if c.count%c.frameRate == 0 {
		c.out.Write(i, 1)
	} else {
		c.out.Write(i, -1)
	}
}

func (c *clock) Close() error {
	if c.stream != nil {
		if err := c.stream.Close(); err != nil {
			return err
		}
		c.stream = nil
		go func() { c.stopEvent <- struct{}{} }()
	}
	return nil
}
