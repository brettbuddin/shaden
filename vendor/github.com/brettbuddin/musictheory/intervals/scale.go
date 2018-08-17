package intervals

import mt "github.com/brettbuddin/musictheory"

// Scales
var (
	Aeolian,
	Chromatic,
	DominantBebop,
	Dorian,
	DoubleHarmonic,
	HarmonicMinor,
	HarmonicMinorBebop,
	InSen,
	Ionian,
	Lydian,
	Major,
	MajorBebop,
	MajorPentatonic,
	MelodicMinorBebop,
	Minor,
	MinorPentatonic,
	Mixolydian,
	Phrygian,
	Locrian,
	WholeTone []mt.Interval
)

func init() {
	var (
		A4 = mt.Augmented(4)
		A5 = mt.Augmented(5)
		M2 = mt.Major(2)
		M3 = mt.Major(3)
		M6 = mt.Major(6)
		M7 = mt.Major(7)
		P1 = mt.Perfect(1)
		P4 = mt.Perfect(4)
		P5 = mt.Perfect(5)
		d5 = mt.Diminished(5)
		d7 = mt.Diminished(7)
		m2 = mt.Minor(2)
		m3 = mt.Minor(3)
		m6 = mt.Minor(6)
		m7 = mt.Minor(7)
	)

	Aeolian = []mt.Interval{P1, M2, m3, P4, P5, m6, m7}
	Chromatic = []mt.Interval{P1, m2, M2, m3, M3, P4, A4, P5, m6, M6, m7, M7}
	DominantBebop = []mt.Interval{P1, M2, M3, P4, P5, M6, m7, M7}
	Dorian = []mt.Interval{P1, M2, m3, P4, P5, M6, m7}
	DoubleHarmonic = []mt.Interval{m2, M3, P4, P5, m6, M7}
	HarmonicMinor = []mt.Interval{P1, M2, m3, P4, P5, m6, M7}
	HarmonicMinorBebop = []mt.Interval{P1, M2, m3, P4, P5, M6, d7, m7}
	InSen = []mt.Interval{P1, m2, P4, P5, m7}
	Ionian = []mt.Interval{P1, M2, M3, P4, P5, M6, M7}
	Locrian = []mt.Interval{P1, m2, m3, P4, d5, m6, m7}
	Lydian = []mt.Interval{P1, M2, M3, A4, P5, M6, M7}
	MajorBebop = []mt.Interval{P1, M2, M3, P4, P5, A5, M6, M7}
	MajorPentatonic = []mt.Interval{P1, M2, M3, P5, M6}
	MelodicMinorBebop = []mt.Interval{P1, M2, m3, P4, P5, m6, M6, M7}
	MinorPentatonic = []mt.Interval{P1, m3, P4, P5, m7}
	Mixolydian = []mt.Interval{P1, M2, M3, P4, P5, M6, m7}
	Phrygian = []mt.Interval{P1, m2, m3, P4, P5, m6, m7}
	WholeTone = []mt.Interval{P1, M2, M3, d5, m6, m7}

	Major = Ionian
	Minor = Aeolian
}
