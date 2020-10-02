package benchmarking

import (
	"log"
	"sync"
	"time"
)

func triggerRelativeAsyncBurstGroups(burstRelativeDeltas []time.Duration, requests int, execMilliseconds int, payloadLengthBytes int) {
	var burstsWaitGroup sync.WaitGroup
	for burstId, relativeDelta := range burstRelativeDeltas {
		log.Printf("Scheduling burst %d (%v) for %v.",
			burstId,
			burstRelativeDeltas[burstId],
			time.Now().Add(burstRelativeDeltas[burstId]).Format(time.StampNano))
		burstsWaitGroup.Add(1)
		go burst(&burstsWaitGroup, burstRelativeDeltas, burstId, relativeDelta, requests, execMilliseconds, payloadLengthBytes)
	}
	burstsWaitGroup.Wait()
}

func burst(burstsWaitGroup *sync.WaitGroup, burstRelativeDeltas []time.Duration, burstId int,
	relativeDelta time.Duration, requests int, execMilliseconds int, payloadLengthBytes int) {
	defer burstsWaitGroup.Done()
	time.Sleep(relativeDelta)

	log.Printf("Starting burst %d (%v) on %v: making %d requests to API Gateway.",
		burstId,
		burstRelativeDeltas[burstId],
		time.Now().Format(time.StampNano),
		requests,
	)

	var requestsWaitGroup sync.WaitGroup
	for i := 0; i < requests; i++ {
		requestsWaitGroup.Add(1)
		go SafeWriterInstance.GenerateLatencyRecord(&requestsWaitGroup, execMilliseconds, payloadLengthBytes, burstId)
	}
	requestsWaitGroup.Wait()
	log.Printf("Received all responses for burst %d.", burstId)
}
