// Package midi provides midi units.
package midi

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rakyll/portmidi"

	"buddin.us/shaden/unit"
)

const sendInterval = 10 * time.Millisecond

var defaultStreamCreator = streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
	s, err := portmidi.NewInputStream(deviceID, frameSize)
	return &stream{Stream: s, stop: make(chan struct{})}, err
})

// UnitBuilders returns the list of units provided by this package.
func UnitBuilders() map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"midi-clock": newClock(defaultStreamCreator, nonBlockingReceiver),
		"midi-input": newInput(defaultStreamCreator, nonBlockingReceiver),
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

// eventReceiver is a function that receives values from a channel
type eventReceiver func(<-chan portmidi.Event) portmidi.Event

// nonBlockingReceiver is an eventReceiver that doesn't wait for an event to become available if the sender isn't ready;
// in the event the sender isn't ready it will return an empty event.
func nonBlockingReceiver(events <-chan portmidi.Event) portmidi.Event {
	select {
	case e := <-events:
		return e
	default:
		return portmidi.Event{}
	}
}

// blockingReceiver is a eventReceiver that always waits for the sender to be ready. Used in testing.
func blockingReceiver(events <-chan portmidi.Event) portmidi.Event { return <-events }

// streamCreator provides new Streams
type streamCreator interface {
	NewStream(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error)
}

// eventStream is a stream of PortMIDI events that can be closed
type eventStream interface {
	Channel(time.Duration) <-chan portmidi.Event
	io.Closer
}

type stream struct {
	*portmidi.Stream
	stop chan struct{}
}

// Channel returns a channel that emits PortMIDI events. Every call to Channel should be terminated by a call to Close;
// failure to do so will result in a leaked goroutine.
func (s stream) Channel(interval time.Duration) <-chan portmidi.Event {
	ch := make(chan portmidi.Event)
	go func() {
		t := time.NewTicker(interval)

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

type streamCreatorFunc func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error)

func (f streamCreatorFunc) NewStream(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
	return f(deviceID, frameSize)
}
