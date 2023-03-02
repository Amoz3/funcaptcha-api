package signals

import (
	"io"
	"os"

	"github.com/youpy/go-wav"
)

func ZeroCrossing(filename string) []int {
	crossings := make([]int, 0)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := wav.NewReader(file)

	defer file.Close()

	data := []wav.Sample{}
	for {
		sample, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		data = append(data, sample...)
	}

	lastOne := -1
	sinceCrossingCount := 0
	for _, sample := range data {
		entry := reader.IntValue(sample, 0)
		sinceCrossingCount++
		if lastOne > 0 && entry <= 0 {
			// fmt.Println("Pos to Neg crossing")
			crossings = append(crossings, sinceCrossingCount)
			lastOne = entry
			sinceCrossingCount = 0
			continue
		}

		if lastOne < 0 && entry >= 0 {
			// fmt.Println("Neg to Pos crossing")
			crossings = append(crossings, sinceCrossingCount)
			sinceCrossingCount = 0
			lastOne = entry
			continue
		}

		lastOne = entry
		sinceCrossingCount++
	}
	return crossings
}

func AvgCrossingTime(filename string) int {
	ZeroCrossings := ZeroCrossing(filename)
	sum := 0
	for _, zeroCrossing := range ZeroCrossings {
		sum += zeroCrossing
	}
	return sum / len(ZeroCrossings)
}
