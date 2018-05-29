package midi

import (
	"time"

	"github.com/rakyll/portmidi"

	"github.com/brettbuddin/shaden/unit"
)

func newClock(creator streamCreator, receiver eventReceiver) func(*unit.IO, unit.Config) (*unit.Unit, error) {
	return func(io *unit.IO, c unit.Config) (*unit.Unit, error) {
		var config struct {
			Rate      int
			Device    int
			FrameRate int
		}
		if err := c.Decode(&config); err != nil {
			return nil, err
		}

		stream, err := creator.NewStream(portmidi.DeviceID(config.Device), int64(c.FrameSize))
		if err != nil {
			return nil, err
		}

		if config.FrameRate == 0 {
			config.FrameRate = 24
		}

		if config.Rate == 0 {
			config.Rate = 10
		}

		return unit.NewUnit(io, &clock{
			stream:    stream,
			eventChan: stream.Channel(time.Duration(config.Rate) * time.Millisecond),
			receiver:  receiver,
			frameRate: config.FrameRate,
			out:       io.NewOut("out"),
			reset:     io.NewOut("reset"),
			start:     io.NewOut("start"),
			stop:      io.NewOut("stop"),
			spp:       io.NewOut("spp"),
		}), nil
	}
}

type clock struct {
	out, reset, start, stop, spp *unit.Out
	stream                       eventStream
	eventChan                    <-chan portmidi.Event
	receiver                     eventReceiver
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
	e := c.receiver(c.eventChan)
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
