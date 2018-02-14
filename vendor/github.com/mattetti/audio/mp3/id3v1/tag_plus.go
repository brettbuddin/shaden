package id3v1

const (
	TagPlusSize = 227
)

var (
	TagPlusCode = []byte{84, 65, 71, 43} // "TAG+"
)

type TagPlugs struct {
	Title  [60]byte
	Artist [60]byte
	Album  [60]byte
	// 0=unset, 1=slow, 2= medium, 3=fast, 4=hardcore
	Speed uint8
	// A free-text field for the genre
	Genre [30]byte
	// the start of the music as mmm:ss
	StartTime [6]byte
	// the end of the music as mmm:ss
	EndTime [6]byte
}
