package midi

// EventMap takes a event byte and returns its matching event name.
// http://www.midi.org/techspecs/midimessages.php
var EventMap = map[byte]string{
	0x8: "NoteOff",
	0x9: "NoteOn",
	0xA: "AfterTouch",
	0xB: "ControlChange",
	0xC: "ProgramChange",
	0xD: "ChannelAfterTouch",
	0xE: "PitchWheelChange",
	0xF: "Meta",
}

// EventByteMap takes an event name and returns its equivalent MIDI byte
var EventByteMap = map[string]byte{
	"NoteOff":           0x8,
	"NoteOn":            0x9,
	"AfterTouch":        0xA,
	"ControlChange":     0xB,
	"ProgramChange":     0xC,
	"ChannelAfterTouch": 0xD,
	"PitchWheelChange":  0xE,
	"Meta":              0xF,
}

// MetaCmdMap maps metadata binary command to their names
var MetaCmdMap = map[byte]string{
	0x0:  "Sequence number",
	0x01: "Text event",
	0x02: "Copyright",
	0x03: "Sequence/Track name",
	0x04: "Instrument name",
	0x05: "Lyric",
	0x06: "Marker",
	0x07: "Cue Point",
	0x20: "MIDI Channel Prefix",
	0x2f: "End of Track",
	0x51: "Tempo",
	0x58: "Time Signature",
	0x59: "Key Signature",
	0x7F: "Sequencer specific",
	0x8F: "Timing Clock",
	0xFA: "Start current sequence",
	0xFB: "Continue stopped sequence where left off",
	0xFC: "Stop sequence",
}

// MetaByteMap maps metadata command names to their binary cmd code
var MetaByteMap = map[string]byte{
	"Sequence number":                          0x0,
	"Text event":                               0x01,
	"Copyright":                                0x02,
	"Sequence/Track name":                      0x03,
	"Instrument name":                          0x04,
	"Lyric":                                    0x05,
	"Marker":                                   0x06,
	"Cue Point":                                0x07,
	"MIDI Channel Prefix":                      0x20,
	"End of Track":                             0x2f,
	"Tempo":                                    0x51,
	"Time Signature":                           0x58,
	"Key Signature":                            0x59,
	"Sequencer specific":                       0x7F,
	"Timing Clock":                             0x8F,
	"Start current sequence":                   0xFA,
	"Continue stopped sequence where left off": 0xFB,
	"Stop sequence":                            0xFC,
}
