// MIT License
//
// Copyright (c) 2021 Theodor Amariucai and EASE Lab, Michal Baczun
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

package p

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func isUsingStorage(requestGRPC *InvokeChainRequest, requestHTTP *http.Request) bool {
	if requestHTTP != nil {
		value := requestHTTP.URL.Query().Get("StorageTransfer")

		if len(value) == 0 {
			return false
		}

		storageTransfer, err := strconv.ParseBool(value)
		if err != nil {
			log.Errorf("Could not parse StorageTransfer: %s", err.Error())
		}

		return storageTransfer
	}

	// gRPC
	return requestGRPC.StorageTransfer == true
}

func loadObjectFromStorage(requestHTTP *http.Request, requestGRPC *InvokeChainRequest) string {
	var objectKey, objectBucket string
	if requestHTTP != nil {
		objectKey = requestHTTP.URL.Query().Get("Key")
		objectBucket = requestHTTP.URL.Query().Get("Bucket")
	} else {
		objectKey = requestGRPC.Key
		objectBucket = requestGRPC.Bucket
	}

	if requestHTTP != nil {
		client, ctx := authenticateCloudStorageClient()

		bucket := client.Bucket(objectBucket)
		object := bucket.Object(objectKey)

		reader, err := object.NewReader(ctx)
		if err != nil {
			log.Infof("Failed to obtain reader for object %v: %v", objectKey, err)
			return ""
		}
		defer reader.Close()

		payload, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Infof("Error reading object body: %v", err)
			return ""
		}

		return string(payload)
	}

	// when using anything but HTTP, at the moment, automatically resort to minio
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	storageClient := getMinioClient()
	object, err := storageClient.GetObject(
		ctx,
		objectBucket,
		objectKey,
		minio.GetObjectOptions{},
	)
	if err != nil {
		log.Infof("Object %q not found in bucket %q: %s", objectKey, objectBucket, err.Error())
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, object); err != nil {
		log.Infof("Error reading object body: %v", err)
		return ""
	}

	return buf.String()
}

func saveObjectToStorage(requestHTTP *http.Request, stringPayload string, requestGRPC *InvokeChainRequest) {
	if requestHTTP != nil {
		key := saveObject(stringPayload, requestHTTP.URL.Query().Get("Bucket"), false)
		requestHTTP.URL.RawQuery += fmt.Sprintf("&Key=%v", key)
	} else {
		// when using anything but HTTP, at the moment, automatically resort to minio
		key := saveObject(stringPayload, requestGRPC.Bucket, true)
		requestGRPC.Key = key
	}
}

func saveObject(payload string, bucket string, useMinio bool) string {
	key := fmt.Sprintf("transfer-payload-%dbytes-%s", int64(len(payload)), GeneratePayloadFromGlobalRandom(20))
	log.Infof(`Using storage, saving transfer payload (~%d bytes) as %q to %q.`, len(payload), key, bucket)

	var uploadResult string
	if useMinio {
		storageClient := getMinioClient()

		uploadOutput, err := storageClient.PutObject(
			context.Background(),
			bucket,
			key,
			strings.NewReader(payload),
			int64(len(payload)),
			minio.PutObjectOptions{},
		)
		if err != nil {
			log.Fatalf("Unable to upload %q to %q, %v", key, bucket, err.Error())
		}
		uploadResult = uploadOutput.Location
	} else {
		client, ctx := authenticateCloudStorageClient()

		bucketVar := client.Bucket(bucket)
		object := bucketVar.Object(key)
		w := object.NewWriter(ctx)
		if _, err := fmt.Fprintf(w, payload); err != nil {
			log.Fatalf("Unable to upload %q to %q, %v", key, bucket, err.Error())
		}
		if err := w.Close(); err != nil {
			log.Fatalf("Error in closing google object writer: %v", err)
		}
		// No access to anything similar to s3 uploadOutput.Location
		uploadResult = "GCS location unavailable"
	}

	log.Infof("Successfully uploaded %q to bucket %q (%s)", key, bucket, uploadResult)
	return key
}
