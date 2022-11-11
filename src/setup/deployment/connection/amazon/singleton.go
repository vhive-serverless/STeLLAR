// MIT License
//
// Copyright (c) 2020 Theodor Amariucai and EASE Lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package amazon

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"stellar/util"
)

const (
	//AWSRegion is the region that AWS operates in
	AWSRegion = endpoints.UsWest1RegionID
	//AWSBucketName is the name of the bucket where the client operates
	AWSBucketName      = "stellar"
	deploymentStage    = "prod"
	maxFunctionTimeout = 900
	namingPrefix       = "vHive-bench_"
)

//AWSSingletonInstance is an object used to interact with AWS through the methods it exports.
var AWSSingletonInstance *awsSingleton

//UserARNNumber is used in AWS benchmarking for client authentication
var UserARNNumber string

type awsSingleton struct {
	// RequestSigner is the AWS object used for signing HTTP requests
	RequestSigner *v4.Signer
	// S3Key is the bucket location in which this specific deployment will be uploaded
	S3Key string
	// S3Bucket is the bucket in which this specific deployment will be uploaded
	S3Bucket string
	// ImageURI is the location where the docker image is located
	ImageURI                string
	s3Uploader              *s3manager.Uploader
	s3Svc                   *s3.S3
	lambdaSvc               *lambda.Lambda
	apiGatewaySvc           *apigateway.APIGateway
	ecrSvc                  *ecr.ECR
	apiTemplateFileContents []byte
	localZipFileContents    []byte
}

//InitializeSingleton will create a new Amazon awsSingleton to interact with different AWS services.
func InitializeSingleton(apiTemplatePath string) {
	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Region:                         aws.String(AWSRegion),
		CredentialsChainVerboseErrors:  aws.Bool(true),
		DisableRestProtocolURICleaning: aws.Bool(true),
	}))

	apiTemplateByteValue, err := io.ReadAll(util.ReadFile(apiTemplatePath))
	if err != nil {
		log.Fatalf("Could not read API template JSON when initializing AWS connection: %s", err.Error())
	}

	AWSSingletonInstance = &awsSingleton{
		RequestSigner:           v4.NewSigner(sessionInstance.Config.Credentials),
		lambdaSvc:               lambda.New(sessionInstance),
		apiGatewaySvc:           apigateway.New(sessionInstance),
		s3Svc:                   s3.New(sessionInstance),
		s3Uploader:              s3manager.NewUploader(sessionInstance),
		ecrSvc:                  ecr.New(sessionInstance),
		apiTemplateFileContents: apiTemplateByteValue,
		S3Bucket:                AWSBucketName,
	}
}

//UploadZIPToS3 helps get around the 50MB image size limit for AWS functions.
func UploadZIPToS3(localZipPath string, sizeMB float64) {
	log.Infof(`Deploying to AWS and package size (~%vMB) > 50 MB, will now attempt to upload to Amazon S3.`, sizeMB)
	AWSSingletonInstance.S3Key = fmt.Sprintf("benchmarking%vMB.zip", sizeMB)

	if _, err := AWSSingletonInstance.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(AWSSingletonInstance.S3Bucket),
		Key:    aws.String(AWSSingletonInstance.S3Key),
	}); err == nil {
		log.Infof("Object %q was already found in S3 bucket %q, skipping upload.", AWSSingletonInstance.S3Key, AWSSingletonInstance.S3Bucket)
		return
	}

	zipFile, err := os.Open(localZipPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %v", localZipPath, err)
	}

	uploadOutput, err := AWSSingletonInstance.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWSSingletonInstance.S3Bucket),
		Key:    aws.String(AWSSingletonInstance.S3Key),
		Body:   zipFile,
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", AWSSingletonInstance.S3Key, AWSSingletonInstance.S3Bucket, err.Error())
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", AWSSingletonInstance.S3Key, AWSSingletonInstance.S3Bucket, uploadOutput.Location)
}

//SetLocalZip sets the location of the zipped binary file for the function to be deployed.
func SetLocalZip(path string) {
	zipBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read local zipped binary: %s", err.Error())
	}
	AWSSingletonInstance.localZipFileContents = zipBytes
}

//GetECRAuthorizationToken helps the client get authorization for container AWS deployment.
func GetECRAuthorizationToken() string {
	log.Info("Requesting ECR authorization token.")

	result, err := AWSSingletonInstance.ecrSvc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return GetECRAuthorizationToken()
		}

		log.Fatalf("Cannot obtain ECR authorization token: %s", err.Error())
	}
	log.Debugf("Get ECR authorization token result: %s", result.String())

	authToken, err := base64.StdEncoding.DecodeString(*result.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		log.Fatalf("Could not decode base64-encoded ECR authorization token: %s", err.Error())
	}
	return strings.Split(string(authToken), ":")[1]
}
