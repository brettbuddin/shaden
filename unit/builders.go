package unit

import "github.com/mitchellh/mapstructure"

var (
	unaryModules = map[string]unaryOp{
		"abs":    unaryAbs,
		"ceil":   unaryCeil,
		"floor":  unaryFloor,
		"invert": unaryInv,
		"not":    unaryNOT,
		"noop":   unaryNoop,
	}

	binaryModules = map[string]binaryOp{
		"and":   binaryAND,
		"diff":  binaryDiff,
		"div":   binaryDiv,
		"gt":    binaryGT,
		"imply": binaryIMPLY,
		"lt":    binaryLT,
		"max":   binaryMax,
		"min":   binaryMin,
		"mod":   binaryMod,
		"mult":  binaryMult,
		"nand":  binaryNAND,
		"nor":   binaryNOR,
		"or":    binaryOR,
		"sum":   binarySum,
		"xnor":  binaryXNOR,
		"xor":   binaryXOR,
	}

	normalModules = map[string]nameBuildFunc{
		"adjust":             newAdjust,
		"chebyshev":          newChebyshev,
		"clip":               newClip,
		"clock":              newClock,
		"clock-div":          newClockDiv,
		"clock-mult":         newClockMult,
		"cluster":            newCluster,
		"cond":               newCond,
		"count":              newCount,
		"debug":              newDebug,
		"decimate":           newDecimate,
		"delay":              newDelay,
		"demux":              newDemux,
		"dynamics":           newDynamics,
		"euclid":             newEuclid,
		"filter":             newFilter,
		"fold":               newFold,
		"gate":               newGate,
		"gate-mix":           newGateMix,
		"gate-series":        newGateSeries,
		"gen":                newGen,
		"low-gen":            newLowGen,
		"lag":                newLag,
		"latch":              newLatch,
		"lerp":               newInterpolate,
		"midi-hz":            newMIDIToHz,
		"mix":                newMix,
		"mux":                newMux,
		"overload":           newOverload,
		"pan":                newPan,
		"panmix":             newPanMix,
		"pitch":              newPitch,
		"quantize":           newQuantize,
		"random-series":      newRandomSeries,
		"reverb":             newReverb,
		"slope":              newSlope,
		"smooth":             newSmooth,
		"stages":             newStages,
		"switch":             newSwitch,
		"tape":               newTape,
		"toggle":             newToggle,
		"transpose":          newTranspose,
		"transpose-interval": newTransposeInterval,
		"val-gate":           newValToGate,
		"sample":             newWAVSample,
		"xfade":              newCrossfade,
		"xfeed":              newCrossfeed,
	}
)

// BuildFunc is a constructor function for Units.
type BuildFunc func(Config) (*Unit, error)

// Config is a map that's used to provide configuration options to BuildFunc.
type Config map[string]interface{}

// Decode loads a struct with the contents of the raw Config object.
func (c Config) Decode(v interface{}) error {
	return mapstructure.Decode(c, v)
}

// Builders returns all BuildFuncs for all Units provided by this package.
func Builders() map[string]BuildFunc {
	m := map[string]BuildFunc{}
	for k, v := range normalModules {
		m[k] = func(name string, f nameBuildFunc) BuildFunc {
			return func(cfg Config) (*Unit, error) {
				return f(name, cfg)
			}
		}(k, v)
	}
	for k, v := range binaryModules {
		m[k] = newBinary(k, v)
	}
	for k, v := range unaryModules {
		m[k] = newUnary(k, v)
	}
	return m
}

type nameBuildFunc func(string, Config) (*Unit, error)
