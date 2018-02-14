package midi

import (
	"encoding/binary"
	"fmt"
	"io"
)

type timeFormat int
type nextChunkType int

const (
	MetricalTF timeFormat = iota + 1
	TimeCodeTF
)

const (
	eventChunk nextChunkType = iota + 1
	trackChunk
)

/*
  Decoder
 Format documented there: http://www.music.mcgill.ca/~ich/classes/mumt306/midiformat.pdf
 <Header Chunk> = <chunk type><length><format><ntrks><division>


				Division, specifies the meaning of the delta-times.
				It has two formats, one for metrical time, and one for time-code-based
				time:
				 +---+-----------------------------------------+
				 | 0 | ticks per quarter-note                  |
				 ==============================================|
				 | 1 | negative SMPTE format  | ticks per frame|
				 +---+-----------------------+-----------------+
				 |15 |14                    8 |7             0 |
				If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
				of delta time "ticks" which make up a quarter-note. For instance, if
				division is 96, then a time interval of an eighth-note between two
				events in the file would be 48.
		    If bit 15 of <division> is a one, delta times in a file correspond
		    to subdivisions of a second, in a way consistent with SMPTE and MIDI
		    Time Code. Bits 14 thru 8 contain one of the four values -24, -25, -29,
		    or -30, corresponding to the four standard SMPTE and MIDI Time Code
		    formats (-29 corresponds to 30 drop frome), and represents the
		    number of frames per second. These negative numbers are stored in
		    two's compliment form. The second byte (stored positive) is the
		    resolution within a frame: typical values may be 4 (MIDI Time Code
		    resolution), 8, 10, 80 (bit resolution), or 100. This stream allows
		    exact specifications of time-code-based tracks, but also allows
		    milisecond-based tracks by specifying 25|frames/sec and a
		    resolution of 40 units per frame. If the events in a file are stored
		    with a bit resolution of thirty-framel time code, the division word
		    would be E250 hex.

*/
type Decoder struct {
	r io.Reader

	Ch chan *Track
	/*
	   Format describes the tracks format

	   0	-	single-track
	   Format 0 file has a header chunk followed by one track chunk. It
	   is the most interchangable representation of data. It is very useful
	   for a simple single-track player in a program which needs to make
	   synthesizers make sounds, but which is primarily concerened with
	   something else such as mixers or sound effect boxes. It is very
	   desirable to be able to produce such a format, even if your program
	   is track-based, in order to work with these simple programs. On the
	   other hand, perhaps someone will write a format conversion from
	   format 1 to format 0 which might be so easy to use in some setting
	   that it would save you the trouble of putting it into your program.


	   Synchronous multiple tracks means that the tracks will all be vertically synchronous, or in other words,
	    they all start at the same time, and so can represent different parts in one song.
	    1	-	multiple tracks, synchronous
	    Asynchronous multiple tracks do not necessarily start at the same time, and can be completely asynchronous.
	    2	-	multiple tracks, asynchronous
	*/
	Format uint16

	// NumTracks represents the number of tracks in the midi file
	NumTracks uint16

	TicksPerQuarterNote uint16

	TimeFormat timeFormat
	Tracks     []*Track
}

func (d *Decoder) CurrentTrack() *Track {
	if d == nil || len(d.Tracks) == 0 {
		return nil
	}
	return d.Tracks[len(d.Tracks)-1]
}

func (d *Decoder) Parse() error {
	var err error
	var code [4]byte

	if err := binary.Read(d.r, binary.BigEndian, &code); err != nil {
		return err
	}
	if code != headerChunkID {
		return fmt.Errorf("%s - %s", ErrFmtNotSupported, code)
	}
	var headerSize uint32
	if err := binary.Read(d.r, binary.BigEndian, &headerSize); err != nil {
		return err
	}

	if headerSize != 6 {
		return fmt.Errorf("%s - expected header size to be 6, was %d", ErrFmtNotSupported, headerSize)
	}

	if err := binary.Read(d.r, binary.BigEndian, &d.Format); err != nil {
		return err
	}

	if err := binary.Read(d.r, binary.BigEndian, &d.NumTracks); err != nil {
		return err
	}

	var division uint16
	if err := binary.Read(d.r, binary.BigEndian, &division); err != nil {
		return err
	}

	// If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note. For instance, if
	// division is 96, then a time interval of an eighth-note between two
	// events in the file would be 48.
	if (division & 0x8000) == 0 {
		d.TicksPerQuarterNote = division & 0x7FFF
		d.TimeFormat = MetricalTF
	} else {
		/*
			If bit 15 of <division> is a one, delta times in a file correspond
			to subdivisions of a second, in a way consistent with SMPTE and MIDI
			Time Code. Bits 14 thru 8 contain one of the four values -24, -25, -29,
			or -30, corresponding to the four standard SMPTE and MIDI Time Code
			formats (-29 corresponds to 30 drop frome), and represents the
			number of frames per second. These negative numbers are stored in
			two's compliment form. The second byte (stored positive) is the
			resolution within a frame: typical values may be 4 (MIDI Time Code
			resolution), 8, 10, 80 (bit resolution), or 100. This stream allows
			exact specifications of time-code-based tracks, but also allows
			milisecond-based tracks by specifying 25|frames/sec and a
			resolution of 40 units per frame. If the events in a file are stored
			with a bit resolution of thirty-framel time code, the division word
			would be E250 hex.
		*/
		d.TimeFormat = TimeCodeTF
	}

	_, nextChunk, err := d.parseTrack()
	if err != nil {
		return err
	}

	for err != io.EOF {
		switch nextChunk {
		case eventChunk:
			nextChunk, err = d.parseEvent()
		case trackChunk:
			_, nextChunk, err = d.parseTrack()
		}

		if err != nil && err != io.EOF {
			return err
		}
	}

	// All done
	return nil
}

func (d *Decoder) parseTrack() (uint32, nextChunkType, error) {
	id, size, err := d.IDnSize()
	if err != nil {
		return 0, trackChunk, err
	}
	if id != trackChunkID {
		return 0, trackChunk, fmt.Errorf("%s - Expected track chunk ID %s, got %s", ErrUnexpectedData, trackChunkID, id)
	}
	d.Tracks = append(d.Tracks, &Track{Size: size})
	return size, eventChunk, nil
}

// IDnSize returns the next ID + block size
func (d *Decoder) IDnSize() ([4]byte, uint32, error) {
	var ID [4]byte
	var blockSize uint32
	if err := binary.Read(d.r, binary.BigEndian, &ID); err != nil {
		return ID, blockSize, err
	}
	if err := binary.Read(d.r, binary.BigEndian, &blockSize); err != err {
		return ID, blockSize, err
	}
	return ID, blockSize, nil
}

// VarLen returns the variable length value at the exact parser location.
func (d *Decoder) VarLen() (val uint32, readBytes uint32, err error) {
	buf := []byte{}
	var lastByte bool
	var n uint32

	for !lastByte {
		b, err := d.ReadByte()
		if err != nil {
			return 0, n, err
		}
		buf = append(buf, b)
		lastByte = (b>>7 == 0x0)
		n++
	}

	val, nUsed := DecodeVarint(buf)
	return val, uint32(nUsed), nil
}

// VarLengthTxt Returns a variable length text string
// as well as the amount of bytes read
func (d *Decoder) VarLenTxt() (string, uint32, error) {
	var l uint32
	var err error
	var n uint32

	if l, n, err = d.VarLen(); err != nil {
		return "", n, err
	}
	buf := make([]byte, l)
	err = d.Read(buf)
	return string(buf), n, err
}

func (d *Decoder) ReadByte() (byte, error) {
	var b byte
	err := binary.Read(d.r, binary.BigEndian, &b)
	return b, err
}

// read reads n bytes from the parser's reader and stores them into the provided dst,
// which must be a pointer to a fixed-size value.
func (d *Decoder) Read(dst interface{}) error {
	return binary.Read(d.r, binary.BigEndian, dst)
}

// Uint7 reads a byte and converts the first 7 bits into an uint8
func (d *Decoder) Uint7() (uint8, error) {
	b, err := d.ReadByte()
	if err != nil {
		return 0, err
	}
	return (b & 0x7f), nil
}

// Uint24 reads 3 bytes and convert them into a uint32
func (d *Decoder) Uint24() (uint32, error) {
	bytes := make([]byte, 3)
	if err := d.Read(bytes); err != nil {
		return 0, err
	}

	var output uint32
	output |= uint32(bytes[2]) << 0
	output |= uint32(bytes[1]) << 8
	output |= uint32(bytes[0]) << 16

	return output, nil
}
