package experiment

import (
	log "github.com/sirupsen/logrus"
	"lambda-benchmarking/client/experiment/configuration"
	"math"
	"math/rand"
	"time"
)

func generateIAT(experiment configuration.Experiment) []time.Duration {
	step := 1.0
	maxStep := experiment.CooldownSeconds
	runningDelta := math.Min(maxStep, experiment.CooldownSeconds)

	log.Debugf("Experiment %d: Generating %s IATs", experiment.Id, experiment.IATType)
	burstDeltas := make([]time.Duration, experiment.Bursts)
	for i := range burstDeltas {
		switch experiment.IATType {
		case "stochastic":
			burstDeltas[i] = time.Duration(getSpinTime(experiment.CooldownSeconds)*1000) * time.Millisecond
		case "deterministic":
			burstDeltas[i] = time.Duration(experiment.CooldownSeconds) * time.Second
		case "step":
			// TODO: TEST THIS and allow customization for runningDelta & step
			if i == 0 {
				burstDeltas[0] = time.Duration(runningDelta) * time.Second
			} else {
				burstDeltas[i] = time.Duration(math.Min(maxStep, runningDelta)) * time.Second
			}
			runningDelta += step
		default:
			log.Errorf("Unrecognized inter-arrival time type %s, using default: stochastic", experiment.IATType)
			burstDeltas[i] = time.Duration(getSpinTime(experiment.CooldownSeconds)*1000) * time.Millisecond
		}
	}
	return burstDeltas
}

// use a shifted and scaled exponential distribution to guarantee a minimum sleep time
func getSpinTime(frequencySeconds float64) float64 {
	rateParameter := 1 / math.Log(frequencySeconds) // experimentally deduced formula
	return frequencySeconds + rand.ExpFloat64()/rateParameter
}
