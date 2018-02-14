// filters implement easy to use audio filters.
package filters

import (
	"github.com/mattetti/audio"
	"github.com/mattetti/audio/dsp/filters"
	"github.com/mattetti/audio/dsp/windows"
)

// LowPass is a basic LowPass filter cutting off
// CutOffFreq is where the filter would be at -3db.
// TODO: param to say how efficient we want the low pass to be.
// matlab: lpFilt = designfilt('lowpassfir','PassbandFrequency',0.25, ...
//         'StopbandFrequency',0.35,'PassbandRipple',0.5, ...
//         'StopbandAttenuation',65,'DesignMethod','kaiserwin');
func LowPass(buf *audio.PCMBuffer, cutOffFreq float64) (err error) {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	s := &filters.Sinc{
		// TODO: find the right taps number to do a proper
		// audio low pass based in the sample rate
		// there should be a magical function to get that number.
		Taps:         62,
		SamplingFreq: buf.Format.SampleRate,
		CutOffFreq:   cutOffFreq,
		Window:       windows.Hamming,
	}
	fir := &filters.FIR{Sinc: s}
	buf.Floats, err = fir.LowPass(buf.AsFloat64s())
	buf.SwitchPrimaryType(audio.Float)
	return err
}

// HighPass is a basic LowPass filter cutting off
// the audio buffer frequencies below the cutOff frequency.
func HighPass(buf *audio.PCMBuffer, cutOff float64) (err error) {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	s := &filters.Sinc{
		Taps:         62,
		SamplingFreq: buf.Format.SampleRate,
		CutOffFreq:   cutOff,
		Window:       windows.Blackman,
	}
	fir := &filters.FIR{Sinc: s}
	buf.Floats, err = fir.HighPass(buf.AsFloat64s())
	buf.SwitchPrimaryType(audio.Float)
	return err
}
