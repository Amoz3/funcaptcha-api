package signals

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/youpy/go-wav"
)

func SavePattern(inputfile, outputFile string) {
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}

	// the text body we will save to output file
	fileData := ""

	wavFile, err := os.Open(inputfile)
	if err != nil {
		panic(err)
	}
	defer wavFile.Close()
	wavReader := wav.NewReader(wavFile)

	samples := []wav.Sample{}
	for {
		sample, err := wavReader.ReadSamples()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		samples = append(samples, sample...)
	}

	for _, sample := range samples {
		entry := wavReader.IntValue(sample, 0)
		fileData += fmt.Sprintf("%d,\n", entry)
	}

	ioutil.WriteFile(file.Name(), []byte(fileData), 0777)
}
