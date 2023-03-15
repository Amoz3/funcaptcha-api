package gamevariants

import (
	"fmt"
	"funcaptcha_api/signals"

	"github.com/fatih/color"
	"github.com/sajari/regression"
)

type AudioSolver struct {
	// all these may be returning , error in the future unsure rn
	// might be justified to have a pre processing function for ensuring the file format is correct but provider and consumer are the same party so idc

	// takes path to audio file, splits and returns 3 paths, 3 paths should be in chronological order
	// file to split, dir to store splits in.
	Splitter func(string, string) []string
	// should take varargs string (output from Splitter) returns guess, 1, 2, 3
	Classifier func(...string) int
}

func (a *AudioSolver) Solve(filePath, dirPath string) int {
	splits := a.Splitter(filePath, dirPath)
	return a.Classifier(splits...)
}

type RegressionSolver struct {
	*AudioSolver
	model   *regression.Regression
	posData [][]float64
	negData [][]float64
}

func NewRegressionSolver(splitter func(string, string) []string, posData, negData [][]float64) *RegressionSolver {
	posDataPoints := regression.DataPoints{}
	for _, data := range posData {
		posDataPoints = append(posDataPoints, regression.DataPoint(1, data))
	}

	negDataPoints := regression.DataPoints{}
	for _, data := range negData {
		negDataPoints = append(posDataPoints, regression.DataPoint(0, data))
	}
	allPoints := append(posDataPoints, negDataPoints...)

	r := new(regression.Regression)
	r.Train(allPoints...)
	r.Run()

	fmt.Printf("Regression formula:\n%v\n", r.Formula)
	fmt.Printf("Regression:\n%s\n", r)
	return &RegressionSolver{
		&AudioSolver{
			Splitter:   splitter,
			Classifier: func(s ...string) int { 
				bestPred := 0.0
				bestIndex := 1
				for count, filename := range s {
					sig := signals.GetWaveSignature(filename)
					data := []float64{sig.StdDev, float64(sig.NonZeroCount), float64(sig.Mean), float64(sig.PosSum), float64(sig.ZeroCount)}
					pred, err := r.Predict(data)
					if err != nil {
						color.Red("error getting prediciton")
						panic(err)
					}
					if pred > bestPred {
						bestPred = pred
						bestIndex = count + 1
					}
				}
				return bestIndex 
			},
		},
		r,
		posData,
		negData,
	}
}

