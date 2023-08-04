package building

import log "github.com/sirupsen/logrus"

// Builder struct keeps track of the functions built by stellar
type Builder struct {
	functionsBuilt []string
}

func (b *Builder) BuildFunction(functionPath string, runtime string) {
	// TODO: Implement function

	// First we check whether the function has not been built already
	// TODO: Check if function path is in functionsBuilt if yes, skip the build. If no, continue the building process and add the functionPath to the list.

	b.functionsBuilt = append(b.functionsBuilt, functionPath)
	switch runtime {
	case "java":
		buildJava(functionPath)
	case "golang":
		buildGolang(functionPath)
	default:
		// building not supported
		log.Warnf("Building runtime %s is not necessary, or not supported. Continuing without building.", runtime)
	}
}

// buildJava builds the java zip artifact for serverless deployment using Gradle
func buildJava(functionPath string) {
	// TODO: Implement function.
	log.Warn(functionPath)
}

// buildGolang builds the Golang binary for serverless deployment
func buildGolang(functionPath string) {
	// TODO: Implement function.
	log.Warn(functionPath)
}
