package portaudio

import (
	"bytes"
	"fmt"

	portaudio "github.com/gordonklaus/portaudio"
	"github.com/pkg/errors"
)

const (
	latencyLow  = "low"
	latencyHigh = "high"
)

// PortAudio is a wrapper for a portaudio client.
type PortAudio struct {
	inDevice, outDevice *portaudio.DeviceInfo
	stream              *portaudio.Stream
	params              portaudio.StreamParameters
}

// Intialize initializes portaudio and returns the list of devices on the machine.
func Initialize() (DeviceList, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, errors.Wrap(err, "initializing portaudio")
	}
	list, err := portaudio.Devices()
	return DeviceList(list), err
}

// DeviceList is a list of portaudio devices.
type DeviceList []*portaudio.DeviceInfo

func (l DeviceList) String() string {
	out := bytes.NewBuffer(nil)
	if len(l) > 0 {
		for i, d := range l {
			fmt.Fprintf(out, "%d: %s\n", i, d.Name)
		}
	} else {
		fmt.Fprintln(out, "(none)")
	}
	return out.String()
}

// Terminate terminates portaudio. This is called after all client's have shut down.
func Terminate() error {
	return portaudio.Terminate()
}

// New returns a new PortAudio.
func New(inDeviceIndex, outDeviceIndex int, latency string, frameSize, sampleRate int) (*PortAudio, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}
	if inDeviceIndex >= len(devices) {
		return nil, fmt.Errorf("input device index out of range")
	}
	if outDeviceIndex >= len(devices) {
		return nil, fmt.Errorf("output device index out of range")
	}

	var (
		params  portaudio.StreamParameters
		in, out = devices[inDeviceIndex], devices[outDeviceIndex]
	)

	switch latency {
	case latencyHigh:
		params = portaudio.HighLatencyParameters(in, out)
	case latencyLow:
		params = portaudio.LowLatencyParameters(in, out)
	default:
		return nil, fmt.Errorf("invalid latency setting: %q", latency)
	}
	params.Input.Channels = 1
	params.Output.Channels = 2
	params.SampleRate = float64(sampleRate)
	params.FramesPerBuffer = frameSize

	return &PortAudio{
		params:    params,
		inDevice:  in,
		outDevice: out,
	}, nil
}

// Devices returns the devices currently in use by portaudio.
func (pa *PortAudio) Devices() (in *portaudio.DeviceInfo, out *portaudio.DeviceInfo) {
	return pa.inDevice, pa.outDevice
}

// FrameSize returns the low-level frame size used by PortAudio.
func (pa *PortAudio) FrameSize() int {
	return pa.params.FramesPerBuffer
}

// Start starts the portaudio stream.
func (pa *PortAudio) Start(callback func([]float32, [][]float32)) error {
	var err error
	pa.stream, err = portaudio.OpenStream(pa.params, callback)
	if err != nil {
		return err
	}
	return pa.stream.Start()
}

// Stop stops the portaudio stream.
func (pa *PortAudio) Stop() error {
	if err := pa.stream.Stop(); err != nil {
		return err
	}
	return pa.stream.Close()
}
