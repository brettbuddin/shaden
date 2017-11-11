package intervals

import mt "buddin.us/musictheory"

// Scales
var (
	Chromatic,
	Major,
	Minor,
	MajorPentatonic,
	MinorPentatonic,
	Ionian,
	Dorian,
	Phrygian,
	Aeolian,
	Lydian,
	Mixolydian,
	Locrian []mt.Interval
)

func init() {
	P1 := mt.Perfect(1)
	P4 := mt.Perfect(4)
	P5 := mt.Perfect(5)

	M2 := mt.Major(2)
	M3 := mt.Major(3)
	M6 := mt.Major(6)
	M7 := mt.Major(7)

	m2 := mt.Minor(2)
	m3 := mt.Minor(3)
	m6 := mt.Minor(6)
	m7 := mt.Minor(7)

	A4 := mt.Augmented(4)
	d5 := mt.Diminished(5)

	Chromatic = []mt.Interval{P1, m2, M2, m3, M3, P4, A4, P5, m6, M6, m7, M7}

	Ionian = []mt.Interval{P1, M2, M3, P4, P5, M6, M7}
	Major = Ionian

	Dorian = []mt.Interval{P1, M2, m3, P4, P5, M6, m7}
	Phrygian = []mt.Interval{P1, m2, m3, P4, P5, m6, m7}
	Lydian = []mt.Interval{P1, M2, M3, A4, P5, M6, M7}
	Mixolydian = []mt.Interval{P1, M2, M3, P4, P5, M6, m7}

	Aeolian = []mt.Interval{P1, M2, m3, P4, P5, m6, m7}
	Minor = Aeolian

	MajorPentatonic = []mt.Interval{P1, M2, M3, P5, M6}
	MinorPentatonic = []mt.Interval{P1, m3, P4, P5, m7}

	Locrian = []mt.Interval{P1, m2, m3, P4, d5, m6, m7}
}
