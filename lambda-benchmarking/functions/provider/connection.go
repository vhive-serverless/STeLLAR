package provider

import (
	"functions/provider/amazon"
	"functions/writer"
	log "github.com/sirupsen/logrus"
)

//Connection creates an interface through which to interact with various providers
type Connection struct {
	ProviderName string
}

//DeployFunction will create a new serverless function in the specified language, with id `i`. An API for it will
//then be created, as well as corresponding interactions between them and specific permissions.
func (c Connection) DeployFunction(i int, language string, zipLocation string) {
	switch c.ProviderName {
	case "aws":
		awsInterface := amazon.Initialize()

		apiID, memoryAssigned := awsInterface.DeployFunction(i, language, zipLocation)
		writer.GatewaysWriterSingleton.WriteGatewayID(apiID, memoryAssigned)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

//RemoveFunction will remove the serverless function with id `i`.
func (c Connection) RemoveFunction(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := amazon.Initialize()

		awsInterface.RemoveFunction(i)
		awsInterface.RemoveAPI(i)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

//UpdateFunction will update the source code of the serverless function with id `i`.
func (c Connection) UpdateFunction(i int, zipLocation string) {
	switch c.ProviderName {
	case "aws":
		awsInterface := amazon.Initialize()

		awsInterface.UpdateFunction(i, zipLocation)
		apiID, memoryAssigned := awsInterface.UpdateFunctionConfiguration(i)
		writer.GatewaysWriterSingleton.WriteGatewayID(apiID, memoryAssigned)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}