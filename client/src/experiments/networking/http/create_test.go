package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestCreateAWSRequest(t *testing.T) {
	randomPayloadLength := 7
	randomEndpointID := "uicnaywo3rb3nsci"
	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("aws", randomPayloadLength, randomEndpointID, randomAssignedIncrement)

	expectedHostname := fmt.Sprintf("%s.execute-api.%s.amazonaws.com", randomEndpointID, awsRegion)
	equals(t, expectedHostname, req.Host)
	equals(t, expectedHostname, req.URL.Host)
	equals(t, http.MethodPost, req.Method)
	equals(t, "https", req.URL.Scheme)
}

func TestCreateExternalRequest(t *testing.T) {
	randomPayloadLength := 7
	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("www.google.com", randomPayloadLength, "", randomAssignedIncrement)

	equals(t, "www.google.com", req.Host)
	equals(t, "www.google.com", req.URL.Host)
	equals(t, http.MethodGet, req.Method)
	equals(t, "https", req.URL.Scheme)
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
