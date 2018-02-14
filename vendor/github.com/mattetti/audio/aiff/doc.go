/*
Package aiff is a AIFF/AIFC decoder and encoder.
It extracts the basic information of the file and provide decoded frames for AIFF files.


This package also allows for quick access to the AIFF LPCM raw audio data:

    in, err := os.Open("audiofile.aiff")
    if err != nil {
    	log.Fatal("couldn't open audiofile.aiff %v", err)
    }
	d := NewDecoder(in)
    frames, err := d.Frames()
    in.Close()

A frame is a slice where each entry is a channel and each value is the sample value.
For instance, a frame in a stereo file will have 2 entries (left and right) and each entry will
have the value of the sample.

Note that this approach isn't memory efficient at all. In most cases
you want to access the decoder's Clip.
You can read the clip's frames in a buffer that you decode in small chunks.

Finally, the encoder allows the encoding of LPCM audio data into a valid AIFF file.
Look at the encoder_test.go file for a more complete example.

Currently only AIFF is properly supported, AIFC files will more than likely not be properly processed.

*/
package aiff
