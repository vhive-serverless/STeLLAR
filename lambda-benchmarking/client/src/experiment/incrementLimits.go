package experiment

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func GetIncrementLimits(experimentIndex int, serviceTimesStrings []string) []int {
	var lambdaIncrementLimits []int
	for _, serviceTime := range serviceTimesStrings {
		parsedServiceTime, err := time.ParseDuration(serviceTime)
		if err != nil {
			log.Fatal(err)
		}
		lambdaIncrementLimits = append(lambdaIncrementLimits, inferIncrementLimit(experimentIndex, parsedServiceTime))
	}
	return lambdaIncrementLimits
}

func inferIncrementLimit(experimentIndex int, serviceTime time.Duration) int {
	log.Debugf("Experiment %d: inferring increment limit (on client) for a service time of %v",
		experimentIndex, serviceTime)

	start := time.Now()
	serviceTimeTicker := time.NewTimer(serviceTime)

	increment := 0
	receivedSignal := false
	for !receivedSignal {
		select {
		case <-serviceTimeTicker.C:
			receivedSignal = true
			break
		default:
			increment++
		}
	} // 1.35e9 for 400ms 300 000 000

	log.Debugf("Experiment %d: inferring increment %d from service time %v took %v", experimentIndex, increment,
		serviceTime, time.Since(start))
	return increment
}
