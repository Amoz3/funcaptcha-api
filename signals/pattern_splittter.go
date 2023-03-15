package signals

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/youpy/go-wav"
)

const SAMPLE_RATE = 44100

// time to cut is calculated off the starting samples index converted to time, so it includes the actual pattern
// the pattern will always be the same time so these values will work fine to cut out the "option x" audio, but will
// need to be modified if the patterns ever change significantly
const OPTION_ONE_LENGTH = 1
const OPTION_TWO_LENGTH = 1.8
const OPTION_THREE_LENGTH = 1.5

func PatternSplit(filePath, dirPath string) []string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	wavFile := dirPath + fmt.Sprintf("convertedToWav%d.wav", time.Now().Nanosecond())
	// convert mp3 to wav, maybe want mime type check here
	ffmpeg.Input(filePath).
		Output(wavFile).
		OverWriteOutput().
		Run()
	time.Sleep(3 * time.Second)

	// get split times based on pattern start index
	fristIndex, _ := PatternDetect(OptionOnePattern, wavFile) // i should get rid of this i think 1st is always the same
	secondIndex, _ := PatternDetect(OptionTwoPattern, wavFile)
	thirdIndex, ttlSamples := PatternDetect(OptionThreePattern, wavFile)

	if secondIndex > thirdIndex {
		color.Red("problem detecting third option in wav file: %s", wavFile)
		panic("")
	}

	// split and return file names
	guessOne := splitWav(wavFile, dirPath+timestamp+"guess_one.wav", indexToTime(fristIndex)+OPTION_ONE_LENGTH, timeBetweenTwoIndexes(fristIndex, secondIndex)-0.3)
	// -0.4 is just because the three pattern starts late
	guessTwo := splitWav(wavFile, dirPath+timestamp+"guess_two.wav", indexToTime(secondIndex)+OPTION_TWO_LENGTH, timeBetweenTwoIndexes(secondIndex, thirdIndex)-1)
	guessThree := splitWav(wavFile, dirPath+timestamp+"guess_three.wav", indexToTime(thirdIndex)+OPTION_THREE_LENGTH, timeBetweenTwoIndexes(thirdIndex, ttlSamples))

	return []string{guessOne, guessTwo, guessThree}
}

func splitWav(inputFilePath, output string, starttime, length float32) string {
	color.Cyan("splitting to output: %s w/ length: %f", output, length)
	// okay to copy codec because its been transformed to a wav prior
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
			if patternIndex > 5000 {
				// fmt.Printf("\n expected: %d got %d @ sample %d pattern index: %d\n", expectedVal, entry, sampleIndex, patternIndex)
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

	fmt.Printf("Longest pattern start: %d Pattern start %d, Pattern end %d, Longest pattern: %d Pattern length: %d, total samples length %d\n ",
		bestStart, patternStart, patternEnd, highestIndex, len(pattern), len(samples))
	return bestStart, len(samples)
}

func InBetween(i, min, max int) bool {
	return (i >= min) && (i <= max)
}
