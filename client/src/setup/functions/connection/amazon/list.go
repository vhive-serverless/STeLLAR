package amazon

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (amazon instance) ListFunctions() []*lambda.FunctionConfiguration {
	log.Info("Querying Lambda functions...")

	args := &lambda.ListFunctionsInput{
		MaxItems: aws.Int64(maxAPIs),
	}

	result, err := amazon.lambdaSvc.ListFunctions(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return amazon.ListFunctions()
		}

		log.Fatalf("Cannot list Lambda functions: %s", err.Error())
	}
	log.Debugf("List Lambda functions result: %s", result.String())

	log.Infof("Found %d Lambda functions.", len(result.Functions))
	return result.Functions
}
