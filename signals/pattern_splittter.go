package signals

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/youpy/go-wav"
)

const SAMPLE_RATE = 44100

// time to cut is calculated off the starting samples index converted to time, so it includes the actual pattern
// the pattern will always be the same time so these values will work fine to cut out the "option x" audio, but will
// need to be modified if the patterns ever change significantly
const OPTION_ONE_LENGTH = 3.15
const OPTION_TWO_LENGTH = 1.15
const OPTION_THREE_LENGTH = 1.4

func PatternSplit(filePath, dirPath string) []string {
	wavFile := dirPath + "convertedToWav.wav"
	// convert mp3 to wav, maybe want mime type check here
	ffmpeg.Input(filePath).
		Output(wavFile).
		OverWriteOutput().
		Run()

	// get split times based on pattern start index
	fristIndex, _ := PatternDetect(OptionOnePattern, wavFile) // i should get rid of this i think 1st is always the same
	secondIndex, _ := PatternDetect(OptionTwoPattern, wavFile)
	thirdIndex, ttlSamples := PatternDetect(OptionThreePattern, wavFile)

	// split and return file names
	guessOne := splitWav(wavFile, dirPath+"guess_one.wav", indexToTime(fristIndex)+OPTION_ONE_LENGTH, timeBetweenTwoIndexes(fristIndex, secondIndex)-OPTION_ONE_LENGTH)
	guessTwo := splitWav(wavFile, dirPath+"guess_two.wav", indexToTime(secondIndex)+OPTION_TWO_LENGTH, timeBetweenTwoIndexes(secondIndex, thirdIndex)-OPTION_TWO_LENGTH)
	guessThree := splitWav(wavFile, dirPath+"guess_three.wav", indexToTime(thirdIndex)+OPTION_THREE_LENGTH, timeBetweenTwoIndexes(thirdIndex, ttlSamples))

	return []string{guessOne, guessTwo, guessThree}
}

func splitWav(inputFilePath, output string, starttime, length float32) string {
	color.Cyan("splitting to output: %s w/ length: %f", output, length)
	ffmpeg.Input(inputFilePath).
		Output(output, ffmpeg.KwArgs{"acodec": "copy"}, ffmpeg.KwArgs{"ss": starttime}, ffmpeg.KwArgs{"t": length}).
		OverWriteOutput().
		Run()
	return output
}

func timeBetweenTwoIndexes(indexOne, indexTwo int) float32 {
	if indexOne > indexTwo {
		return float32((indexOne - indexTwo) / SAMPLE_RATE)
	}
	return float32((indexTwo - indexOne) / SAMPLE_RATE)
}

func indexToTime(sampleIndex int) float32 {
	return float32(sampleIndex / SAMPLE_RATE)
}

// 2nd return is the size of all the samples so i can get from option 3 to eof
func PatternDetect(pattern []int, inputPath string) (int, int) {

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
		tolerence := 10
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
	return bestStart, len(samples)
}

func InBetween(i, min, max int) bool {
	return (i >= min) && (i <= max)
}
