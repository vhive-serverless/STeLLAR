package provider

import (
	"functions/provider/aws"
	"functions/writer"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	ProviderName string
}

func (c Connection) DeployFunction(i int, language string) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		apiID := awsInterface.DeployFunction(i, language)
		writer.GatewaysWriterSingleton.WriteRowToFile(apiID)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

func (c Connection) RemoveFunction(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		awsInterface.RemoveFunction(i)
		apiID := awsInterface.GetAPIID(i)
		awsInterface.RemoveAPI(i, apiID)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}

func (c Connection) UpdateFunction(i int) {
	switch c.ProviderName {
	case "aws":
		awsInterface := aws.Initialize()

		awsInterface.UpdateFunction(i)
		awsInterface.UpdateFunctionConfiguration(i)
	default:
		log.Fatalf("Unrecognized provider %s", c.ProviderName)
	}
}
