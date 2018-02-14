// given a PCM audio file, convert it to mono and decimates it.
// Because of nyquist law, we can't simply average or drop samples otherwise we will have alisasing.
// A proper decimation is needed http://dspguru.com/dsp/faqs/multirate/decimation
// Decimation is useful to reduce the amount of data to process.
// Max hearing frequency would be around 20kHz, so we need to low pass to remove anything above 20kHz so
// we don't get any aliasing.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattetti/audio"
	"github.com/mattetti/audio/generator"
	"github.com/mattetti/audio/transforms"
	"github.com/mattetti/audio/wav"
)

var (
	fileFlag   = flag.String("file", "", "file to downsample (copy will be done)")
	factorFlag = flag.Int("factor", 2, "The decimator factor divides the sampling rate")
	outputFlag = flag.String("format", "aiff", "output format, aiff or wav")
)

func main() {
	flag.Parse()

	if *fileFlag == "" {
		freq := audio.RootA
		fs := 44100
		fmt.Printf("Resampling from %dHz to %dHz back to %dHz\n", fs, fs / *factorFlag, fs)

		// generate a wave sine
		osc := generator.NewOsc(generator.WaveSine, freq, fs)
		data := osc.Signal(fs * 4)
		buf := audio.NewPCMFloatBuffer(data, audio.FormatMono4410016bBE)
		// drop the sample rate
		if err := transforms.Decimate(buf, *factorFlag); err != nil {
			panic(err)
		}

		fmt.Println("bit crushing to 8 bit sound")
		// the bitcrusher switches the data range to PCM scale
		transforms.BitCrush(buf, 8)
		transforms.Resample(buf, float64(fs))

		// sideBuf := buf.Clone()
		// sideBuf.SwitchPrimaryType(audio.Float)
		// truncate
		// sideBuf.Floats = sideBuf.Floats[:1024]
		// transforms.NormalizeMax(sideBuf)
		// export as a gnuplot binary file
		// if err := presenters.GnuplotText(sideBuf, "decimator.dat"); err != nil {
		// 	panic(err)
		// }
		// if err := presenters.CSV(sideBuf, "data.csv"); err != nil {
		// 	panic(err)
		// }

		// encode the sound file
		o, err := os.Create("resampled.wav")
		if err != nil {
			panic(err)
		}
		defer o.Close()
		e := wav.NewEncoder(o, buf.Format.SampleRate, buf.Format.BitDepth, 1, 1)
		if err := e.Write(buf); err != nil {
			panic(err)
		}
		e.Close()
		fmt.Println("checkout resampled.wav")
		return
	}

	/*
		ext := filepath.Ext(*fileFlag)
		var codec string
		switch strings.ToLower(ext) {
		case ".aif", ".aiff":
			codec = "aiff"
		case ".wav", ".wave":
			codec = "wav"
		default:
			fmt.Printf("files with extension %s not supported\n", ext)
		}

		f, err := os.Open(*fileFlag)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var monoFrames audio.Frames
		var sampleRate int
		var sampleSize int
		switch codec {
		case "aiff":
			d := aiff.NewDecoder(f)
			frames, err := d.Frames()
			if err != nil {
				panic(err)
			}
			sampleRate = d.SampleRate
			sampleSize = int(d.BitDepth)
			monoFrames = audio.ToMonoFrames(frames)

		case "wav":
			info, frames, err := wav.NewDecoder(f, nil).ReadFrames()
			if err != nil {
				panic(err)
			}
			sampleRate = int(info.SampleRate)
			sampleSize = int(info.BitsPerSample)
			monoFrames = audio.ToMonoFrames(frames)
		}

		fmt.Printf("undersampling -> %s file at %dHz to %d samples (%d)\n", codec, sampleRate, sampleRate / *factorFlag, sampleSize)

		switch sampleRate {
		case 44100:
		case 48000:
		default:
			log.Fatalf("input sample rate of %dHz not supported", sampleRate)
		}

		amplitudesF := make([]float64, len(monoFrames))
		for i, f := range monoFrames {
			amplitudesF[i] = float64(f[0])
		}

		// low pass filter before we drop some samples to avoid aliasing
		s := &filters.Sinc{Taps: 62, SamplingFreq: sampleRate, CutOffFreq: float64(sampleRate / 2), Window: windows.Blackman}
		fir := &filters.FIR{Sinc: s}
		filtered, err := fir.LowPass(amplitudesF)
		if err != nil {
			panic(err)
		}
		frames := make([][]int, len(amplitudesF) / *factorFlag)
		for i := 0; i < len(frames); i++ {
			frames[i] = []int{int(filtered[i**factorFlag])}
		}

		of, err := os.Create("resampled.aiff")
		if err != nil {
			panic(err)
		}
		defer of.Close()
		aiffe := aiff.NewEncoder(of, sampleRate / *factorFlag, sampleSize, 1)
		aiffe.Frames = frames
		if err := aiffe.Write(); err != nil {
			panic(err)
		}
	*/
}
