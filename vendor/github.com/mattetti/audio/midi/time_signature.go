package midi

import (
	"fmt"
	"math"
)

// Time signature
// FF 58 04 nn dd cc bb Time Signature
// The time signature is expressed as four numbers. nn and dd
// represent the numerator and denominator of the time signature as it
// would be notated. The denominator is a negative power of two: 2
// represents a quarter-note, 3 represents an eighth-note, etc.
// The cc parameter expresses the number of MIDI clocks in a
// metronome click. The bb parameter expresses the number of
// notated 32nd-notes in a MIDI quarter-note (24 MIDI clocks). This
// was added because there are already multiple programs which allow a
// user to specify that what MIDI thinks of as a quarter-note (24 clocks)
// is to be notated as, or related to in terms of, something else.
type TimeSignature struct {
	Numerator                   uint8
	Denominator                 uint8
	ClocksPerTick               uint8
	ThirtySecondNotesPerQuarter uint8
}

// Denum returns the notation denominator (which is not how it's stored in MIDI)
func (ts *TimeSignature) Denum() int {
	return int(math.Exp2(float64(ts.Denominator)))
}

func (ts *TimeSignature) String() string {
	return fmt.Sprintf("%d/%d - %d clocks per tick - %d", ts.Numerator, ts.Denum(), ts.ClocksPerTick, ts.ThirtySecondNotesPerQuarter)
}
