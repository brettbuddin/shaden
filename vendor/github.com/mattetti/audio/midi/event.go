package midi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

// NoteOn returns a pointer to a new event of type NoteOn (without the delta timing data)
func NoteOn(channel, key, vel int) *Event {
	return &Event{
		MsgChan:  uint8(channel),
		MsgType:  uint8(EventByteMap["NoteOn"]),
		Note:     uint8(key),
		Velocity: uint8(vel),
	}
}

// NoteOff returns a pointer to a new event of type NoteOff (without the delta timing data)
func NoteOff(channel, key int) *Event {
	return &Event{
		MsgChan:  uint8(channel),
		MsgType:  uint8(EventByteMap["NoteOff"]),
		Note:     uint8(key),
		Velocity: 64,
	}
}

// EndOfTrack indicates the end of the midi track. Note that this event is
// automatically added when encoding a normal track.
func EndOfTrack() *Event {
	return &Event{
		MsgType: uint8(EventByteMap["Meta"]),
		MsgChan: uint8(15),
		Cmd:     MetaByteMap["End of Track"],
	}
}

// AfterTouch returns a pointer to a new aftertouch event
func Aftertouch(channel, key, vel int) *Event {
	return &Event{
		MsgChan:  uint8(channel),
		MsgType:  uint8(EventByteMap["AfterTouch"]),
		Note:     uint8(key),
		Velocity: uint8(vel),
	}
}

// ControlChange sets a new value for a given controller
// The controller number is between 0-119.
// The new controller value is between 0-127.
func ControlChange(channel, controller, newVal int) *Event {
	return &Event{
		MsgChan:    uint8(channel),
		MsgType:    uint8(EventByteMap["ControlChange"]),
		Controller: uint8(controller),
		NewValue:   uint8(newVal),
	}
}

// ProgramChange sets a new value the same way as ControlChange
// but implements Mode control and special message by using reserved controller numbers 120-127.
func ProgramChange(channel, controller, newVal int) *Event {
	return &Event{
		MsgChan:    uint8(channel),
		MsgType:    uint8(EventByteMap["ProgramChange"]),
		Controller: uint8(controller),
		NewValue:   uint8(newVal),
	}
}

// ChannelAfterTouch is a global aftertouch with a value from 0 to 127
func ChannelAfterTouch(channel, vel int) *Event {
	return &Event{
		MsgChan:  uint8(channel),
		MsgType:  uint8(EventByteMap["ChannelAfterTouch"]),
		Pressure: uint8(vel),
	}
}

// PitchWheelChange is sent to indicate a change in the pitch bender.
// The possible value goes between 0 and 16383 where 8192 is the center.
func PitchWheelChange(channel, int, val int) *Event {
	return &Event{
		MsgChan:      uint8(channel),
		MsgType:      uint8(EventByteMap["PitchWheelChange"]),
		AbsPitchBend: uint16(val),
	}
}

// CopyrightEvent returns a copyright event with the passed string in it.
func CopyrightEvent(txt string) *Event {
	return &Event{
		MsgType:   uint8(EventByteMap["Meta"]),
		Cmd:       uint8(MetaByteMap["Copyright"]),
		Copyright: txt,
	}
}

// TempoEvent returns a new tempo event of the passed value.
func TempoEvent(bpmF float64) *Event {
	ms := uint32(60000000 / bpmF)

	// bpm is expressed in in microseconds per MIDI quarter-note
	// This event indicates a tempo change.  Another way of putting
	// "microseconds per quarter-note" is "24ths of a microsecond per MIDI
	// clock".  Representing tempos as time per beat instead of beat per time
	// allows absolutely exact dword-term synchronization with a time-based sync
	// protocol such as SMPTE time code or MIDI time code.  This amount of
	// accuracy provided by this tempo resolution allows a four-minute piece at
	// 120 beats per minute to be accurate within 500 usec at the end of the
	// piece.  Ideally, these events should only occur where MIDI clocks would
	// be located Q this convention is intended to guarantee, or at least
	// increase the likelihood, of compatibility with other synchronization
	// devices so that a time signature/tempo map stored in this format may
	// easily be transferred to another device.
	return &Event{
		MsgType:        uint8(EventByteMap["Meta"]),
		Cmd:            uint8(MetaByteMap["Tempo"]),
		MsPerQuartNote: ms,
	}
}

// TODO
func Meta(channel int) *Event {
	return nil
}

// Event
// <event> = <MIDI event> | <sysex event> | <meta-event>
// <MIDI event> is any MIDI channel message.
// Running status is used:
// status bytes of MIDI channel messages may be omitted if the preceding
// event is a MIDI channel message with the same status. The first event
// in each MTrk chunk must specifyy status. Delta-time is not
// considered an event itself: it is an integral part of the syntax for
// an MTrk event. Notice that running status occurs across delta-times.
// See http://www.indiana.edu/~emusic/etext/MIDI/chapter3_MIDI4.shtml
type Event struct {
	TimeDelta    uint32
	MsgType      uint8
	MsgChan      uint8
	Note         uint8
	Velocity     uint8
	Pressure     uint8
	Controller   uint8
	NewValue     uint8
	NewProgram   uint8
	Channel      uint8
	AbsPitchBend uint16
	RelPitchBend int16
	// Meta
	Cmd            uint8
	SeqNum         uint16
	Text           string
	Copyright      string
	SeqTrackName   string
	InstrumentName string
	Lyric          string
	Marker         string
	CuePoint       string
	MsPerQuartNote uint32
	Bpm            uint32
	TimeSignature  *TimeSignature
	// A positive value for the key specifies the number of sharps and a negative value specifies the number of flats.
	Key int32 //-7 to +7
	// A value of 0 for the scale specifies a major key and a value of 1 specifies a minor key.
	Scale uint32 // 0 or 1
	//
	SmpteOffset *SmpteOffset
}

// Copy returns an exact copy of the event
func (e *Event) Copy() *Event {
	newEv := &Event{
		TimeDelta:    e.TimeDelta,
		MsgType:      e.MsgType,
		MsgChan:      e.MsgChan,
		Note:         e.Note,
		Velocity:     e.Velocity,
		Pressure:     e.Pressure,
		Controller:   e.Controller,
		NewValue:     e.NewValue,
		NewProgram:   e.NewProgram,
		Channel:      e.Channel,
		AbsPitchBend: e.AbsPitchBend,
		RelPitchBend: e.RelPitchBend,
		// Meta
		Cmd:            e.Cmd,
		SeqNum:         e.SeqNum,
		Text:           e.Text,
		Copyright:      e.Copyright,
		SeqTrackName:   e.SeqTrackName,
		InstrumentName: e.InstrumentName,
		Lyric:          e.Lyric,
		Marker:         e.Marker,
		CuePoint:       e.CuePoint,
		MsPerQuartNote: e.MsPerQuartNote,
		Bpm:            e.Bpm,
		Key:            e.Key,
		Scale:          e.Scale,
	}
	if e.TimeSignature != nil {
		newEv.TimeSignature = &TimeSignature{
			Numerator:                   e.TimeSignature.Numerator,
			Denominator:                 e.TimeSignature.Denominator,
			ClocksPerTick:               e.TimeSignature.ClocksPerTick,
			ThirtySecondNotesPerQuarter: e.TimeSignature.ThirtySecondNotesPerQuarter,
		}
	}
	if e.SmpteOffset != nil {
		newEv.SmpteOffset = &SmpteOffset{
			Hour:  e.SmpteOffset.Hour,
			Min:   e.SmpteOffset.Min,
			Sec:   e.SmpteOffset.Sec,
			Fr:    e.SmpteOffset.Fr,
			SubFr: e.SmpteOffset.SubFr,
		}
	}
	return newEv
}

// String implements the stringer interface
func (e *Event) String() string {
	if e == nil {
		return ""
	}
	var k string
	var ok bool
	if k, ok = EventMap[e.MsgType]; !ok {
		k = fmt.Sprintf("%#X", e.MsgType)
	}
	out := fmt.Sprintf("Ch %d @ %d \t%s", e.MsgChan, e.TimeDelta, k)
	if e.Velocity > 0 {
		out += fmt.Sprintf(" Vel: %d", e.Velocity)
	}
	if e.MsgType == EventByteMap["NoteOn"] || e.MsgType == EventByteMap["NoteOff"] {
		out += fmt.Sprintf(" Note: %s", NoteToName(int(e.Note)))
	}
	if e.Cmd != 0 {
		out = fmt.Sprintf("Ch %d @ %d \t%s", e.MsgChan, e.TimeDelta, MetaCmdMap[e.Cmd])
		switch e.Cmd {
		case MetaByteMap["Sequence/Track name"]:
			out = fmt.Sprintf("%s -> %s", out, e.SeqTrackName)
		case MetaByteMap["Time Signature"]:
			out = fmt.Sprintf("%s -> %s", out, e.TimeSignature)
		case MetaByteMap["Copyright"]:
			out = fmt.Sprintf("%s -> %s", out, e.Copyright)
		case MetaByteMap["Tempo"]:
			out = fmt.Sprintf("%s -> %d", out, e.Bpm)
		}
	}

	return out
}

// Encode converts an Event into a slice of bytes ready to be written to a file.
func (e *Event) Encode() []byte {
	buff := bytes.NewBuffer(nil)
	buff.Write(EncodeVarint(e.TimeDelta))

	// msg type and chan are stored together
	msgData := []byte{(e.MsgType << 4) | e.MsgChan}
	if e.MsgType == EventByteMap["Meta"] {
		msgData = []byte{0xFF}
	}
	buff.Write(msgData)
	switch e.MsgType {
	// unknown but found in the wild (seems to come with 1 data bytes)
	case 0x2, 0x3, 0x4, 0x5, 0x6:
		buff.Write([]byte{0x0})
	// Note Off/On
	case 0x8, 0x9, 0xA:
		// note
		binary.Write(buff, binary.BigEndian, e.Note)
		// velocity
		binary.Write(buff, binary.BigEndian, e.Velocity)
		// Control Change / Channel Mode
		// This message is sent when a controller value changes.
		// Controllers include devices such as pedals and levers.
		// Controller numbers 120-127 are reserved as "Channel Mode Messages".
		// The controller number is between 0-119.
		// The new controller value is between 0-127.
	case 0xB:
		binary.Write(buff, binary.BigEndian, e.Controller)
		binary.Write(buff, binary.BigEndian, e.NewValue)
		/*
			channel mode messages
			Documented, not technically exposed

			This the same code as the Control Change, but implements Mode control and
			special message by using reserved controller numbers 120-127. The commands are:

			All Sound Off
			c = 120, v = 0
			When All Sound Off is received all oscillators will turn off,
			and their volume envelopes are set to zero as soon as possible.

			Reset All Controllers
			c = 121, v = x
			When Reset All Controllers is received, all controller values are reset to their default values.
			Value must only be zero unless otherwise allowed in a specific Recommended Practice.

			Local Control.
			c = 122, v = 0: Local Control Off
			c = 122, v = 127: Local Control On
			When Local Control is Off, all devices on a given channel will respond only to data received over MIDI.
			Played data, etc. will be ignored. Local Control On restores the functions of the normal controllers.

			All Notes Off.
			c = 123, v = 0: All Notes Off (See text for description of actual mode commands.)
			c = 124, v = 0: Omni Mode Off
			c = 125, v = 0: Omni Mode On
			c = 126, v = M: Mono Mode On (Poly Off) where M is the number of channels (Omni Off) or 0 (Omni On)
			c = 127, v = 0: Poly Mode On (Mono Off) (Note: These four messages also cause All Notes Off)
			When an All Notes Off is received, all oscillators will turn off.
			Program Change
					This message sent when the patch number changes. Value is the new program number.
		*/
	case 0xC:
		binary.Write(buff, binary.BigEndian, e.NewProgram)
		binary.Write(buff, binary.BigEndian, e.NewValue)
		// Channel Pressure (Aftertouch)
		// This message is most often sent by pressing down on the key after it "bottoms out".
		// This message is different from polyphonic after-touch.
		// Use this message to send the single greatest pressure value (of all the current depressed keys).
		// Value is the pressure value.
		// Most MIDI controllers don't generate Polyphonic Key AfterTouch because that requires a pressure sensor for each individual key
		// on a MIDI keyboard, and this is an expensive feature to implement.
		// For this reason, many cheaper units implement Channel Pressure instead of Aftertouch, as the former only requires
		// one sensor for the entire keyboard's pressure.
	case 0xD:
		binary.Write(buff, binary.BigEndian, e.Pressure)
		// Pitch Bend Change.
		// This message is sent to indicate a change in the pitch bender (wheel or lever, typically).
		// The pitch bender is measured by a fourteen bit value. Center (no pitch change) is 2000H.
		// Sensitivity is a function of the transmitter.
		// Last 7 bits of the first byte are the least significant 7 bits.
		// Last 7 bits of the second byte are the most significant 7 bits.
	case 0xE:
		// pitchbend
		lsb := byte(e.AbsPitchBend & 0x7F)
		msb := byte((e.AbsPitchBend & (0x7F << 7)) >> 7)
		binary.Write(buff, binary.BigEndian, []byte{lsb, msb})
		//  Meta
		// All meta-events start with FF followed by the command (xx), the length,
		// or number of bytes that will contain data (nn), and the actual data (dd).
		// meta_event = 0xFF + <meta_type> + <v_length> + <event_data_bytes>
	case 0xF:
		binary.Write(buff, binary.BigEndian, e.Cmd)
		switch e.Cmd {
		// Copyright Notice
		case 0x02:
			copyright := []byte(e.Copyright)
			binary.Write(buff, binary.BigEndian, EncodeVarint(uint32(len(copyright))))
			binary.Write(buff, binary.BigEndian, copyright)
			// BPM / tempo event
		case 0x51:
			binary.Write(buff, binary.BigEndian, EncodeVarint(3))
			binary.Write(buff, binary.BigEndian, Uint24(e.MsPerQuartNote))
		case 0x2f: // end of track
			buff.WriteByte(0x0)
		}
	default:
		fmt.Printf("didn't encode %#v because didn't know how to\n", e)
	}

	return buff.Bytes()
}

// Size represents the byte size to encode the event
func (e *Event) Size() uint32 {
	switch e.MsgType {
	case 0x2, 0x3, 0x4, 0x5, 0x6, 0xC, 0xD:
		return 1
	// Note Off, On, aftertouch, control change
	case 0x8, 0x9, 0xA, 0xB, 0xE:
		return 2
	case 0xF:
		// meta event
		switch e.Cmd {
		// Copyright Notice
		case 0x02:
			copyright := []byte(e.Copyright)
			varintBytes := EncodeVarint(uint32(len(copyright)))
			return uint32(len(copyright) + len(varintBytes))
			// BPM (size + encoded in uint24)
		case 0x51: // tempo
			return 4
		case 0x2f: // end of track
			return 1
		default:
			// NOT currently support, blowing up on purpose
			log.Fatal(errors.New("Can't encode this meta event, it is not supported yet"))
		}
	}
	return 0
}
