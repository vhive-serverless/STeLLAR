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
	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (instance awsSingleton) ListFunctions(marker *string) []*lambda.FunctionConfiguration {
	log.Info("Querying Lambda functions...")
	var queriedFunctions []*lambda.FunctionConfiguration

	args := &lambda.ListFunctionsInput{
		Marker: marker,
	}

	result, err := instance.lambdaSvc.ListFunctions(args)
	if err != nil {
		if strings.Contains(err.Error(), "TooManyRequestsException") {
			log.Warnf("Facing AWS rate-limiting error, retrying...")
			return instance.ListFunctions(nil)
		}

		log.Fatalf("Cannot list Lambda functions: %s", err.Error())
	}

	if result == nil {
		log.Fatalf("List Lambda functions result was nil.")
		return nil
	}

	queriedFunctions = append(queriedFunctions, result.Functions...)
	if (*result).NextMarker != nil {
		queriedFunctions = append(queriedFunctions, instance.ListFunctions((*result).NextMarker)...)
	}

	return queriedFunctions
}
