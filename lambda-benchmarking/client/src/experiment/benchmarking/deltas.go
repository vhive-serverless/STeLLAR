package benchmarking

import (
	"math"
	"math/rand"
	"time"
)

func CreateBurstDeltas(frequencySeconds int, burstsNumber int, randomization bool) []time.Duration {
	var burstDeltas []time.Duration
	if frequencySeconds != -1 {
		// latency profiler run, delta is constant
		burstDeltas = make([]time.Duration, burstsNumber)
		for i := range burstDeltas {
			if randomization {
				// scale and shift the standard normal distribution by frequencySeconds
				// make sure result is positive
				sampleBurst := math.Max(rand.NormFloat64()*float64(frequencySeconds)+float64(frequencySeconds), 0)
				burstDeltas[i] = time.Duration(sampleBurst*1000) * time.Millisecond
			} else {
				burstDeltas[i] = time.Duration(frequencySeconds) * time.Second
			}
		}
	} else {
		// cold start delta identifier run, delta varies so that the exact timeout can be identified
		burstDeltas = []time.Duration{
			time.Duration(0), // COLD START
			500 * time.Millisecond,
			time.Second, // WAIT A BIT MORE FOR INITIALIZATION
			time.Second,
			time.Second,
			5 * time.Second,
			15 * time.Second,
			30 * time.Second,
			45 * time.Second,
			time.Minute,
			5 * time.Minute,
			8 * time.Minute,
			10 * time.Minute,
			12 * time.Minute,
			20 * time.Minute,
			30 * time.Minute,
			time.Hour,
			2 * time.Hour,
		}
		burstDeltas = burstDeltas[:burstsNumber]
	}
	return burstDeltas
}
