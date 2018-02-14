/*
caf package implements an Apple Core Audio Format (CAF) parser.

Apple’s Core Audio Format (CAF) is a file format (container) for storing and transporting digital audio data.
It simplifies the management and manipulation of many types of audio data without the file-size limitations
of other audio file formats.
CAF provides high performance and flexibility and is scalable to future ultra-high resolution audio recording, editing, and playback.
CAF files can contain audio or even define patches, or musical voice configurations.

The primary goal of this package is to allow the transcoding/conversion of CAF files to other formats.
Note that because CAF fiels contain other metadata than just audio data, the conversion will be lossy, not in the sound
quality meaning of the sense but in the data senses.


That said here is some information about CAF provided by Apple.

CAF files have several advantages over other standard audio file formats:

* Unrestricted file size
Whereas AIFF, AIFF-C, and WAV files are limited in size to 4 gigabytes, which might represent as little as 15 minutes of audio,
CAF files use 64-bit file offsets, eliminating practical limits.
A standard CAF file can hold audio data with a playback duration of hundreds of years.

* Safe and efficient recording
Applications writing AIFF and WAV files must either update the data header’s size field at the end of recording—which can result in an unusable file
if recording is interrupted before the header is finalized—or they must update the size field after recording each packet of data, which is inefficient. With CAF files,
in contrast, an application can append new audio data to the end of the file in a manner that allows it
to determine the amount of data even if the size field in the header has not been finalized.

* Support for many data formats
CAF files serve as wrappers for a wide variety of audio data formats.
The flexibility of the CAF file structure and the many types of metadata that can be recorded enable CAF files to be used with practically any type of audio data.
Furthermore, CAF files can store any number of audio channels.

* Support for many types of auxiliary data
In addition to audio data, CAF files can store text annotations, markers, channel layouts, and many other types of information that can help in the interpretation,
analysis, or editing of the audio.

* Support for data dependencies
Certain metadata in CAF files is linked to the audio data by an edit count value.
You can use this value to determine when metadata has a dependency on the audio data and, furthermore,
when the audio data has changed since the metadata was written.


Vocab:

* Sample
One number for one channel of digitized audio data.

* Frame
A set of samples representing one sample for each channel. The samples in a frame are intended to be played together (that is, simultaneously).
Note that this definition might be different from the use of the term “frame” by codecs, video files, and audio or video processing applications.

* Packet
The smallest, indivisible block of data. For linear PCM (pulse-code modulated) data, each packet contains exactly one frame.
For compressed audio data formats, the number of frames in a packet depends on the encoding.
For example, a packet of AAC represents 1024 frames of PCM. In some formats, the number of frames per packet varies.

* Sample rate
The number of complete frames of samples per second of noncompressed or decompressed data.



The format is documented by Apple at
https://developer.apple.com/library/mac/documentation/MusicAudio/Reference/CAFSpec/CAF_spec/CAF_spec.html#//apple_ref/doc/uid/TP40001862-CH210-TPXREF101

*/
package caf
