package setup

import log "github.com/sirupsen/logrus"

type Serverless struct {
	// TODO: add serverless.yml fields
}

// AddFunctionConfig Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfig(subex *SubExperiment, index int) {
	// TODO: implement function
	log.Warnf("Adding function configg of Subexperiment %s, index %d", subex.Function, index)
}

// CreateServerlessConfigFile dumps the contents of the Serverless struct into a yml file.
func (s *Serverless) CreateServerlessConfigFile() {
	// TODO: implement function
}

// RemoveService removes the service defined in serverless.yml
func RemoveService() string {
	// TODO: implement function
	return ""
}

// Deploys the functions defined in the serverless.com file

func DeployService() string {
	// TODO: implement function
	return ""
}
