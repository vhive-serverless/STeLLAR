// MIT License
//
// Copyright (c) 2021 Theodor Amariucai
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

package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
)

var minioClientSingleton *minio.Client

func authenticateStorageClient(useS3 bool) *minio.Client {
	if minioClientSingleton != nil {
		return minioClientSingleton
	}

	var serverAddress, accessKey, secretKey string
	if useS3 {
		serverAddress = "s3.amazonaws.com"
		accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	} else {
		serverAddress = "http://10.96.0.46:9000"
		accessKey = "minio"
		secretKey = "minio123"
	}

	minioClient, err := minio.New(serverAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalf("Could not create minio client: %s", err.Error())
	}

	minioClientSingleton = minioClient
	return minioClientSingleton
}

func usingStorage(requestGRPC *proto_gen.InvokeChainRequest, requestHTTP *events.APIGatewayProxyRequest) bool {
	if requestHTTP != nil {
		value, hasUseS3Field := requestHTTP.QueryStringParameters["UseS3"]

		if !hasUseS3Field {
			return false
		}

		useS3, err := strconv.ParseBool(value)
		if err != nil {
			log.Errorf("Could not parse UseS3: %s", err.Error())
		}

		return useS3
	}

	// gRPC
	return requestGRPC.UseS3 == true
}

func loadObjectFromStorage(requestHTTP *events.APIGatewayProxyRequest, requestGRPC *proto_gen.InvokeChainRequest) string {
	var s3key, s3bucket string
	if requestHTTP != nil {
		s3key = requestHTTP.QueryStringParameters["S3Key"]
		s3bucket = requestHTTP.QueryStringParameters["S3Bucket"]
	} else {
		s3key = requestGRPC.S3Key
		s3bucket = requestGRPC.S3Bucket
	}

	s3Client := authenticateStorageClient(true) // always use S3 for now
	object, err := s3Client.GetObject(
		context.Background(),
		s3bucket,
		s3key,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Infof("Object %q not found in S3 bucket %q: %s", s3key, s3bucket, err.Error())
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, object); err != nil {
		log.Infof("Error reading object body: %v", err)
		return ""
	}

	return buf.String()
}

func saveObjectToStorage(requestHTTP *events.APIGatewayProxyRequest, stringPayload string, requestGRPC *proto_gen.InvokeChainRequest) {
	if requestHTTP != nil {
		s3key := saveObject(stringPayload, requestHTTP.QueryStringParameters["S3Bucket"])
		requestHTTP.QueryStringParameters["S3Key"] = s3key
	} else {
		s3key := saveObject(stringPayload, requestGRPC.S3Bucket)
		requestGRPC.S3Key = s3key
	}
}

func saveObject(payload string, s3bucket string) string {
	log.Infof(`Using S3, saving transfer payload (~%d bytes) to AWS S3.`, len(payload))
	s3key := fmt.Sprintf("transfer-payload-%s", randStringBytes(20))

	s3Client := authenticateStorageClient(true) // always use S3 for now

	uploadOutput, err := s3Client.PutObject(
		context.Background(),
		s3bucket,
		s3key,
		strings.NewReader(payload),
		-1,
		minio.PutObjectOptions{},
	)
	if err != nil {
		log.Fatalf("Unable to upload %q to %q, %v", s3key, s3bucket, err.Error())
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", s3key, s3bucket, uploadOutput.Location)
	return s3key
}
