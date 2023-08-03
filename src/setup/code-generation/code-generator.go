package code_generation

import log "github.com/sirupsen/logrus"

func GenerateCode(functionName string, provider string) {

	switch provider {
	case "aws":
		fallthrough
	case "azure":
		fallthrough
	case "google":
		fallthrough
	default:
		log.Warnf("Code generation of %s function for %s is not supported.")
	}
}
