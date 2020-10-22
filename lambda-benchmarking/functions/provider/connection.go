package provider

import (
	"functions/provider/aws"
	"functions/writer"
	"log"
)

type Connection struct {
	ProviderName string
}

func (p Connection) DeployFunction(i int, language string) {
	switch p.ProviderName {
	case "aws":
		awsInterface := aws.Interface{}
		apiID := awsInterface.DeployFunction(i, language)
		writer.GatewaysWriterSingleton.WriteRowToFile(apiID)
	default:
		log.Fatalf("Unrecognized provider %s", p.ProviderName)
	}
}

func (p Connection) RemoveFunction(i int) {
	switch p.ProviderName {
	case "aws":
		awsInterface := aws.Interface{}
		awsInterface.RemoveFunction(i)
		apiID := awsInterface.GetAPIID(i)
		awsInterface.RemoveAPI(i, apiID)
	default:
		log.Fatalf("Unrecognized provider %s", p.ProviderName)
	}
}

func (p Connection) UpdateFunction(i int) {
	switch p.ProviderName {
	case "aws":
		awsInterface := aws.Interface{}
		awsInterface.UpdateFunction(i)
		awsInterface.UpdateFunctionConfiguration(i)
	default:
		log.Fatalf("Unrecognized provider %s", p.ProviderName)
	}
}
