package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/brettbuddin/shaden/engine/portaudio"
	"github.com/brettbuddin/shaden/errors"
)

const (
	backendPortAudio = "portaudio"
	backendStdout    = "stdout"
)

// Config is a structure for storing all the parsed flags.
type Config struct {
	Seed                 int64
	HTTPAddr             string
	REPL                 bool
	FrameSize            int
	SampleRate           float64
	SingleSampleDisabled bool
	FadeIn               int
	Gain                 float64

	Backend string

	DeviceList      bool
	deviceIn        string
	deviceOut       string
	DeviceLatency   string
	DeviceFrameSize int

	ScriptPath string
}

func (c Config) DeviceIn() (portaudio.DeviceSelection, error) {
	return deviceSelection(c.deviceIn)
}

func (c Config) DeviceOut() (portaudio.DeviceSelection, error) {
	return deviceSelection(c.deviceOut)
}

func deviceSelection(v string) (portaudio.DeviceSelection, error) {
	if v == "none" {
		return portaudio.DeviceNone, nil
	}

	vInt, err := strconv.Atoi(v)
	if err != nil {
		return portaudio.DeviceNone, err
	}

	return portaudio.DeviceSelection(vInt), nil
}

func parseArgs(args []string) (Config, error) {
	var cfg Config

	set := flag.NewFlagSet("shaden", flag.ContinueOnError)

	set.Int64Var(&cfg.Seed, "seed", 0, "random seed")
	set.IntVar(&cfg.FrameSize, "frame", 256, "frame size used within the synthesis engine")
	set.StringVar(&cfg.HTTPAddr, "addr", ":5000", "http address to serve")
	set.BoolVar(&cfg.REPL, "repl", false, "REPL")
	set.Float64Var(&cfg.SampleRate, "samplerate", 44.1, "sample rate (8, 22.05, 44.1, 48.0)")
	set.BoolVar(&cfg.SingleSampleDisabled, "disable-single-sample", false, "disables single-sample mode for feedback loops")
	set.IntVar(&cfg.FadeIn, "fade-in", 100, "Duration of fade-in (milliseconds) once output signal is detected")
	set.Float64Var(&cfg.Gain, "gain", 0, "gain decibels (dB)")

	set.BoolVar(&cfg.DeviceList, "device-list", false, "list all devices")
	set.StringVar(&cfg.deviceIn, "device-in", "0", "input device")
	set.StringVar(&cfg.deviceOut, "device-out", "1", "output device")
	set.StringVar(&cfg.DeviceLatency, "device-latency", "low", "latency setting for audio device")
	set.IntVar(&cfg.DeviceFrameSize, "device-frame", 1024, "frame size used when writing to audio device")

	set.StringVar(&cfg.Backend, "backend", "portaudio", "driver (portaudio, stdout)")

	err := set.Parse(args)

	if len(set.Args()) > 0 {
		cfg.ScriptPath = set.Arg(0)
	}

	switch cfg.Backend {
	case "portaudio":
	case "stdout":
	default:
		return cfg, errors.Errorf("unknown backend %q", cfg.Backend)
	}

	if cfg.HTTPAddr == "" {
		return cfg, errors.Errorf("addr cannot be empty")
	}

	if cfg.DeviceFrameSize < cfg.FrameSize {
		return cfg, errors.Errorf("device frame size cannot be less than %d", cfg.FrameSize)
	}

	if cfg.DeviceFrameSize%cfg.FrameSize != 0 {
		return cfg, errors.Errorf("frame size (%d) must be a multiple of %d", cfg.DeviceFrameSize, cfg.FrameSize)
	}

	cfg.SampleRate = cfg.SampleRate * 1000

	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}

	return cfg, err
}
