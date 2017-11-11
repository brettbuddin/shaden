package musictheory

// Scale is a series of Transposers
type Scale []Transposer

// Transpose transposes a scale by the specified Interval
func (s Scale) Transpose(i Interval) Transposer {
	scale := Scale{}
	for _, transposer := range s {
		scale = append(scale, transposer.Transpose(i))
	}
	return scale
}

// NewScale returns a Scale built using a set of intervals
func NewScale(root Transposer, intervals []Interval, octaves int) Scale {
	scale := Scale{}
	for i := 0; i < octaves; i++ {
		for _, v := range intervals {
			scale = append(scale, root.Transpose(v))
		}
		root = root.Transpose(Octave(1))
	}
	return scale
}
