package gamevariants

import (
	"fmt"
	"funcaptcha_api/signals"
	"time"

	"github.com/fatih/color"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var RepeatingSongSolver AudioSolver = AudioSolver{
	Splitter:   signals.PatternSplit,
	Classifier: classifyDrums,
}

func SplitRepeatingSongs(filename, dirPath string) []string {
	timestamp := time.Now().UnixMicro()
	outputOne := fmt.Sprintf("%sguess_0_%d.wav", dirPath, timestamp)
	outputTwo := fmt.Sprintf("%sguess_1_%d.wav", dirPath, timestamp)
	outputThree := fmt.Sprintf("%sguess_2_%d.wav", dirPath, timestamp)
	color.Magenta(outputOne)
	ffmpeg.Input(filename).
		Output(outputOne, ffmpeg.KwArgs{"ss": 3.2}, ffmpeg.KwArgs{"t": 4}).
		OverWriteOutput().
		Run()
	ffmpeg.Input(filename).
		Output(outputTwo, ffmpeg.KwArgs{"ss": 8.5}, ffmpeg.KwArgs{"t": 4}).
		OverWriteOutput().
		Run()
	ffmpeg.Input(filename).
		Output(outputThree, ffmpeg.KwArgs{"ss": 13.7}, ffmpeg.KwArgs{"t": 4}).
		OverWriteOutput().
		Run()

	return []string{
		outputOne,
		outputTwo,
		outputThree,
	}
}

// this classifier is terrible signal lib needs a better huerstic for repetition
func classifyRepeatingSong(splits ...string) int {
	for count, fileName := range splits {
		sig := signals.GetWaveSignature(fileName)
		if sig.StdDev >= 5000 && sig.StdDev <= 8000 {
			return count + 1
		}
	}
	return 1 // random
}
