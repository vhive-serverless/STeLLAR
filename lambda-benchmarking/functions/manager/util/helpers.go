package util

import (
	"log"
	"os"
)

func CheckAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Fatalf("Environment variable %s is not set.", key)
	}
	return envVar
}