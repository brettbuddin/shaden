package midi

import (
	"testing"

	"github.com/brettbuddin/shaden/unit"
	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	ch := make(chan portmidi.Event)
	creator := streamCreatorFunc(func(deviceID portmidi.DeviceID, frameSize int64) (eventStream, error) {
		return streamMock{
			events: ch,
		}, nil
	})

	go func() {
		ch <- portmidi.Event{Status: 251, Timestamp: 1}
		ch <- portmidi.Event{Status: 248, Timestamp: 2}
		ch <- portmidi.Event{Status: 252, Timestamp: 3}
		ch <- portmidi.Event{Status: 242, Data1: 1, Data2: 3, Timestamp: 4}
		ch <- portmidi.Event{Status: 250, Timestamp: 5}
	}()

	u, err := newClock(creator, blockingReceiver)(unit.NewIO("midi-clock", frameSize), newUnitConfig(nil))
	require.NoError(t, err)
	require.NotNil(t, u)

	u.ProcessFrame(5)

	var (
		out   = u.Out["out"].(*unit.Out)
		start = u.Out["start"].(*unit.Out)
		reset = u.Out["reset"].(*unit.Out)
		stop  = u.Out["stop"].(*unit.Out)
		spp   = u.Out["spp"].(*unit.Out)
	)

	require.Equal(t, 1.0, out.Read(0))
	require.Equal(t, 1.0, start.Read(0))
	require.Equal(t, -1.0, reset.Read(0))
	require.Equal(t, -1.0, stop.Read(0))
	require.Equal(t, 0.0, spp.Read(0))

	require.Equal(t, -1.0, out.Read(1))
	require.Equal(t, -1.0, start.Read(1))
	require.Equal(t, -1.0, reset.Read(1))
	require.Equal(t, -1.0, stop.Read(1))
	require.Equal(t, 0.0, spp.Read(1))

	require.Equal(t, -1.0, out.Read(2))
	require.Equal(t, -1.0, start.Read(2))
	require.Equal(t, -1.0, reset.Read(2))
	require.Equal(t, 1.0, stop.Read(2))
	require.Equal(t, 0.0, spp.Read(2))

	require.Equal(t, -1.0, out.Read(3))
	require.Equal(t, -1.0, start.Read(3))
	require.Equal(t, -1.0, reset.Read(3))
	require.Equal(t, -1.0, stop.Read(3))
	require.Equal(t, 382.0, spp.Read(3))

	require.Equal(t, 1.0, out.Read(4))
	require.Equal(t, -1.0, start.Read(4))
	require.Equal(t, 1.0, reset.Read(4))
	require.Equal(t, -1.0, stop.Read(4))
	require.Equal(t, 382.0, spp.Read(4))

	u.Close()
}
