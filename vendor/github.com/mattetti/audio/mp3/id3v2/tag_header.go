package id3v2

import "bytes"

// TagHeader is the header content for an id3v2 tag
type TagHeader [10]byte

// IsValidID checks that the tag header starts by the right ID
func (th TagHeader) IsValidID() bool {
	return (bytes.Compare(th[:3], HeaderTagID) == 0)
}

// ReadVersion extracts the version information.
func (th TagHeader) ReadVersion() Version {
	return Version{
		Major:    uint8(th[3]),
		Revision: uint8(th[4]),
	}
}

// ReadFlags reads the header flags
func (th TagHeader) ReadFlags() Flags {
	flags := Flags{}
	flags.Unsynchronisation = (th[5] & (1 << 0)) != 0     // 3.1.a
	flags.ExtendedHeader = (th[5] & (1 << 1)) != 0        // 3.1.b
	flags.ExperimentalIndicator = (th[5] & (1 << 2)) != 0 // 3.1.c
	flags.FooterPresent = (th[5] & (1 << 3)) != 0         // 3.1.d

	return flags
}

// ReadSize reads the size of the tag header
func (th TagHeader) ReadSize() (int, error) {
	return synchSafe(th[6:])
}
