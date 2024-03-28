package code_generation

import log "github.com/sirupsen/logrus"

// GenerateCode generates source code of the given function for the given provider
func GenerateCode(functionName string, provider string) {

	//TODO: Implement function.

	switch provider {
	case "aws":
		fallthrough
	case "azure":
		fallthrough
	case "google":
		fallthrough
	default:
		log.Warnf("Code generation of %s function for %s is not supported.", functionName, provider)
	}
}
