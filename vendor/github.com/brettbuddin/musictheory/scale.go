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
	scale := Scale{}
	for i := 0; i < octaves; i++ {
		for _, v := range intervals {
			scale = append(scale, root.Transpose(v))
		}
		root = root.Transpose(Octave(1))
	}
	return scale
}
