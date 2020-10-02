package benchmarking

import (
	"time"
)

func RunProfiler(burstDeltas []time.Duration, requests int, execMilliseconds int, payloadLengthBytes int) {
	burstRelativeDeltas := []time.Duration{burstDeltas[0]}
	for _, burstDelta := range burstDeltas[1:] {
		burstRelativeDeltas = append(burstRelativeDeltas,
			burstRelativeDeltas[len(burstRelativeDeltas)-1]+burstDelta)
	}

	triggerRelativeAsyncBurstGroups(burstRelativeDeltas, requests, execMilliseconds, payloadLengthBytes)
}
