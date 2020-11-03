package experiment

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

func getIncrementLimits(incrementLimitStrings []string) []int64 {
	var functionIncrementLimits []int64
	for _, incrementLimitString := range incrementLimitStrings {
		parsedIncrementLimit, err := strconv.ParseInt(incrementLimitString, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		functionIncrementLimits = append(functionIncrementLimits, parsedIncrementLimit)
	}
	return functionIncrementLimits
}
