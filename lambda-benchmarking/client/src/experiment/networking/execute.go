package networking

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"time"
)

const (
	timeout = 15 * time.Minute
)

//ExecuteHTTPRequest will send an HTTP request, check its status code and return the response body.
func ExecuteHTTPRequest(req http.Request) ([]byte, time.Time, time.Time) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	resp, reqSentTime, reqReceivedTime := sendTimedRequest(ctx, req)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Could not read HTTP response body: %s", err.Error())
		}
		log.Errorf("Response from %s had status %s:\n %s", req.URL.Hostname(), resp.Status, string(bodyBytes))
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	return bytes, reqSentTime, reqReceivedTime
}

//https://stackoverflow.com/questions/48077098/getting-ttfb-time-to-first-byte-value-in-golang/48077762#48077762
func sendTimedRequest(ctx context.Context, req http.Request) (*http.Response, time.Time, time.Time) {
	var receivedFirstByte time.Time

	trace := &httptrace.ClientTrace{
		GotFirstResponseByte: func() {
			receivedFirstByte = time.Now()
		},
	}

	reqSendTime := time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req.WithContext(httptrace.WithClientTrace(ctx, trace)))
	if err != nil {
		log.Fatalf("HTTP request failed with error %s", err.Error())
	}
	// For total time, return resp, reqSendTime, time.Now()
	return resp, reqSendTime, receivedFirstByte
}
