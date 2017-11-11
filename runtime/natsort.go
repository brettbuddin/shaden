package runtime

import (
	"regexp"
	"sort"
	"strconv"
)

var natsortChunk = regexp.MustCompile(`(\d+|\D+)`)

func natsort(s []string) {
	sort.Slice(s, func(i, j int) bool {
		return natsortCompare(s[i], s[j])
	})
}

func natsortCompare(a, b string) bool {
	var (
		chunksA  = natsortChunk.FindAllString(a, -1)
		chunksB  = natsortChunk.FindAllString(b, -1)
		nChunksA = len(chunksA)
		nChunksB = len(chunksB)
	)

	for i := range chunksA {
		aInt, aErr := strconv.Atoi(chunksA[i])
		bInt, bErr := strconv.Atoi(chunksB[i])

		if aErr == nil && bErr == nil {
			if aInt == bInt {
				if i == nChunksA-1 {
					return true
				} else if i == nChunksB-1 {
					return false
				}
				continue
			}
			return aInt < bInt
		}

		if chunksA[i] == chunksB[i] {
			if i == nChunksA-1 {
				return true
			} else if i == nChunksB-1 {
				return false
			}
			continue
		}
		return chunksA[i] < chunksB[i]
	}

	return false
}
