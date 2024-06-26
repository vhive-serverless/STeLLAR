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

package benchhttp

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

const (
	timeout = 15 * time.Minute
)

// ExecuteRequest will send an HTTP request, check its status code and return the response body.
func ExecuteRequest(req http.Request) (bool, []byte, time.Time, time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ok := true
	defer cancel()

	err, resp, reqSentTime, reqReceivedTime := sendTimedRequest(ctx, req)
	if err != nil {
		ok = false
		log.Errorf("Could not send HTTP request: %s", err.Error())
		return ok, nil, reqSentTime, reqReceivedTime
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ok = false
		log.Errorf("Could not read HTTP response body: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		ok = false
		log.Errorf("Response from %s had status %s: %s", req.URL.Hostname(), resp.Status, string(bodyBytes))
	}

	return ok, bodyBytes, reqSentTime, reqReceivedTime
}

// https://stackoverflow.com/questions/48077098/getting-ttfb-time-to-first-byte-value-in-golang/48077762#48077762
func sendTimedRequest(ctx context.Context, req http.Request) (error, *http.Response, time.Time, time.Time) {
	var receivedFirstByte time.Time

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			receivedFirstByte = time.Now()
		},
	}

	reqSentTime := time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req.WithContext(httptrace.WithClientTrace(ctx, trace)))

	// For total time, return resp, reqSentTime, time.Now()
	return err, resp, reqSentTime, receivedFirstByte
}
