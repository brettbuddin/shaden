package midi

import (
	"time"

	"buddin.us/shaden/unit"
	"github.com/rakyll/portmidi"
)

func UnitBuilders() map[string]unit.BuildFunc {
	return map[string]unit.BuildFunc{
		"midi-clock": newClock,
		"midi-input": newInput,
	}
}

func Initialize() ([]*portmidi.DeviceInfo, error) {
	if err := portmidi.Initialize(); err != nil {
		return nil, err
	}
	info := []*portmidi.DeviceInfo{}
	for i := 0; i < portmidi.CountDevices(); i++ {
		info = append(info, portmidi.Info(portmidi.DeviceID(i)))
	}
	return info, nil
}

func Terminate() error {
	return portmidi.Terminate()
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
