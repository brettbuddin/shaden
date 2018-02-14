// id3v1 is a package allowing the extraction of id3v1 tags.
// See http://en.wikipedia.org/wiki/ID3#ID3v1
package id3v1

var (
	HeaderTagID = []byte{0x54, 0x41, 0x47}
)

const (
	Size       = 128
	HeaderCode = "TAG"
)
