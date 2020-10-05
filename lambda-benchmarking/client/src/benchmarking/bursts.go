package benchmarking

import (
	"log"
	"sync"
	"time"
)

func TriggerRelativeAsyncBurstGroups(gatewayEndpoint string, burstRelativeDeltas []time.Duration, requests int, lambdaIncrementLimit int, payloadLengthBytes int) {
	var burstsWaitGroup sync.WaitGroup
	for burstId, relativeDelta := range burstRelativeDeltas {
		log.Printf("Scheduling burst %d (%v) for %v.",
			burstId,
			burstRelativeDeltas[burstId],
			time.Now().Add(burstRelativeDeltas[burstId]).Format(time.StampNano))
		burstsWaitGroup.Add(1)
		go burst(gatewayEndpoint, &burstsWaitGroup, burstRelativeDeltas, burstId, relativeDelta, requests, lambdaIncrementLimit, payloadLengthBytes)
	}
	burstsWaitGroup.Wait()
}

func burst(gatewayEndpoint string, burstsWaitGroup *sync.WaitGroup, burstRelativeDeltas []time.Duration, burstId int,
	relativeDelta time.Duration, requests int, lambdaIncrementLimit int, payloadLengthBytes int) {
	defer burstsWaitGroup.Done()
	time.Sleep(relativeDelta)

	log.Printf("Starting burst %d (%v) on %v: making %d requests to API Gateway (%s).",
		burstId,
		burstRelativeDeltas[burstId],
		time.Now().Format(time.StampNano),
		requests,
		gatewayEndpoint,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go SafeWriterInstance.GenerateLatencyRecord(gatewayEndpoint, &requestsWaitGroup, lambdaIncrementLimit, payloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Printf("Received all responses for burst %d.", burstId)
}
