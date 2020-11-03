package networking

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const Region = "us-west-1"

var signerSingleton *v4.Signer

func GetAWSSignerSingleton() *v4.Signer {
	if signerSingleton != nil {
		return signerSingleton
	}

	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			checkAndReturnEnvVar("AWS_ACCESS_KEY_ID"),
			checkAndReturnEnvVar("AWS_SECRET_ACCESS_KEY"),
			""),
		Region: aws.String(Region),
	}))
	signerSingleton = v4.NewSigner(sessionInstance.Config.Credentials)
	return signerSingleton
}

func checkAndReturnEnvVar(key string) string {
	envVar, isSet := os.LookupEnv(key)
	if !isSet {
		log.Warnf("Environment variable %s is not set.", key)
	}
	return envVar
}

type LambdaFunctionResponse struct {
	AwsRequestID string `json:"AwsRequestID"`
	Payload      []byte `json:"Payload"`
}

func GetAWSRequestID(resp *http.Response) string {
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var lambdaFunctionResponse LambdaFunctionResponse
	if err := json.Unmarshal(bytes, &lambdaFunctionResponse); err != nil {
		log.Error(err)
	}
	return lambdaFunctionResponse.AwsRequestID
}
