package randtest

import "math/rand"

// Static returns a deterministic *rand.Rand for use in tests.
func Static() *rand.Rand {
	return rand.New(rand.NewSource(1))
}
