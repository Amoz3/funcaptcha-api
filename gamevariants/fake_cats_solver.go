package gamevariants

import (
	"funcaptcha_api/signals"

	"github.com/fatih/color"
)

var FakeCatsSolver AudioSolver = AudioSolver{
	Splitter:   signals.PatternSplit,
	Classifier: classifyFakeCat,
}

func classifyFakeCat(splits ...string) int {
	for count, fileName := range splits {
		avgCrossing := signals.AvgCrossingTime(fileName)
		color.Red("avgcrossing - fake cats: %d", avgCrossing)
		if avgCrossing >= 80 {
			return count + 1
		}
	}
	return 1 // random
}
