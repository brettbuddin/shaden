package engine

import (
	"fmt"

	"buddin.us/lumen/dsp"
	"github.com/gordonklaus/portaudio"
)

const (
	latencyLow  = "low"
	latencyHigh = "high"
)

// Initialize initializes PortAudio and returns the list of devices available on the machine
func Initialize() ([]*portaudio.DeviceInfo, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, nil
	}
	return portaudio.Devices()
}

// Terminate closes PortAudio
func Terminate() error {
	return portaudio.Terminate()
}

type portAudio struct {
	inDevice, outDevice *portaudio.DeviceInfo
	stream              *portaudio.Stream
	params              portaudio.StreamParameters
}

func newPortAudio(inDeviceIndex, outDeviceIndex int, latency string, frameSize int) (*portAudio, error) {
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
	params.SampleRate = float64(dsp.SampleRate)
	params.FramesPerBuffer = frameSize

	return &portAudio{
		params:    params,
		inDevice:  in,
		outDevice: out,
	}, nil
}

func (pa *portAudio) Start(callback interface{}) error {
	var err error
	pa.stream, err = portaudio.OpenStream(pa.params, callback)
	if err != nil {
		return err
	}
	return pa.stream.Start()
}

func (pa *portAudio) Stop() error {
	if err := pa.stream.Stop(); err != nil {
		return err
	}
	return pa.stream.Close()
}
