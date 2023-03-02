package signals

import (
	"io"
	"math"
	"os"

	"github.com/fatih/color"
	"github.com/youpy/go-wav"
)

// todo add more meaningful data points, these work fine for what im doing so im leaving it for now
type WaveSignature struct {
	StdDev       float64
	NonZeroCount int
	Mean         int
	PosSum       int
	ZeroCount    int
}

// signature is just some math from the wave samples
func GetWaveSignature(filename string) WaveSignature {
	file, _ := os.Open(filename)
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

	highest := 0
	sum := 0
	posSum := 0
	zeroCount := 0
	nonZeroSize := 0
	for _, sample := range data {
		entry := reader.IntValue(sample, 0)
		// fmt.Println(fmt.Sprintf("Sample: %d", entry))
		if entry > highest {
			highest = entry
		}
		if entry == 0 {
			zeroCount++
			continue
		}
		nonZeroSize++
		sum += entry
		if entry > 0 {
			posSum += entry
		}
	}
	if nonZeroSize == 0 {
		color.Red("non zero size = 0")
		return WaveSignature{}
	}
	mean := sum / nonZeroSize
	var stdSum float64 = 0
	for _, sample := range data {
		entry := reader.IntValue(sample, 0)
		if entry == 0 {
			continue
		}
		stdSum += math.Pow(float64((entry - mean)), 2)

	}

	stdDev := math.Sqrt(stdSum / float64(nonZeroSize))

	// fmt.Printf("StdDev: %f \t nonZeroCount: %d \t Mean: %d \t posSum: %d \t zeroCount: %d\n", stdDev, nonZeroSize, mean, posSum, zeroCount)
	// fmt.Printf("StdDev: %f \t Mean: %d\n", stdDev, mean)
	return WaveSignature{
		StdDev:       stdDev,
		NonZeroCount: nonZeroSize,
		Mean:         mean,
		PosSum:       posSum,
		ZeroCount:    zeroCount,
	}
}
