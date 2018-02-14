package mp3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

type Frame struct {
	buf []byte
	// SkippedBytes is the amount of bytes we had to skip before getting to the frame
	SkippedBytes int
	// Counter gets incremented if the same frame is reused to parse a file
	Counter int
	Header  FrameHeader
}

type (
	// FrameVersion is the MPEG version given in the frame header
	FrameVersion byte
	// FrameLayer is the MPEG layer given in the frame header
	FrameLayer byte
	// FrameEmphasis is the Emphasis value from the frame header
	FrameEmphasis byte
	// FrameChannelMode is the Channel mode from the frame header
	FrameChannelMode byte
	// FrameBitRate is the bit rate from the frame header
	FrameBitRate int
	// FrameSampleRate is the sample rate from teh frame header
	FrameSampleRate int
	// FrameSideInfo holds the SideInfo bytes from the frame
	FrameSideInfo []byte
)

// Duration calculates the time duration of this frame based on the samplerate and number of samples
func (f *Frame) Duration() time.Duration {
	if !f.Header.IsValid() {
		return 0
	}
	ms := (1000 / float64(f.Header.SampleRate())) * float64(f.Header.Samples())
	dur := time.Duration(int(float64(time.Millisecond) * ms))
	if dur < 0 {
		// we have bad data, let's ignore it
		dur = 0
	}
	return dur
}

// CRC returns the CRC word stored in this frame
func (f *Frame) CRC() uint16 {
	var crc uint16
	if !f.Header.Protection() {
		return 0
	}
	crcdata := bytes.NewReader(f.buf[4:6])
	binary.Read(crcdata, binary.BigEndian, &crc)
	return crc
}

// SideInfo returns the  side info for this frame
func (f *Frame) SideInfo() FrameSideInfo {
	if f.Header.Protection() {
		return FrameSideInfo(f.buf[6:])
	} else {
		return FrameSideInfo(f.buf[4:])
	}
}

// Frame returns a string describing this frame, header and side info
func (f *Frame) String() string {
	str := ""
	str += fmt.Sprintf("Header: \n%s", f.Header)
	str += fmt.Sprintf("CRC: %x\n", f.CRC())
	str += fmt.Sprintf("Samples: %v\n", f.Header.Samples())
	str += fmt.Sprintf("Size: %v\n", f.Header.Size())
	str += fmt.Sprintf("Duration: %v\n", f.Duration())
	return str
}

// NDataBegin is the number of bytes before the frame header at which the sample data begins
// 0 indicates that the data begins after the side channel information. This data is the
// data from the "bit resevoir" and can be up to 511 bytes
func (i FrameSideInfo) NDataBegin() uint16 {
	return (uint16(i[0]) << 1 & (uint16(i[1]) >> 7))
}
