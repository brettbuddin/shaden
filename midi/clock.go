package midi

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rakyll/portmidi"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/unit"
)

func newClock(creator StreamCreator) unit.BuildFunc {
	return func(c unit.Config) (*unit.Unit, error) {
		var config struct {
			Device    int
			FrameRate int
		}
		if err := mapstructure.Decode(c, &config); err != nil {
			return nil, err
		}

		stream, err := creator.NewStream(portmidi.DeviceID(config.Device), int64(dsp.FrameSize))
		if err != nil {
			return nil, err
		}

		if config.FrameRate == 0 {
			config.FrameRate = 24
		}

		io := unit.NewIO()
		clk := &clock{
			stream:    stream,
			events:    stream.Channel(),
			frameRate: config.FrameRate,
			out:       io.NewOut("out"),
			reset:     io.NewOut("reset"),
			start:     io.NewOut("start"),
			stop:      io.NewOut("stop"),
			spp:       io.NewOut("spp"),
		}
		return unit.NewUnit(io, "midi-clock", clk), nil
	}
}

type clock struct {
	out, reset, start, stop, spp *unit.Out
	stream                       Stream
	events                       <-chan portmidi.Event
	frameRate, count             int
	lastSPP                      float64
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
		stop  = -1.0
		start = -1.0
		reset = -1.0
	)
	e := <-c.events
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
		c.lastSPP = float64(e.Data1 + (e.Data2 * 127))
	}

	c.start.Write(i, start)
	c.stop.Write(i, stop)
	c.reset.Write(i, reset)
	c.spp.Write(i, c.lastSPP)

	if c.count%c.frameRate == 0 {
		c.out.Write(i, 1)
	} else {
		c.out.Write(i, -1)
	}
}

func (c *clock) Close() error {
	err := c.stream.Close()
	c.stream = nil
	return err
}
