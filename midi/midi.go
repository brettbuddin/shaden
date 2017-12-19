// Package midi provides midi units.
package midi

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"buddin.us/shaden/unit"
	"github.com/rakyll/portmidi"
)

var defaultStreamCreator = streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (Stream, error) {
	s, err := portmidi.NewInputStream(deviceID, frameSize)
	return &stream{Stream: s, stop: make(chan struct{})}, err
})

// UnitBuilders returns the list of units provided by this package.
func UnitBuilders() map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"midi-clock": newClock(defaultStreamCreator),
		"midi-input": newInput(defaultStreamCreator),
	}
}

// Initialize initializes portmidi and returns the list of devices on the system.
func Initialize() (DeviceList, error) {
	if err := portmidi.Initialize(); err != nil {
		return nil, err
	}
	info := []*portmidi.DeviceInfo{}
	for i := 0; i < portmidi.CountDevices(); i++ {
		info = append(info, portmidi.Info(portmidi.DeviceID(i)))
	}
	return DeviceList(info), nil
}

// Terminate terminates portmidi.
func Terminate() error {
	return portmidi.Terminate()
}

// DeviceList is a list of MIDI devices.
type DeviceList []*portmidi.DeviceInfo

func (l DeviceList) String() string {
	out := bytes.NewBuffer(nil)
	if len(l) > 0 {
		for i, d := range l {
			dirs := []string{}
			if d.IsInputAvailable {
				dirs = append(dirs, "input")
			}
			if d.IsOutputAvailable {
				dirs = append(dirs, "output")
			}
			fmt.Fprintf(out, "%d: %s (%s)\n", i, d.Name, strings.Join(dirs, "/"))
		}
	} else {
		fmt.Fprintln(out, "(none)")
	}
	return out.String()
}

// StreamCreator provides new Streams
type StreamCreator interface {
	NewStream(deviceID portmidi.DeviceID, frameSize int64) (Stream, error)
}

// Stream is a stream of PortMIDI events that can be closed
type Stream interface {
	Channel() <-chan portmidi.Event
	io.Closer
}

type stream struct {
	*portmidi.Stream
	stop chan struct{}
}

// Channel returns a channel that emits PortMIDI events. Every call to Channel should be terminated by a call to Close;
// failure to do so will result in a leaked goroutine.
func (s stream) Channel() <-chan portmidi.Event {
	ch := make(chan portmidi.Event)
	go func() {
		t := time.NewTicker(10 * time.Millisecond)

		for {
			select {
			case <-s.stop:
				close(s.stop)
				return
			case <-t.C:
				events, err := s.Stream.Read(1024)
				if err != nil {
					continue
				}
				for i := range events {
					ch <- events[i]
				}
			default:
				select {
				case <-s.stop:
					return
				case ch <- portmidi.Event{}:
				}
			}
		}
	}()
	return ch
}

// Close closes the underlying PortMIDI stream
func (s stream) Close() error {
	s.stop <- struct{}{}
	return s.Stream.Close()
}

type streamCreatorFunc func(deviceID portmidi.DeviceID, frameSize int64) (Stream, error)

func (f streamCreatorFunc) NewStream(deviceID portmidi.DeviceID, frameSize int64) (Stream, error) {
	return f(deviceID, frameSize)
}
