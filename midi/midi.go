package midi

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"buddin.us/shaden/unit"
	"github.com/rakyll/portmidi"
)

// UnitBuilders returns the list of units provided by this package.
func UnitBuilders() map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"midi-clock": newClock,
		"midi-input": newInput,
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

func eventStream(s *portmidi.Stream, stop <-chan struct{}) <-chan portmidi.Event {
	ch := make(chan portmidi.Event)
	go func() {
		t := time.NewTicker(10 * time.Millisecond)

		for {
			if s == nil {
				continue
			}
			select {
			case <-stop:
				return
			case <-t.C:
				events, err := s.Read(1024)
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
