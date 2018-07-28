package main

import (
	"flag"
	"time"

	"github.com/brettbuddin/shaden/errors"
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

	DeviceList      bool
	DeviceIn        int
	DeviceOut       int
	DeviceLatency   string
	DeviceFrameSize int

	ScriptPath string
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

	set.BoolVar(&cfg.DeviceList, "device-list", false, "list all devices")
	set.IntVar(&cfg.DeviceIn, "device-in", 0, "input device")
	set.IntVar(&cfg.DeviceOut, "device-out", 1, "output device")
	set.StringVar(&cfg.DeviceLatency, "device-latency", "low", "latency setting for audio device")
	set.IntVar(&cfg.DeviceFrameSize, "device-frame", 1024, "frame size used when writing to audio device")

	err := set.Parse(args)

	if len(set.Args()) > 0 {
		cfg.ScriptPath = set.Arg(0)
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
