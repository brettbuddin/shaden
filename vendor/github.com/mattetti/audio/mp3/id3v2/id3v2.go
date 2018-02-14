package id3v2

import (
	"errors"
	"fmt"
)

var (
	// HeaderTagID are the 3 bytes starting the ID3 v2 tag
	HeaderTagID = []byte{0x49, 0x44, 0x33}

	// ErrInvalidTagHeader
	ErrInvalidTagHeader = errors.New("invalid tag header")
)

type Header struct {
	Version Version
	Flags   Flags
	Size    int
}

type Version struct {
	Major    uint8
	Revision uint8
}

type Flags struct {
	Unsynchronisation     bool
	ExtendedHeader        bool
	ExperimentalIndicator bool
	FooterPresent         bool
}

type Frame struct {
	Header [10]byte
	Data   []byte
}

type Tag struct {
	Header         *Header
	extendedHeader []byte
	frameSets      map[string][]*Frame
}

// ReadHeader reads the 10 bytes header and parses the data which gets stored in
// the tag header.
func (t *Tag) ReadHeader(th TagHeader) error {
	// id3 tag ID
	if !th.IsValidID() {
		return ErrInvalidTagHeader
	}
	// version
	t.Header = &Header{
		Version: th.ReadVersion(),
		Flags:   th.ReadFlags(),
	}
	size, err := th.ReadSize()
	if err != nil {
		return err
	}
	t.Header.Size = size

	return nil
}

// synchsafe integers
// https://en.wikipedia.org/wiki/Synchsafe
func synchSafe(buf []byte) (int, error) {
	n := int(0)
	for _, b := range buf {
		if (b & (1 << 7)) != 0 {
			return 0, fmt.Errorf("invalid synchsafe integer")
		}
		n |= (n << 7) | int(b)
	}
	return n, nil
}
