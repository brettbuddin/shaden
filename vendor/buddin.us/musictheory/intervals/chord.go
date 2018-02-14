package intervals

import mt "buddin.us/musictheory"

var (
	AugmentedMajorSeventh,
	AugmentedSeventh,
	AugmentedSixth,
	AugmentedTriad,
	DiminishedMajorSeventh,
	DiminishedSeventh,
	DiminishedTriad,
	DominantSeventh,
	HalfDiminishedSeventh,
	MajorSeventh,
	MajorSixth,
	MajorTriad,
	MinorMajorSeventh,
	MinorSeventh,
	MinorSixth,
	MinorTriad []mt.Interval
)

func init() {
	P1 := mt.Perfect(1)
	P5 := mt.Perfect(5)

	M3 := mt.Major(3)
	M6 := mt.Major(6)
	M7 := mt.Major(7)

	m3 := mt.Minor(3)
	m6 := mt.Minor(6)
	m7 := mt.Minor(7)

	A4 := mt.Augmented(4)
	A5 := mt.Augmented(5)
	d5 := mt.Diminished(5)
	d7 := mt.Diminished(7)

	MajorTriad = []mt.Interval{P1, M3, P5}
	MajorSeventh = []mt.Interval{P1, M3, P5, M7}
	MajorSixth = []mt.Interval{P1, M3, P5, M6}
	DominantSeventh = []mt.Interval{P1, M3, P5, m7}

	MinorTriad = []mt.Interval{P1, m3, P5}
	MinorMajorSeventh = []mt.Interval{P1, m3, P5, M7}
	MinorSeventh = []mt.Interval{P1, m3, P5, m7}
	MinorSixth = []mt.Interval{P1, m3, P5, m6}

	DiminishedTriad = []mt.Interval{P1, m3, d5}
	DiminishedMajorSeventh = []mt.Interval{P1, m3, d5, M7}
	DiminishedSeventh = []mt.Interval{P1, m3, d5, d7}
	HalfDiminishedSeventh = []mt.Interval{P1, m3, d5, m7}

	AugmentedTriad = []mt.Interval{P1, M3, A5}
	AugmentedMajorSeventh = []mt.Interval{P1, M3, A5, M7}
	AugmentedSeventh = []mt.Interval{P1, M3, A5, m7}
	AugmentedSixth = []mt.Interval{P1, A4, m6}
}
