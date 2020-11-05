package provider

import (
	"functions/provider/aws"
	"functions/writer"
	log "github.com/sirupsen/logrus"
)

//Connection creates an interface through which to interact with various providers
type Connection struct {
	ProviderName string
}

//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (c Connection) DeployFunction(i int, language string) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		apiID := awsInterface.DeployFunction(i, language)
		writer.GatewaysWriterSingleton.WriteGatewayID(apiID)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

//RemoveFunction will remove the serverless function with id `i`.
func (c Connection) RemoveFunction(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		awsInterface.RemoveFunction(i)
		awsInterface.RemoveAPI(i)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

//UpdateFunction will update the source code of the serverless function with id `i`.
func (c Connection) UpdateFunction(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		awsInterface.UpdateFunction(i)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

//UpdateFunctionConfiguration  will update the configuration (e.g. timeout) of the serverless function with id `i`.
func (c Connection) UpdateFunctionConfiguration(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		awsInterface.UpdateFunctionConfiguration(i)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}
