package configuration

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

var cachedServiceTimeIncrement map[string]int64

func determineFunctionIncrementLimits(subExperiment *SubExperiment, standardIncrement int64, standardDurationMs int64) {
	for _, serviceTime := range subExperiment.DesiredServiceTimes {
		if cachedIncrement, ok := cachedServiceTimeIncrement[serviceTime]; ok {
			log.Infof("Using cached increment %d for desired %v", cachedIncrement, serviceTime)
			subExperiment.FunctionIncrementLimits = append(subExperiment.FunctionIncrementLimits, cachedIncrement)
			continue
		}

		parsedDesiredDuration, err := time.ParseDuration(serviceTime)
		if err != nil {
			log.Fatalf("Could not parse desired function run duration %s from configuration file.", serviceTime)
		}

		desiredDurationMs := parsedDesiredDuration.Milliseconds()
		log.Infof("Determining function increment for a duration of %dms...", desiredDurationMs)

		ratio := big.NewRat(desiredDurationMs, standardDurationMs)
		currentIncrement := big.NewRat(standardIncrement, 1)
		currentIncrement.Mul(currentIncrement, ratio)

		suggestedIncrementFloat, _ := currentIncrement.Float64()
		suggestedIncrement := int64(suggestedIncrementFloat)
		suggestedDurationMs := timeSession(suggestedIncrement).Milliseconds()
		if !almostEqual(suggestedDurationMs, desiredDurationMs, float64(desiredDurationMs)*0.02) {
			log.Warnf("Suggested increment %d (duration %dms) is not within 2%% of desired duration %dms",
				suggestedIncrement, suggestedDurationMs, desiredDurationMs)

			log.Print("Please enter a better increment (leave empty for unchanged): ")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil && strings.Compare(err.Error(), "unexpected newline") != 0 {
				log.Fatalf("Could not read response: %s.", err.Error())
			} else if err == nil {
				parsedManualIncrement, err := strconv.ParseInt(response, 10, 64)
				if err != nil {
					log.Fatalf("Could not parse integer %s: %s.", response, err.Error())
				}
				suggestedIncrement = parsedManualIncrement
			}
		}

		log.Infof("Using increment %d (timed ~%dms) for desired %dms", suggestedIncrement, suggestedDurationMs, desiredDurationMs)
		cachedServiceTimeIncrement[serviceTime] = suggestedIncrement
		subExperiment.FunctionIncrementLimits = append(subExperiment.FunctionIncrementLimits, suggestedIncrement)
	}
}

func timeSession(increment int64) time.Duration {
	start := time.Now()
	for i := int64(0); i < increment; i++ {
	}
	return time.Since(start)
}

func almostEqual(a, b int64, float64EqualityThreshold float64) bool {
	return math.Abs(float64(a-b)) <= float64EqualityThreshold
}
