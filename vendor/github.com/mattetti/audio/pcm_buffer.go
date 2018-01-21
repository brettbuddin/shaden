package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	ErrInvalidBuffer = errors.New("invalid buffer")
)

// DataFormat is an enum type to indicate the underlying data format used.
type DataFormat int

const (
	// Unknown refers to an unknown format
	Unknown DataFormat = iota
	// Integer represents the int type.
	// it represents the native int format used in audio buffers.
	Integer
	// Float represents the float64 type.
	// It represents the native float format used in audio buffers.
	Float
	// Byte represents the byte type.
	Byte
)

// Format is a high level representation of the underlying data.
type Format struct {
	// NumChannels is the number of channels contained in the data
	NumChannels int
	// SampleRate is the sampling rate in Hz
	SampleRate int
	// BitDepth is the number of bits of data for each sample
	BitDepth int
	// Endianess indicate how the byte order of underlying bytes
	Endianness binary.ByteOrder
}

// PCMBuffer encapsulates uncompressed audio data
// and provides useful methods to read/manipulate this PCM data.
type PCMBuffer struct {
	// Format describes the format of the buffer data.
	Format *Format
	// Ints is a store for audio sample data as integers.
	Ints []int
	// Floats is a store for audio samples data as float64.
	Floats []float64
	// Bytes is a store for audio samples data as raw bytes.
	Bytes []byte
	// DataType indicates the primary format used for the underlying data.
	// The consumer of the buffer might want to look at this value to know what store
	// to use to optimaly retrieve data.
	DataType DataFormat
	// framePos is the position of the last frame we read
	framePos int64
}

// NewPCMIntBuffer returns a new PCM buffer backed by the passed integer samples
func NewPCMIntBuffer(data []int, format *Format) *PCMBuffer {
	return &PCMBuffer{
		Format:   format,
		DataType: Integer,
		Ints:     data,
	}
}

// NewPCMFloatBuffer returns a new PCM buffer backed by the passed float samples
func NewPCMFloatBuffer(data []float64, format *Format) *PCMBuffer {
	return &PCMBuffer{
		Format:   format,
		DataType: Float,
		Floats:   data,
	}
}

// NewPCMByteBuffer returns a new PCM buffer backed by the passed float samples
func NewPCMByteBuffer(data []byte, format *Format) *PCMBuffer {
	return &PCMBuffer{
		Format:   format,
		DataType: Byte,
		Bytes:    data,
	}
}

// Len returns the length of the underlying data.
func (b *PCMBuffer) Len() int {
	if b == nil {
		return 0
	}

	switch b.DataType {
	case Integer:
		return len(b.Ints)
	case Float:
		return len(b.Floats)
	case Byte:
		return len(b.Bytes)
	default:
		return 0
	}
}

// Size returns the number of frames contained in the buffer.
func (b *PCMBuffer) Size() (numFrames int) {
	if b == nil || b.Format == nil {
		return 0
	}
	numChannels := b.Format.NumChannels
	if numChannels == 0 {
		numChannels = 1
	}
	switch b.DataType {
	case Integer:
		numFrames = len(b.Ints) / numChannels
	case Float:
		numFrames = len(b.Floats) / numChannels
	case Byte:
		sampleSize := int((b.Format.BitDepth-1)/8 + 1)
		numFrames = (len(b.Bytes) / sampleSize) / numChannels
	}
	return numFrames
}

// Clone creates a clean clone that can be modified without
// changing the source buffer.
func (b *PCMBuffer) Clone() *PCMBuffer {
	if b == nil {
		return nil
	}
	newB := &PCMBuffer{DataType: b.DataType}
	switch b.DataType {
	case Integer:
		newB.Ints = make([]int, len(b.Ints))
		copy(newB.Ints, b.Ints)
	case Float:
		newB.Floats = make([]float64, len(b.Floats))
		copy(newB.Floats, b.Floats)
	case Byte:
		newB.Bytes = make([]byte, len(b.Bytes))
		copy(newB.Bytes, b.Bytes)
	}
	newB.Format = &Format{
		NumChannels: b.Format.NumChannels,
		SampleRate:  b.Format.SampleRate,
		BitDepth:    b.Format.BitDepth,
		Endianness:  b.Format.Endianness,
	}
	return newB
}

// AsInt16s returns the buffer samples as int16 sample values.
func (b *PCMBuffer) AsInt16s() (out []int16) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer, Float:
		out = make([]int16, len(b.Ints))
		for i := 0; i < len(b.Ints); i++ {
			out[i] = int16(b.Ints[i])
		}
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int16, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// AsInt32s returns the buffer samples as int32 sample values.
func (b *PCMBuffer) AsInt32s() (out []int32) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer, Float:
		out = make([]int32, len(b.Ints))
		for i := 0; i < len(b.Ints); i++ {
			out[i] = int32(b.Ints[i])
		}
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int32, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// AsInt64s returns the buffer samples as int64 sample values.
func (b *PCMBuffer) AsInt64s() (out []int64) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer, Float:
		out = make([]int64, len(b.Ints))
		for i := 0; i < len(b.Ints); i++ {
			out[i] = int64(b.Ints[i])
		}
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int64, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// AsInts returns the content of the buffer values as ints.
func (b *PCMBuffer) AsInts() (out []int) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer:
		return b.Ints
	case Float:
		out = make([]int, len(b.Floats))
		for i := 0; i < len(b.Floats); i++ {
			out[i] = int(b.Floats[i])
		}
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// AsFloat32s returns the buffer samples as float32 sample values.
func (b *PCMBuffer) AsFloat32s() (out []float32) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer:
		out = make([]float32, len(b.Ints))
		for i := 0; i < len(b.Ints); i++ {
			out[i] = float32(b.Ints[i])
		}
	case Float:
		out = make([]float32, len(b.Floats))
		for i := 0; i < len(b.Floats); i++ {
			out[i] = float32(b.Floats[i])
		}
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// AsFloat64s returns the buffer samples as float64 sample values.
func (b *PCMBuffer) AsFloat64s() (out []float64) {
	if b == nil {
		return nil
	}
	switch b.DataType {
	case Integer:
		out = make([]float64, len(b.Ints))
		for i := 0; i < len(b.Ints); i++ {
			out[i] = float64(b.Ints[i])
		}
	case Float:
		return b.Floats
	case Byte:
		// if the format isn't defined, we can't read the byte data
		if b.Format == nil || b.Format.Endianness == nil || b.Format.BitDepth == 0 {
			return out
		}
		bytesPerSample := int((b.Format.BitDepth-1)/8 + 1)
		buf := bytes.NewBuffer(b.Bytes)
		out := make([]int, len(b.Bytes)/bytesPerSample)
		binary.Read(buf, b.Format.Endianness, &out)
	}
	return out
}

// CacheInts ensures that the underlying int store is filled up
// so Ints() can be called knowing that the data is available.
// Note that if the underlying data is changed, it is the caller responsibility
// to refresh the cache.
func (b *PCMBuffer) CacheInts() {
	if b == nil || b.DataType == Integer {
		return
	}
	b.Ints = b.AsInts()
}

// CacheFloat64s ensures that the underlying int store is filled up
// so Floats() can be called knowing that the data is available.
// Note that if the underlying data is changed, it is the caller responsibility
// to refresh the cache.
func (b *PCMBuffer) CacheFloat64s() {
	if b == nil || b.DataType == Float {
		return
	}
	b.Ints = b.AsInts()
}

// SwitchPrimaryType is a convenience method to switch the primary data type.
// Use this if you process/swap a different type than the original type.
func (b *PCMBuffer) SwitchPrimaryType(t DataFormat) {
	if b == nil || t == b.DataType {
		return
	}
	switch t {
	case Integer:
		b.Ints = b.AsInts()
		b.Floats = nil
		b.Bytes = nil
	case Float:
		b.Floats = b.AsFloat64s()
		b.Ints = nil
		b.Bytes = nil
	case Byte:
		// TODO: get as bytes
		b.Floats = nil
		b.Ints = nil
	}
	b.DataType = t
}
