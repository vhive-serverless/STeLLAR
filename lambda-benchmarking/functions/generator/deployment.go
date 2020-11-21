package generator

import (
	"functions/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
)

//SetupDeployment will create the serverless function zip deployment for the given provider,
//in the given language and of the given size in bytes.
func SetupDeployment(action string, language string, provider string, sizeBytes int) string {
	switch action {
	case "deploy":
		fallthrough
	case "update":
		generateDeploymentPackage(provider, language, sizeBytes)
	case "remove":
		// No setup required for removing functions
	default:
		log.Fatalf("Unrecognized function action %s", action)
	}
	return ""
}

func generateDeploymentPackage(provider string, language string, sizeBytes int) {
	zippedBinarySize := createBinary(provider, language)

	if sizeBytes < zippedBinarySize {
		log.Fatalf("Total size (~%dMB) cannot be smaller than zipped binary size (~%dMB).",
			util.BytesToMB(sizeBytes),
			util.BytesToMB(zippedBinarySize))
	}

	randomFileName := "random.file"
	generateRandomFile(sizeBytes-zippedBinarySize, randomFileName)
	generateZIP(provider, randomFileName, sizeBytes)
}

func generateRandomFile(sizeBytes int, randomFileName string) {
	if fileExists(randomFileName) {
		log.Infof("Random file `%s` already exists, removing...", randomFileName)
		if err := os.Remove(randomFileName); err != nil {
			log.Fatalf("Failed to remove random file `%s`", randomFileName)
		}
	}

	buffer := make([]byte, sizeBytes)
	_, err := rand.Read(buffer) // The slice should now contain random bytes instead of only zeroes (prevents efficient archiving).
	if err != nil {
		log.Fatalf("Failed to fill buffer with random bytes: `%s`", err.Error())
	}

	if err := ioutil.WriteFile(randomFileName, buffer, 0666); err != nil {
		log.Fatalf("Could not generate random file with size %d bytes", sizeBytes)
	}
}
