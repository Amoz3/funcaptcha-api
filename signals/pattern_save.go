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

func PatternDetect(pattern []int, inputPath string) int {

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
	// best start is the start index for the longest pattern
	bestStart := 0

	fmt.Println(len(samples))
	for sampleIndex, sample := range samples {
		if patternIndex > highestIndex {
			bestStart = patternStart
			highestIndex = patternIndex
		}

		if patternIndex >= len(pattern) {
			patternEnd = sampleIndex
			fmt.Println("Pattern found.")
			break
		}

		entry := wavReader.IntValue(sample, 0)
		// bits seem to be off by a bit sometimes this will allow for inperfect matches
		tolerence := 100
		expectedVal := pattern[patternIndex]
		// if entry != expectedVal {
		if !InBetween(entry, expectedVal-tolerence, expectedVal+tolerence) {
			if patternIndex > 5 {
				fmt.Printf("\n expected: %d got %d @ sample %d pattern index: %d\n", expectedVal, entry, sampleIndex, patternIndex)
			}
			patternIndex = 0
		}

		expectedVal = pattern[patternIndex]
		// if entry == expectedVal {
		if InBetween(entry, expectedVal-tolerence, expectedVal+tolerence) {
			if patternIndex == 0 {
				patternStart = sampleIndex
				if bestStart == 0 {
					bestStart = patternStart
				}
			}
			patternIndex++
		}
	}

	fmt.Printf("Pattern start %d, Pattern end %d, Longest pattern: %d Pattern length: %d, total samples length %d\n ",
		patternStart, patternEnd, highestIndex, len(pattern), len(samples))
	return bestStart
}

func InBetween(i, min, max int) bool {
	return (i >= min) && (i <= max)
}
