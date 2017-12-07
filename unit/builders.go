package unit

import "github.com/mitchellh/mapstructure"

var (
	builders = map[string]nameBuildFunc{
		"abs":      unaryBuildFunc(unaryAbs),
		"bipolar":  unaryBuildFunc(unaryBipolar),
		"ceil":     unaryBuildFunc(unaryCeil),
		"floor":    unaryBuildFunc(unaryFloor),
		"invert":   unaryBuildFunc(unaryInv),
		"not":      unaryBuildFunc(unaryNOT),
		"noop":     unaryBuildFunc(unaryNoop),
		"unipolar": unaryBuildFunc(unaryUnipolar),

		"and":   binaryBuildFunc(binaryAND),
		"diff":  binaryBuildFunc(binaryDiff),
		"div":   binaryBuildFunc(binaryDiv),
		"gt":    binaryBuildFunc(binaryGT),
		"imply": binaryBuildFunc(binaryIMPLY),
		"lt":    binaryBuildFunc(binaryLT),
		"max":   binaryBuildFunc(binaryMax),
		"min":   binaryBuildFunc(binaryMin),
		"mod":   binaryBuildFunc(binaryMod),
		"mult":  binaryBuildFunc(binaryMult),
		"nand":  binaryBuildFunc(binaryNAND),
		"nor":   binaryBuildFunc(binaryNOR),
		"or":    binaryBuildFunc(binaryOR),
		"sum":   binaryBuildFunc(binarySum),
		"xnor":  binaryBuildFunc(binaryXNOR),
		"xor":   binaryBuildFunc(binaryXOR),

		"adjust":             newAdjust,
		"adsr":               newAdsr,
		"center":             newCenter,
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
		"lag":                newLag,
		"latch":              newLatch,
		"lerp":               newInterpolate,
		"logic":              newLogic,
		"low-gen":            newLowGen,
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
		"sample":             newWAVSample,
		"shift":              newShift,
		"slope":              newSlope,
		"smooth":             newSmooth,
		"stages":             newStages,
		"switch":             newSwitch,
		"toggle":             newToggle,
		"transpose":          newTranspose,
		"transpose-interval": newTransposeInterval,
		"val-gate":           newValToGate,
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
	for k, v := range builders {
		m[k] = func(name string, f nameBuildFunc) BuildFunc {
			return func(cfg Config) (*Unit, error) {
				return f(name, cfg)
			}
		}(k, v)
	}
	return m
}

type nameBuildFunc func(string, Config) (*Unit, error)

func unaryBuildFunc(op unaryOp) nameBuildFunc {
	return func(name string, c Config) (*Unit, error) {
		return newUnary(name, op)(c)
	}
}

func binaryBuildFunc(op binaryOp) nameBuildFunc {
	return func(name string, c Config) (*Unit, error) {
		return newBinary(name, op)(c)
	}
}
