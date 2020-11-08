package networking

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	log "github.com/sirupsen/logrus"
)

const awsRegion = "us-west-1"

var signerSingleton *v4.Signer

func getAWSSignerSingleton() *v4.Signer {
	if signerSingleton != nil {
		return signerSingleton
	}

	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	signerSingleton = v4.NewSigner(sessionInstance.Config.Credentials)
	return signerSingleton
}

type lambdaFunctionResponse struct {
	AwsRequestID string `json:"AwsRequestID"`
	Payload      []byte `json:"Payload"`
}

//GetAWSRequestID will process an HTTP response body coming from an AWS integration, extracting its ID.
func GetAWSRequestID(respBody []byte) string {
	var lambdaFunctionResponse lambdaFunctionResponse
	if err := json.Unmarshal(respBody, &lambdaFunctionResponse); err != nil {
		log.Error(err)
	}
	return lambdaFunctionResponse.AwsRequestID
}
