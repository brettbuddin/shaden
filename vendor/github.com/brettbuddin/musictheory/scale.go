package musictheory

// Scale is a series of Pitches
type Scale []Pitch

// Transpose transposes a scale by the specified Interval
func (s Scale) Transpose(i Interval) Scale {
	scale := Scale{}
	for _, transposer := range s {
		scale = append(scale, transposer.Transpose(i))
	}
	return scale
}

// NewScale returns a Scale built using a set of intervals
func NewScale(root Pitch, intervals []Interval, octaves int) Scale {
	var (
		scale        = Scale{}
		originalRoot = root
		descending   = (octaves < 0)
	)

	// Begin at the base of our octave shift
	if descending {
		root = root.Transpose(Octave(octaves))
	}

	for i := 0; i < abs(octaves); i++ {
		for j, v := range intervals {
			// Ignore the tonic which will become the *last* item in the slice
			// once reversed. This is to maintain consistency with ascending
			// scales: they don't include the final octave of the tonic.
			if descending && i == 0 && j == 0 {
				continue
			}
			scale = append(scale, root.Transpose(v))
		}
		root = root.Transpose(Octave(1))
	}

	// Add the original tonic to the end. It's about to become the beginning of
	// the slice once it's reversed. Reversing the list produces our descending
	// scale.
	if descending {
		scale = append(scale, originalRoot)
		for i := len(scale)/2 - 1; i >= 0; i-- {
			opp := len(scale) - 1 - i
			scale[i], scale[opp] = scale[opp], scale[i]
		}
	}

	return scale
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
