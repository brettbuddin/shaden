package main

func main() {
	/*
		path, _ := filepath.Abs("../../../decimator/beat.aiff")
		ext := filepath.Ext(path)
		var codec string
		switch strings.ToLower(ext) {
		case ".aif", ".aiff":
			codec = "aiff"
		case ".wav", ".wave":
			codec = "wav"
		default:
			fmt.Printf("files with extension %s not supported\n", ext)
		}

		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var monoFrames audio.Frames
		var sampleRate int
		var sampleSize int
		switch codec {
		case "aiff":
			d := aiff.NewDecoder(f)
			frames, err := d.Frames()
			if err != nil {
				panic(err)
			}
			sampleRate = d.SampleRate
			sampleSize = int(d.BitDepth)
			monoFrames = audio.ToMonoFrames(frames)

		case "wav":
			info, frames, err := wav.NewDecoder(f, nil).ReadFrames()
			if err != nil {
				panic(err)
			}
			sampleRate = int(info.SampleRate)
			sampleSize = int(info.BitsPerSample)
			monoFrames = audio.ToMonoFrames(frames)
		}

		data := make([]float64, len(monoFrames))
		for i, f := range monoFrames {
			data[i] = float64(f[0])
		}
		dft := analysis.NewDFT(sampleRate, data)
		sndData := dft.IFFT()
		frames := make([][]int, len(sndData))
		for i := 0; i < len(frames); i++ {
			frames[i] = []int{int(sndData[i])}
		}
		of, err := os.Create("roundtripped.aiff")
		if err != nil {
			panic(err)
		}
		defer of.Close()
		aiffe := aiff.NewEncoder(of, sampleRate, sampleSize, 1)
		aiffe.Frames = frames
		if err := aiffe.Write(); err != nil {
			panic(err)
		}
	*/
}
