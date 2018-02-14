package id3v1

const (
	// TagSize is the size in bytes of an id3v1 tag
	TagSize = 128
)

var (
	//TagCode is the byte representation of "TAG"
	TagCode = []byte{84, 65, 71}
)

// Tag contains the various information stored in a id3 v1 tag.
// Strings are either space or zero-padded.
// Unset string entries are filled using an empty string. ID3v1 is 128 bytes long
type Tag struct {
	Title  [30]byte
	Artist [30]byte
	Album  [30]byte
	Year   [4]byte
	// The track number is stored in the last two bytes of the comment field.
	// If the comment is 29 or 30 characters long, no track number can be stored.
	Comment  [30]byte
	ZeroByte byte
	Track    byte
	Genre    byte
}
