package gamevariants

import "funcaptcha_api/signals"

var FakeBeesSolver AudioSolver = AudioSolver{
	Splitter:   signals.PatternSplit,
	Classifier: classifyFakeBees,
}
// https://media.discordapp.net/attachments/1080648036264448070/1084828383420088420/chart.png?width=547&height=547

func classifyFakeBees(splits ...string) int {
	for count, fileName := range splits {
		avgCrossing := signals.AvgCrossingTime(fileName)
		sig := signals.GetWaveSignature(fileName)
		if avgCrossing >= 30 && sig.StdDev < 5000 {
			return count + 1
		}
	}
	return 1 // random
}
