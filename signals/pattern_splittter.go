package signals

import (
	"fmt"
	"io"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/youpy/go-wav"
)

func PatternSplit(filePath, dirPath string) []string {
	wavFile := dirPath + "convertedToWav.wav"
	// convert mp3 to wav, maybe want mime type check here
	ffmpeg.Input(filePath).
		Output(wavFile).
		OverWriteOutput().
		Run()

	// get split times based on pattern start index
	fristIndex := PatternDetect(OptionOnePattern, wavFile) // i should get rid of this i think 1st is always the same
	secondIndex := PatternDetect(OptionTwoPattern,  wavFile)
	thirdIndex := PatternDetect(OptionThreePattern, wavFile)

	// split and return file names
	guessOne := splitWav(wavFile, dirPath + "guess_one.wav", indexToTime(fristIndex), 1)
	guessTwo := splitWav(wavFile, dirPath + "guess_two.wav", indexToTime(secondIndex), 1)
	guessThree := splitWav(wavFile, dirPath + "guess_three.wav", indexToTime(thirdIndex), 1)
	
	return []string{guessOne, guessTwo, guessThree}
}

func splitWav(inputFilePath, output string, starttime, length float32) string {
	ffmpeg.Input(inputFilePath).
		Output(output, ffmpeg.KwArgs{"acodec": "copy"}, ffmpeg.KwArgs{"ss": 3.2}, ffmpeg.KwArgs{"t": 3}).
		OverWriteOutput().
		Run()
		return output
}

func indexToTime(sampleIndex int) float32 {
	sampleRate := 44100
	return float32(sampleIndex / sampleRate)
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
