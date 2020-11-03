package experiment

import (
	log "github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

func generateIAT(frequencySeconds float64, burstsNumber int, iatType string, experimentId int) []time.Duration {
	step := 1.0
	maxStep := frequencySeconds
	runningDelta := math.Min(maxStep, frequencySeconds)

	log.Debugf("Experiment %d: Generating %s IATs", experimentId, iatType)
	burstDeltas := make([]time.Duration, burstsNumber)
	for i := range burstDeltas {
		switch iatType {
		case "stochastic":
			burstDeltas[i] = time.Duration(getSleepTime(frequencySeconds)*1000) * time.Millisecond
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
			log.Errorf("Unrecognized inter-arrival time type %s, using default: stochastic", iatType)
			burstDeltas[i] = time.Duration(getSleepTime(frequencySeconds)*1000) * time.Millisecond
		}
	}
	return burstDeltas
}

// use a shifted and scaled exponential distribution to guarantee a minimum sleep time
func getSleepTime(frequencySeconds float64) float64 {
	rateParameter := 1 / math.Log(frequencySeconds) // experimentally deduced formula
	return frequencySeconds + rand.ExpFloat64()/rateParameter
}
