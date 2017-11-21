## Shaden

[![Build Status](https://travis-ci.org/brettbuddin/shaden.svg?branch=ci)](https://travis-ci.org/brettbuddin/shaden)
[![GoDoc](https://godoc.org/buddin.us/shaden?status.svg)](https://godoc.org/buddin.us/shaden)
[![Go Report Card](https://goreportcard.com/badge/github.com/brettbuddin/shaden)](https://goreportcard.com/report/github.com/brettbuddin/shaden)
[![Coverage Status](https://coveralls.io/repos/github/brettbuddin/shaden/badge.svg?branch=master)](https://coveralls.io/github/brettbuddin/shaden?branch=master)

Shaden is a modular audio synthesizer. Patches for the synthesizer are written in a Lisp dialect. A REPL and HTTP
interface is provided for interacting with the synthesizer in real-time. I started this project as a way of learning
more about digital signal processing and music theory. **Consider this an art project**. 

The name "shaden" comes from the *Cycle of Galand* book series by Edward W. Robertson.

### Highlights

- Lisp interpreter for creating patches
- [Large collection of builtin Units](https://github.com/brettbuddin/shaden/wiki/Units)
- [Music theory primitives](https://github.com/brettbuddin/shaden/wiki/Values#music-theory)
- MIDI controller and clock input
- Single-sample feedback loops
- Vim plugin for sending snippets of code over to the synth for evaluation

## Dependencies

- [Go 1.9](http://golang.org)+
- [PortAudio](http://www.portaudio.com/)
- [PortMIDI](http://portmedia.sourceforge.net/portmidi/)

On macOS you can install these dependencies with: `brew install go portaudio portmidi`

## Getting Started

### Install

    $ go get -u buddin.us/shaden
	$ shaden -h
	Usage of shaden:
  	-addr string
        	http address to serve (default ":5000")
  	-device-frame int
        	frame size used when writing to audio device (default 1024)
  	-device-in int
        	input device
  	-device-latency string
        	latency setting for audio device (default "low")
  	-device-list
        	list all devices
  	-device-out int
        	output device (default 1)
  	-repl
        	REPL
  	-seed int
        	random seed
	flag: help requested

### CLI Usage

#### REPL

    $ shaden -repl
    > (define gen (unit/gen))
    > (-> gen (table :freq (hz 300)))
    > (emit (<- gen :sine))

#### Load File

    $ shaden examples/frequency-modulation.lisp

#### HTTP

    $ shaden
    $ curl -X POST http://127.0.0.1:5000/eval -d "(define source (unit/gen)) ; ..."

This is my preferred way of interacting with the synthesizer. I've written a small Vim plugin that can send over
snippets of Lisp code to the program for evaluation. You can get [that plugin here](extra/shaden.vim).

The HTTP interface is limited to Lisp evaluation at the moment, but I have hopes of providing an API for direct graph
manipulation via HTTP.

### Lisp

For a more information about the Lisp dialect bundled with Shaden, [check out the wiki](https://github.com/brettbuddin/shaden/wiki).

## Examples

The best way to get to know the way patching works in Shaden is to look at the [examples directory](examples). As far as
sounds that can be created with it:

- The synth was used to create the intro music for the [GothamGo 2017 conference videos](https://www.youtube.com/watch?v=l_FkVIPerzE)
- I frequently post patches I've created with Shaden [on Instagram](https://www.instagram.com/brettbuddin/)
