package musictheory

import (
	"fmt"
	"math"
)

// Quality types
const (
	PerfectType QualityType = iota
	MajorType
	MinorType
	AugmentedType
	DiminishedType
)

// IntervalFunc creates an interval at as specific step/degree
type IntervalFunc func(int) Interval

// Perfect interval
func Perfect(step int) Interval {
	return qualityInterval(step, Quality{PerfectType, 0})
}

// Major interval
func Major(step int) Interval {
	return qualityInterval(step, Quality{MajorType, 0})
}

// Minor interval
func Minor(step int) Interval {
	return qualityInterval(step, Quality{MinorType, 0})
}

// Augmented interval
func Augmented(step int) Interval {
	return qualityInterval(step, Quality{AugmentedType, 1})
}

// DoublyAugmented interval
func DoublyAugmented(step int) Interval {
	return qualityInterval(step, Quality{AugmentedType, 2})
}

// Diminished interval
func Diminished(step int) Interval {
	return qualityInterval(step, Quality{DiminishedType, 1})
}

// DoublyDiminished interval
func DoublyDiminished(step int) Interval {
	return qualityInterval(step, Quality{DiminishedType, 2})
}

// Octave interval
func Octave(step int) Interval {
	return Interval{step, 0, 0}
}

// Semitones is an interval using direct semitones
func Semitones(step int) Interval {
	return Interval{chromaticOctaves(step), chromaticToDiatonic(step), step}
}

func qualityInterval(step int, quality Quality) Interval {
	absStep := int(math.Abs(float64(step)))
	diatonic := normalizeDiatonic(absStep - 1)
	diff := qualityDiff(quality, canBePerfect(diatonic))
	octaves := diatonicOctaves(absStep - 1)

	i := NewInterval(absStep, octaves, diff)
	if step > 0 {
		return i
	}
	return i.Negate()
}

// NewInterval builds a new Interval
func NewInterval(step, octaves, offset int) Interval {
	diatonic := normalizeDiatonic(step - 1)
	chromatic := diatonicToChromatic(diatonic) + offset

	return Interval{octaves, diatonic, chromatic}
}

// Interval represents an interval in 12-tone equal temperament
type Interval struct {
	Octaves   int
	Diatonic  int
	Chromatic int
}

func (i Interval) String() string {
	return fmt.Sprintf("(octaves: %d, diatonic: %d, chromatic: %d)", i.Octaves, i.Diatonic, i.Chromatic)
}

// Semitones returns the total number of semitones that make up the interval
func (i Interval) Semitones() int {
	return i.Octaves*12 + i.Chromatic
}

// Quality returns the Quality
func (i Interval) Quality() Quality {
	quality := diffQuality(i.Chromatic-diatonicToChromatic(i.Diatonic), canBePerfect(i.Diatonic))

	if i.Octaves < 0 {
		return quality.Invert()
	}

	return quality
}

// Ratio returns the interval ratio
func (i Interval) Ratio() float64 {
	return math.Exp2(float64(i.Semitones()) / 12.0)
}

// Transpose returns a new Interval that has been transposed by the given Interval
func (i Interval) Transpose(o Interval) Interval {
	var diatonic int

	// TODO: Accomodate weird behavior of sequential minor second transpositions. We don't need to advance the diatonic
	// every transposition. We're currently modeling things as integers, but maybe we need to model as floats and
	// accumulate over time; whole numbers trigger a move.
	if o.Diatonic == o.Chromatic {
		if diatonicToChromatic(i.Diatonic) == i.Chromatic {
			diatonic = i.Diatonic + o.Diatonic
		} else {
			diatonic = i.Diatonic
		}
	} else {
		diatonic = i.Diatonic + o.Diatonic
	}

	diatonicOctaves := diatonicOctaves(diatonic)
	diatonicRemainder := normalizeDiatonic(diatonic)

	octaves := i.Octaves + o.Octaves + diatonicOctaves
	chromatic := normalizeChromatic(i.Chromatic + o.Chromatic)

	return Interval{octaves, diatonicRemainder, chromatic}
}

// Negate returns a new, negated Interval
func (i Interval) Negate() Interval {
	if i.Diatonic == 0 && i.Chromatic == 0 {
		return Interval{-i.Octaves, i.Diatonic, i.Chromatic}
	}

	return Interval{-(i.Octaves + 1), inverseDiatonic(i.Diatonic), inverseChromatic(i.Chromatic)}
}

// Eq determines if another interval is the same
func (i Interval) Eq(o Interval) bool {
	return i.Semitones() == o.Semitones()
}

// QualityType represents the type a Quality can take
type QualityType int

func (q QualityType) String() string {
	switch q {
	case PerfectType:
		return "perfect"
	case MajorType:
		return "major"
	case MinorType:
		return "minor"
	case AugmentedType:
		return "augmented"
	case DiminishedType:
		return "diminished"
	default:
		return "unknown"
	}
}

// Quality describes the quality of an interval
type Quality struct {
	Type QualityType
	Size int
}

// Invert returns a new, inverted Quality
func (q Quality) Invert() Quality {
	switch q.Type {
	case PerfectType:
		return q
	case MajorType:
		return Quality{MinorType, q.Size}
	case MinorType:
		return Quality{MajorType, q.Size}
	case AugmentedType:
		return Quality{DiminishedType, q.Size}
	case DiminishedType:
		return Quality{AugmentedType, q.Size}
	default:
		panic(fmt.Sprintf("invalid type: %d", q.Type))
	}
}

// Eq checks two Qualities for equality
func (q Quality) Eq(o Quality) bool {
	return q.Type == o.Type && q.Size == o.Size
}

func (q Quality) String() string {
	switch q.Type {
	case PerfectType, MajorType, MinorType:
		return fmt.Sprintf("%s", q.Type)
	case AugmentedType, DiminishedType:
		return fmt.Sprintf("%s(%d)", q.Type, q.Size)
	default:
		return "unknown"
	}
}

func diatonicToChromatic(interval int) int {
	if interval >= len(diatonicToChromaticLookup) {
		panic(fmt.Sprintf("interval out of range: %d", interval))
	}

	return diatonicToChromaticLookup[interval]
}

var diatonicToChromaticLookup = []int{0, 2, 4, 5, 7, 9, 11}

func chromaticToDiatonic(v int) int {
	mag := 1
	if v < 0 {
		mag = -1
		v = -v
	}
	v = normalizeChromatic(v)

	for i, c := range diatonicToChromaticLookup {
		if v == c || v < c {
			return i * mag
		}
	}
	return 6 * mag
}

func qualityDiff(q Quality, perfect bool) int {
	if q.Type == PerfectType || q.Type == MajorType {
		return 0
	} else if q.Type == MinorType {
		return -1
	} else if q.Type == AugmentedType {
		return q.Size
	} else if q.Type == DiminishedType {
		if perfect {
			return -q.Size
		}
		return -(q.Size + 1)
	}
	panic("invalid quality")
}

func diffQuality(diff int, perfect bool) Quality {
	if perfect {
		if diff == 0 {
			return Quality{PerfectType, 0}
		} else if diff > 0 {
			return Quality{AugmentedType, diff}
		}

		return Quality{DiminishedType, -diff}
	}

	if diff == 0 {
		return Quality{MajorType, 0}
	} else if diff == -1 {
		return Quality{MinorType, 0}
	} else if diff > 0 {
		return Quality{AugmentedType, diff}
	}

	return Quality{DiminishedType, -(diff + 1)}
}

func canBePerfect(interval int) bool {
	return interval == 0 || interval == 3 || interval == 4
}
