package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/brettbuddin/shaden/engine"
	"github.com/brettbuddin/shaden/engine/portaudio"
	"github.com/brettbuddin/shaden/errors"
	"github.com/brettbuddin/shaden/midi"
	"github.com/brettbuddin/shaden/runtime"
)

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "", 0)

	if err := run(cfg, logger); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(cfg Config, logger *log.Logger) error {
	devices, err := portaudio.Initialize()
	if err != nil {
		return errors.Wrap(err, "initializing portaudio")
	}
	defer func() {
		if err := portaudio.Terminate(); err != nil {
			logger.Fatal(err)
		}
	}()

	midiDevices, err := midi.Initialize()
	if err != nil {
		return errors.Wrap(err, "initializing portmidi")
	}
	defer func() {
		if err := midi.Terminate(); err != nil {
			logger.Fatal(err)
		}
	}()

	if cfg.DeviceList {
		fmt.Println("Audio Devices")
		fmt.Println(devices)
		fmt.Println("MIDI Devices")
		fmt.Println(midiDevices)
		return nil
	}

	rand.Seed(cfg.Seed)

	// Create the engine
	backend, err := portaudio.New(
		cfg.DeviceIn,
		cfg.DeviceOut,
		cfg.DeviceLatency,
		cfg.DeviceFrameSize,
		int(cfg.SampleRate),
	)
	if err != nil {
		return errors.Wrap(err, "creating portaudio backend")
	}
	opts := []engine.Option{engine.WithFadeIn(cfg.FadeIn)}
	if cfg.SingleSampleDisabled {
		opts = append(opts, engine.WithSingleSampleDisabled())
	}
	e, err := engine.New(backend, cfg.FrameSize, opts...)
	if err != nil {
		return errors.Wrap(err, "engine create failed")
	}
	printPreamble(backend, cfg.Seed)

	// Create the lisp runtime
	run, err := runtime.New(e, logger)
	if err != nil {
		return errors.Wrap(err, "start lisp runtime failed")
	}

	// Start the HTTP server
	go func() {
		mux := http.NewServeMux()
		runtime.AddHandler(mux, run)
		if err := http.ListenAndServe(cfg.HTTPAddr, mux); err != nil {
			logger.Fatal(err)
		}
	}()

	// Start the engine
	go e.Run()
	go func() {
		for err := range e.Errors() {
			logger.Println("engine error:", err)
		}
	}()
	defer e.Stop()

	if cfg.ScriptPath != "" {
		if err := run.Load(cfg.ScriptPath); err != nil {
			return errors.Wrap(err, "file eval failed")
		}
	}

	replDone := make(chan struct{})
	if cfg.REPL {
		go run.REPL(replDone)
	}

	select {
	case <-replDone:
	case <-waitForSignal():
	}

	return nil
}

func waitForSignal() <-chan struct{} {
	sigs := make(chan os.Signal)
	done := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(done)
	}()
	return done
}

func printPreamble(pa *portaudio.PortAudio, seed int64) {
	inDevice, outDevice := pa.Devices()
	fmt.Println("PID:", os.Getpid())
	fmt.Println("Seed:", seed)
	fmt.Printf(
		"Input Device: %s (%s/%s)\n",
		inDevice.Name,
		inDevice.DefaultLowOutputLatency,
		inDevice.DefaultHighInputLatency,
	)
	fmt.Printf(
		"Output Device: %s (%s/%s)\n",
		outDevice.Name,
		outDevice.DefaultLowOutputLatency,
		outDevice.DefaultHighOutputLatency,
	)
}
