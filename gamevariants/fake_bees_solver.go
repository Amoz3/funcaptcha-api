package gamevariants

import (
	"fmt"
	"funcaptcha_api/signals"

	"github.com/fatih/color"
)

var FakeBeesSolver AudioSolver = AudioSolver{
	Splitter:   signals.PatternSplit,
	Classifier: classifyFakeBees,
}

// https://media.discordapp.net/attachments/1080648036264448070/1084828383420088420/chart.png?width=547&height=547

func classifyFakeBees(splits ...string) int {
	for count, fileName := range splits {
		avgCrossing := signals.AvgCrossingTime(fileName)
		sig := signals.GetWaveSignature(fileName)
		color.Cyan(fmt.Sprintf("crossing: %d stddev: %f - fake bees", avgCrossing, sig.StdDev))
		if avgCrossing >= 30 {
			return count + 1
		}
	}
	color.Red("didnt detect fake bees in given samples")
	return 1 // random
}
