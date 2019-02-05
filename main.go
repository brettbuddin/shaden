package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/brettbuddin/shaden/engine"
	"github.com/brettbuddin/shaden/engine/portaudio"
	"github.com/brettbuddin/shaden/engine/stdout"
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

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	rand.Seed(cfg.Seed)

	dest := os.Stdout
	if cfg.REPL {
		dest = os.Stderr
	}

	var (
		backend engine.Backend
		logger  = log.New(dest, "", 0)
	)

	switch cfg.Backend {
	case backendPortAudio:
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
			logger.Println("Audio Devices")
			logger.Println(devices)
			logger.Println("MIDI Devices")
			logger.Println(midiDevices)
			return nil
		}

		// Create the engine
		paBackend, err := portaudio.New(
			cfg.DeviceIn,
			cfg.DeviceOut,
			cfg.DeviceLatency,
			cfg.DeviceFrameSize,
			int(cfg.SampleRate),
		)
		if err != nil {
			return errors.Wrap(err, "creating portaudio backend")
		}

		printPreamble(paBackend, logger, cfg.Seed)

		backend = paBackend
	case backendStdout:
		logger = log.New(os.Stderr, "", 0)
		backend = stdout.New(os.Stdout, cfg.FrameSize, int(cfg.SampleRate))
	default:
		return errors.Errorf("unknown backend %q", cfg.Backend)
	}

	opts := []engine.Option{
		engine.WithFadeIn(cfg.FadeIn),
		engine.WithGain(dbToFloat(cfg.Gain)),
	}
	if cfg.SingleSampleDisabled {
		opts = append(opts, engine.WithSingleSampleDisabled())
	}
	e, err := engine.New(backend, cfg.FrameSize, opts...)
	if err != nil {
		return errors.Wrap(err, "engine create failed")
	}

	// Create the lisp runtime
	run, err := runtime.New(e, logger)
	if err != nil {
		return errors.Wrap(err, "start lisp runtime failed")
	}

	// Start the HTTP server
	go func() {
		mux := http.NewServeMux()

		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		mux.Handle("/debug/pprof/block", pprof.Handler("block"))
		mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

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

func printPreamble(pa *portaudio.PortAudio, logger *log.Logger, seed int64) {
	inDevice, outDevice := pa.Devices()
	logger.Println("PID:", os.Getpid())
	logger.Println("Seed:", seed)

	if inDevice != nil {
		logger.Printf(
			"Input Device: %s (%s/%s)\n",
			inDevice.Name,
			inDevice.DefaultLowOutputLatency,
			inDevice.DefaultHighInputLatency,
		)
	} else {
		logger.Println("Input Device: none")
	}

	logger.Printf(
		"Output Device: %s (%s/%s)\n",
		outDevice.Name,
		outDevice.DefaultLowOutputLatency,
		outDevice.DefaultHighOutputLatency,
	)
}

func dbToFloat(v float64) float32 {
	return float32(math.Pow(10, 0.05*v))
}
