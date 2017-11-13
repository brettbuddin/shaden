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
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/gordonklaus/portaudio"
	"github.com/pkg/errors"
	"github.com/rakyll/portmidi"

	"buddin.us/shaden/dsp"
	"buddin.us/shaden/engine"
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
	)

	if err := set.Parse(args); err != nil {
		return errors.Wrap(err, "parsing flags")
	}

	if *deviceFrameSize < dsp.FrameSize {
		return errors.Errorf("device frame size cannot be less than %d", dsp.FrameSize)
	}

	logger := log.New(os.Stdout, "", 0)

	devices, err := engine.Initialize()
	if err != nil {
		return errors.Wrap(err, "engine initialization failed")
	}
	defer func() {
		if err := engine.Terminate(); err != nil {
			logger.Println(err)
			os.Exit(1)
		}
	}()

	midiDevices, err := midi.Initialize()
	if err != nil {
		return errors.Wrap(err, "midi initialization failed")
	}
	defer func() {
		if err := midi.Terminate(); err != nil {
			logger.Println(err)
			os.Exit(1)
		}
	}()

	if *deviceList {
		printDeviceList(devices, midiDevices)
		return nil
	}

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}
	rand.Seed(*seed)

	e, err := engine.New(*deviceIn, *deviceOut, *deviceLatency, *deviceFrameSize, *singleSampleDisabled)
	if err != nil {
		return errors.Wrap(err, "engine create failed")
	}
	printPreamble(e, *seed)

	go e.Run()
	go func() {
		for err := range e.Errors() {
			logger.Println("engine error:", err)
		}
	}()
	defer e.Stop()

	run, err := runtime.New(e, logger)
	if err != nil {
		return errors.Wrap(err, "start lisp runtime failed")
	}
	if len(set.Args()) > 0 {
		if err := run.Load(set.Arg(0)); err != nil {
			return errors.Wrap(err, "file eval failed")
		}
	}
	if *repl {
		go func() {
			if err := serve(*httpAddr, run); err != nil {
				logger.Fatal(err)
			}
		}()
		run.REPL()
	} else {
		return serve(*httpAddr, run)
	}
	return nil
}

func printPreamble(e *engine.Engine, seed int64) {
	inDevice, outDevice := e.Devices()
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

func printDeviceList(devices []*portaudio.DeviceInfo, midiDevices []*portmidi.DeviceInfo) {
	fmt.Println("Audio Devices")
	if len(devices) > 0 {
		for i, d := range devices {
			fmt.Printf("%d: %s\n", i, d.Name)
		}
	} else {
		fmt.Println("(none)")
	}
	fmt.Println("\nMIDI Devices")
	if len(midiDevices) > 0 {
		for i, d := range midiDevices {
			dirs := []string{}
			if d.IsInputAvailable {
				dirs = append(dirs, "input")
			}
			if d.IsOutputAvailable {
				dirs = append(dirs, "output")
			}
			fmt.Printf("%d: %s (%s)\n", i, d.Name, strings.Join(dirs, "/"))
		}
	} else {
		fmt.Println("(none)")
	}
}
