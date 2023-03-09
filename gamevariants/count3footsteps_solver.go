package gamevariants

import (
	"fmt"
	"funcaptcha_api/signals"
)

var ThreeFootstepSolver AudioSolver = AudioSolver{
	Splitter: signals.PatternSplit,
	Classifier: classifyThreeFootstep,
}

func classifyThreeFootstep(splits ...string) int {
	fmt.Println(splits)
	return 3
}