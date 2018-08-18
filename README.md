## Shaden

[![Build Status](https://travis-ci.org/brettbuddin/shaden.svg?branch=ci)](https://travis-ci.org/brettbuddin/shaden)
[![GoDoc](https://godoc.org/github.com/brettbuddin/shaden?status.svg)](https://godoc.org/github.com/brettbuddin/shaden)
[![Go Report Card](https://goreportcard.com/badge/github.com/brettbuddin/shaden)](https://goreportcard.com/report/github.com/brettbuddin/shaden)

<img src="https://github.com/brettbuddin/shaden/raw/master/extra/shaden-snail.png" width="200">

Shaden is a modular audio synthesizer. Patches for the synthesizer are written in a Lisp dialect. A REPL and HTTP
interface are provided for interacting with the synthesizer in real-time. I started this project as a way of learning
more about digital signal processing and music theory. **Consider this an art project**. 

The name "shaden" comes from the *Cycle of Galand* book series by Edward W. Robertson.

### Highlights

- Lisp interpreter for creating patches
- [Large collection of builtin Units](https://github.com/brettbuddin/shaden/wiki/Units)
- [Music theory primitives](https://github.com/brettbuddin/shaden/wiki/Values#music-theory)
- MIDI controller and clock input
- Single-sample feedback loops
- Editor plugins for sending expressions to the synth for evaluation:
    - [Vim](extra/shaden.vim)
    - [VS Code](https://github.com/semrekkers/shaden-vscode)

## Dependencies

- [Go 1.9](http://golang.org)+
- [PortAudio](http://www.portaudio.com/)
- [PortMIDI](http://portmedia.sourceforge.net/portmidi/)

On macOS you can install these dependencies with: `brew install go portaudio portmidi`

## Getting Started

### Install

    $ go get -u github.com/brettbuddin/shaden

### CLI Usage

#### REPL

    $ shaden -repl
    > (define gen (unit/gen))
    > (-> gen (table :freq (hz 300)))
    > (emit (<- gen :sine))

#### Load File

    $ shaden examples/frequency-modulation.lisp

#### HTTP

    $ shaden examples/krell.lisp
    $ curl -X POST http://127.0.0.1:5000/eval -d "(define source (unit/gen)) ; ..."

This is my preferred way of interacting with the synthesizer. I've written a small Vim plugin that can send over
snippets of Lisp code to the program for evaluation. [You can find that plugin here.](extra/shaden.vim)

The HTTP interface is limited to Lisp evaluation at the moment, but I have hopes of providing an API for direct graph
manipulation via HTTP.

### Lisp

For a more information about the Lisp dialect bundled with Shaden, [check out the wiki](https://github.com/brettbuddin/shaden/wiki).

## Examples

The best way to get to know the way patching works in Shaden is to look at the [examples directory](examples). As far as
sounds that can be created with it:

- The synth was used to create the intro music for the [GothamGo 2017 conference videos](https://www.youtube.com/playlist?list=PLeGxIOPLk9ELp7dx6A0gtvjbc99dU2kq-)
- I frequently post patches I've created [on Instagram](https://www.instagram.com/brettbuddin/) and recordings [on Bandcamp](https://returnnil.bandcamp.com).
