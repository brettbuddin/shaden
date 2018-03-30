package midi

import (
	"testing"
	"time"

	"buddin.us/shaden/unit"
	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/require"
)

type midiOutput interface {
	unit.Output
	unit.CondProcessor
	unit.FrameProcessor
	unit.SampleProcessor
}

// Ensure midi-in outputs conform to the interfaces necessary for processing.
var _ = []midiOutput{
	&pitch{},
	&pitchRaw{},
	&gate{},
	&bend{},
	&cc{},
}

func TestInput_Pitch(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(1)

	pitchOut := u.Out["1/pitch"].(*pitch)
	pitchOut.ProcessFrame(2)

	require.Equal(t, 0.005932552501147361, pitchOut.out.Read(0))
	require.Equal(t, 0.005932552501147361, pitchOut.out.Read(1))

	u.Close()
}

func TestInput_PitchRaw(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(1)

	pitchOut := u.Out["1/pitchraw"].(*pitchRaw)
	pitchOut.ProcessFrame(2)

	require.Equal(t, 60.0, pitchOut.out.Read(0))
	require.Equal(t, 60.0, pitchOut.out.Read(1))

	u.Close()
}

func TestInput_Gate_NoteOff(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 1}
		ch <- portmidi.Event{Status: 128, Data1: 60, Data2: 0, Timestamp: 2}
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 3}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(3)

	gateOut := u.Out["1/gate"].(*gate)
	gateOut.ProcessFrame(7)

	require.Equal(t, -1.0, gateOut.out.Read(0))
	require.Equal(t, 1.0, gateOut.out.Read(1))
	require.Equal(t, -1.0, gateOut.out.Read(2))
	require.Equal(t, 1.0, gateOut.out.Read(3))
	require.Equal(t, 1.0, gateOut.out.Read(4))
	require.Equal(t, 1.0, gateOut.out.Read(5))
	require.Equal(t, 1.0, gateOut.out.Read(6))

	u.Close()
}

func TestInput_Gate_NoNoteOff(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 1}
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 0, Timestamp: 2}
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 3}
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 4}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(4)

	gateOut := u.Out["1/gate"].(*gate)
	gateOut.ProcessFrame(7)

	require.Equal(t, -1.0, gateOut.out.Read(0))
	require.Equal(t, 1.0, gateOut.out.Read(1))
	require.Equal(t, -1.0, gateOut.out.Read(2))
	require.Equal(t, 1.0, gateOut.out.Read(3))
	require.Equal(t, -1.0, gateOut.out.Read(4))
	require.Equal(t, -1.0, gateOut.out.Read(5))
	require.Equal(t, -1.0, gateOut.out.Read(6))

	u.Close()
}

func TestInput_Gate_Rolling(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 1}
		ch <- portmidi.Event{Status: 144, Data1: 61, Data2: 127, Timestamp: 2}
		ch <- portmidi.Event{Status: 144, Data1: 60, Data2: 127, Timestamp: 3}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(3)

	gateOut := u.Out["1/gate"].(*gate)
	gateOut.ProcessFrame(7)

	require.Equal(t, -1.0, gateOut.out.Read(0))
	require.Equal(t, 1.0, gateOut.out.Read(1))
	require.Equal(t, -1.0, gateOut.out.Read(2))
	require.Equal(t, 1.0, gateOut.out.Read(3))
	require.Equal(t, 1.0, gateOut.out.Read(4))
	require.Equal(t, 1.0, gateOut.out.Read(5))
	require.Equal(t, 1.0, gateOut.out.Read(6))

	u.Close()
}

func TestInput_CC(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 176, Data1: 1, Data2: 127, Timestamp: 1}
		ch <- portmidi.Event{Status: 176, Data1: 2, Data2: 127, Timestamp: 2}
		ch <- portmidi.Event{Status: 176, Data1: 1, Data2: 0, Timestamp: 3}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(3)

	one := u.Out["1/cc/1"].(*cc)
	two := u.Out["1/cc/2"].(*cc)
	one.ProcessFrame(3)
	two.ProcessFrame(3)

	require.Equal(t, 1.0, one.out.Read(0))
	require.Equal(t, 0.0, two.out.Read(0))

	require.Equal(t, 1.0, one.out.Read(1))
	require.Equal(t, 1.0, two.out.Read(1))

	require.Equal(t, 0.0, one.out.Read(2))
	require.Equal(t, 1.0, two.out.Read(2))

	u.Close()
}

func TestInput_Bend(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 224, Data1: 0, Data2: 127, Timestamp: 1}
		ch <- portmidi.Event{Status: 224, Data1: 0, Data2: 100, Timestamp: 2}
		ch <- portmidi.Event{Status: 224, Data1: 0, Data2: 0, Timestamp: 3}
	}()

	u, err := newInput(creator, blockingReceiver)(unit.NewIO("midi-input", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(3)
	bendOut := u.Out["1/bend"].(*bend)
	bendOut.ProcessFrame(3)

	require.Equal(t, 1.0, bendOut.out.Read(0))
	require.Equal(t, 0.7874015748031497, bendOut.out.Read(1))
	require.Equal(t, 0.0, bendOut.out.Read(2))

	u.Close()
}

type streamMock struct {
	events chan portmidi.Event
	err    error
}

func (s streamMock) Channel(time.Duration) <-chan portmidi.Event { return s.events }
func (s streamMock) Close() error                                { return s.err }
