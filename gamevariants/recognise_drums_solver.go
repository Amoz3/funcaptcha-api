package gamevariants

import "funcaptcha_api/signals"

var RecogniseDrumsSolver AudioSolver = AudioSolver{
	Splitter:   signals.PatternSplit,
	Classifier: classifyDrums,
}

// todo need to investigate signature aspects more
// https://media.discordapp.net/attachments/1080648036264448070/1084833084454735892/chart.png?width=547&height=547

func classifyDrums(splits ...string) int {
	for count, fileName := range splits {
		sig := signals.GetWaveSignature(fileName)
		if sig.StdDev < 2250 {
			return count + 1
		}
	}
	return 1
}
