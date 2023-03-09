package signals

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

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

func PatternDetect(pattern []int, inputPath string) {

	wavFile, err := os.Open(inputPath)
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

	patternIndex := 0
	highestIndex := 0
	// index for the start and end of the pattern in challenge
	patternStart := 0
	patternEnd := 0

	for sampleIndex, sample := range samples {
		if patternIndex > highestIndex {
			highestIndex = patternIndex
		}
		if patternIndex >= len(pattern) {
			fmt.Println("break")
			patternEnd = sampleIndex
			break
		}

		entry := wavReader.IntValue(sample, 0)
		fmt.Println(entry)
		if pattern[patternIndex] != entry && patternIndex > 3 {
			fmt.Printf("pattern broken %d\n", patternIndex)
			time.Sleep(2 * time.Second)
			// break
		}

		if pattern[patternIndex] == entry {
			fmt.Printf("match %d\n", patternIndex)
			if patternIndex == 0 {
				patternStart = sampleIndex
			}
			patternIndex++
		} else {
			patternIndex = 0
		}
	}

	fmt.Printf("Pattern start %d, Pattern end %d, Longest pattern: %d Pattern length: %d, total samples length %d\n ",
		patternStart, patternEnd, highestIndex, len(pattern), len(samples))
}
