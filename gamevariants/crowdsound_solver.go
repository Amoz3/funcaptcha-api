package gamevariants

import (
	"fmt"
	"funcaptcha_api/signals"
	"time"

	"github.com/fatih/color"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var CrowdSoundSolver AudioSolver = AudioSolver{
	Splitter:   crowdSplitter,
	Classifier: findCrowdSound,
}

// optimally this would be based on patterns of the "Option 1", "Option 2" in all the challenges and we would pass around files
// we would pass like, straight up sample data
func crowdSplitter(filename, dirPath string) []string {
	timestamp := time.Now().UnixMicro()
	outputOne := fmt.Sprintf("%sguess_0_%d.wav", dirPath, timestamp)
	outputTwo := fmt.Sprintf("%sguess_1_%d.wav", dirPath, timestamp)
	outputThree := fmt.Sprintf("%sguess_2_%d.wav", dirPath, timestamp)
	color.Magenta(outputOne)
	ffmpeg.Input(filename).
		Output(outputOne, ffmpeg.KwArgs{"acodec": "copy"}, ffmpeg.KwArgs{"ss": 3.2}, ffmpeg.KwArgs{"t": 3}).
		OverWriteOutput().
		Run()
	ffmpeg.Input(filename).
		Output(outputTwo, ffmpeg.KwArgs{"ss": 7.5}, ffmpeg.KwArgs{"t": 3}).
		OverWriteOutput().
		Run()
	ffmpeg.Input(filename).
		Output(outputThree, ffmpeg.KwArgs{"ss": 11.8}, ffmpeg.KwArgs{"t": 3}).
		OverWriteOutput().
		Run()

	return []string{
		outputOne,
		outputTwo,
		outputThree,
	}
}

// this is just something i copied from the early mocking gopher lib, from memory its like 60-70% accurate (?)
func findCrowdSound(splits ...string) int {
	expectedStdDev := 3250
	expectedZeroCount := 2600

	time.Sleep(time.Second * 4)
	lowestDifference := 90_000_000
	likelyCrowdsound := 0
	for count, split := range splits {
		sig := signals.GetWaveSignature(split)
		if sig.ZeroCount > expectedZeroCount {
			continue
		}

		diff := expectedStdDev - int(sig.StdDev)
		if diff < 0 {
			diff = diff * -1
		}
		if diff < lowestDifference {
			lowestDifference = diff
			likelyCrowdsound = count + 1
		}
	}
	return likelyCrowdsound
}
