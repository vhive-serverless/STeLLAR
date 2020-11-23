// MIT License
//
// Copyright (c) 2020 Theodor Amariucai
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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

const (
	awsRegion = "us-west-1"
	maxAPIs   = 600
)

//AWSSingleton is an object used to interact with AWS through the methods it exports.
var AWSSingleton *instance

type instance struct {
	//There can only be one of localZip vs s3Bucket, s3Key
	localZip             []byte
	s3Bucket             string
	s3Key                string
	APINamePrefix        string
	LambdaFunctionPrefix string
	region               string
	cloneAPIID           string
	stage                string
	session              *session.Session
	s3Svc                *s3.S3
	lambdaSvc            *lambda.Lambda
	apiGatewaySvc        *apigateway.APIGateway
}

//InitializeSingleton will create a new Amazon instance to interact with different AWS services.
func InitializeSingleton() {
	sessionInstance := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	AWSSingleton = &instance{
		s3Bucket:             "benchmarking-aws",
		APINamePrefix:        "vHiveAPI_",
		LambdaFunctionPrefix: "vHiveFunction_",
		region:               awsRegion,
		cloneAPIID:           "hjnwqihyo1",
		stage:                "prod",
		session:              sessionInstance,
		lambdaSvc:            lambda.New(sessionInstance),
		apiGatewaySvc:        apigateway.New(sessionInstance),
		s3Svc:                s3.New(sessionInstance),
	}
}

//UploadZIPToS3 helps get around the 50MB image size limit for AWS functions.
func UploadZIPToS3(localZipPath string, sizeMB float64) {
	log.Infof(`Deploying to AWS and package size (~%vMB) > 50 MB, will now attempt to upload to Amazon S3.`, sizeMB)
	AWSSingleton.s3Key = fmt.Sprintf("benchmarking%vMB.zip", sizeMB)

	//TODO: test this
	_, err := AWSSingleton.s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(AWSSingleton.s3Bucket),
		Key:    aws.String(AWSSingleton.s3Key),
	})
	if err == nil {
		log.Infof("%q was already found in bucket %q.", AWSSingleton.s3Key, AWSSingleton.s3Bucket)
		return
	}

	zipFile, err := os.Open(localZipPath)
	if err != nil {
		log.Fatalf("Failed to open zip file %q: %v", localZipPath, err)
	}

	uploadOutput, err := s3manager.NewUploader(AWSSingleton.session).Upload(&s3manager.UploadInput{
		Bucket: aws.String("benchmarking-aws"),
		Key:    aws.String(AWSSingleton.s3Key),
		Body:   zipFile,
	})
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", AWSSingleton.s3Key, AWSSingleton.s3Bucket, err.Error())
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", AWSSingleton.s3Key, AWSSingleton.s3Bucket, uploadOutput.Location)
}

//SetLocalZip sets the location of the zipped binary file for the function to be deployed.
func SetLocalZip(path string) {
	zipBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read local zipped binary: %s", err.Error())
	}
	AWSSingleton.localZip = zipBytes
}
