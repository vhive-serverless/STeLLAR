package benchmarking

import (
	"time"
)

func MakeBurstDeltasRelative(burstDeltas []time.Duration) []time.Duration {
	burstRelativeDeltas := []time.Duration{burstDeltas[0]}
	for _, burstDelta := range burstDeltas[1:] {
		burstRelativeDeltas = append(burstRelativeDeltas,
			burstRelativeDeltas[len(burstRelativeDeltas)-1]+burstDelta)
	}
	return burstRelativeDeltas
}

func CreateBurstDeltas(frequencySeconds int, burstsNumber int) []time.Duration {
	var burstDeltas []time.Duration
	if frequencySeconds != -1 {
		// latency profiler run, delta is constant
		burstDeltas = make([]time.Duration, burstsNumber)
		for i := range burstDeltas {
			burstDeltas[i] = time.Duration(frequencySeconds) * time.Second
		}
	} else {
		// cold start delta identifier run, delta varies so that the exact timeout can be identified
		burstDeltas = []time.Duration{
			time.Duration(0),
			500 * time.Millisecond,
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
