// shaden is a modular synthesizer.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/pkg/errors"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/engine"
	"buddin.us/shaden/engine/portaudio"
	"buddin.us/shaden/midi"
	"buddin.us/shaden/runtime"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var (
		set                  = flag.NewFlagSet("shaden", flag.ContinueOnError)
		seed                 = set.Int64("seed", 0, "random seed")
		deviceList           = set.Bool("device-list", false, "list all devices")
		deviceIn             = set.Int("device-in", 0, "input device")
		deviceOut            = set.Int("device-out", 1, "output device")
		deviceLatency        = set.String("device-latency", "low", "latency setting for audio device")
		deviceFrameSize      = set.Int("device-frame", 1024, "frame size used when writing to audio device")
		httpAddr             = set.String("addr", ":5000", "http address to serve")
		repl                 = set.Bool("repl", false, "REPL")
		singleSampleDisabled = set.Bool("disable-single-sample", false, "disables single-sample mode for feedback loops")
		logger               = log.New(os.Stdout, "", 0)
	)

	if err := set.Parse(args); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if *deviceFrameSize < dsp.FrameSize {
		return errors.Errorf("device frame size cannot be less than %d", dsp.FrameSize)
	}

	if *deviceFrameSize%dsp.FrameSize != 0 {
		return errors.Errorf("frame size (%d) must be a multiple of %d", *deviceFrameSize, dsp.FrameSize)
	}

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

	if *deviceList {
		fmt.Println("Audio Devices")
		fmt.Println(devices)
		fmt.Println("MIDI Devices")
		fmt.Println(midiDevices)
		return nil
	}

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rand.Seed(*seed)

	// Create the engine
	backend, err := portaudio.New(*deviceIn, *deviceOut, *deviceLatency, *deviceFrameSize)
	if err != nil {
		return errors.Wrap(err, "creating portaudio backend")
	}
	e, err := engine.New(backend, *singleSampleDisabled)
	if err != nil {
		return errors.Wrap(err, "engine create failed")
	}
	printPreamble(backend, *seed)

	// Create the lisp runtime
	run, err := runtime.New(e, logger)
	if err != nil {
		return errors.Wrap(err, "start lisp runtime failed")
	}

	// Start the HTTP server
	go func() {
		if err := serve(*httpAddr, run); err != nil {
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

	if len(set.Args()) > 0 {
		if err := run.Load(set.Arg(0)); err != nil {
			return errors.Wrap(err, "file eval failed")
		}
	}

	replDone := make(chan struct{})
	if *repl {
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

func serve(addr string, run *runtime.Runtime) error {
	http.HandleFunc("/eval", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := run.Eval(body); err != nil {
			fmt.Fprintf(w, "%s", err)
			return
		}
		fmt.Fprintf(w, "OK")
	})
	return http.ListenAndServe(addr, nil)
}
