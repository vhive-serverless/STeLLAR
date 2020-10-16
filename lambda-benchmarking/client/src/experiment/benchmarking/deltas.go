package benchmarking

import (
	"log"
	"math"
	"math/rand"
	"time"
)

func CreateBurstDeltas(frequencySeconds float64, burstsNumber int, iatType string) []time.Duration {
	step := 1.0
	maxStep := frequencySeconds
	runningDelta := math.Min(maxStep, frequencySeconds)

	var burstDeltas []time.Duration
	burstDeltas = make([]time.Duration, burstsNumber)
	for i := range burstDeltas {
		switch iatType {
		case "stochastic":
			// TODO: allow customization for mean (currently frequencySeconds) and stddev (currently frequencySeconds)
			burstDeltas[i] = time.Duration(getGaussianSleepTime(frequencySeconds)*1000) * time.Millisecond
		case "deterministic":
			burstDeltas[i] = time.Duration(frequencySeconds) * time.Second
		case "step":
			// TODO: TEST THIS and allow customization for runningDelta & step
			if i == 0 {
				burstDeltas[0] = time.Duration(runningDelta) * time.Second
			} else {
				burstDeltas[i] = time.Duration(math.Min(maxStep, runningDelta)) * time.Second
			}
			runningDelta += step
		default:
			log.Fatalf("Unrecognized inter-arrival time type %s", iatType)
		}
	}
	return burstDeltas
}

// scale and shift the standard normal distribution, make sure result is positive
func getGaussianSleepTime(frequencySeconds float64) float64 {
	return math.Max(rand.NormFloat64()*frequencySeconds+frequencySeconds, 0)
}
