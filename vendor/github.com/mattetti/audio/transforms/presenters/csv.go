package presenters

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/mattetti/audio"
)

// CSV writes the content of the buffer in a CSV file.
// Can be used to plot with R for instance:
//      Rscript -e 'png("my_plot.png", height = 768, width = 1024);
//      myData <- read.csv("./data.csv"); matplot(myData[,1], type="l")'
func CSV(buf *audio.PCMBuffer, path string) error {
	if buf == nil || buf.Format == nil {
		return audio.ErrInvalidBuffer
	}
	csvf, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvf.Close()
	csvw := csv.NewWriter(csvf)
	buf.SwitchPrimaryType(audio.Float)
	row := make([]string, buf.Format.NumChannels)

	for i := 0; i < buf.Format.NumChannels; i++ {
		row[i] = fmt.Sprintf("Channel %d", i+1)
	}

	if err := csvw.Write(row); err != nil {
		return fmt.Errorf("error writing header to csv: %s", err)
	}

	totalSize := buf.Len()
	for i := 0; i < totalSize; i++ {
		for j := 0; j < buf.Format.NumChannels; j++ {
			row[j] = fmt.Sprintf("%v", buf.Floats[i*buf.Format.NumChannels+j])
			if i >= totalSize {
				break
			}
		}
		if err := csvw.Write(row); err != nil {
			return fmt.Errorf("error writing record to csv: %s", err)
		}
	}
	csvw.Flush()
	return csvw.Error()
}
