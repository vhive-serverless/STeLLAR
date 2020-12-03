package http

import (
	"testing"
)

func TestExecuteExternalHTTPRequest(t *testing.T) {
	randomPayloadLength := 7
	randomAssignedIncrement := int64(1482911482)
	req := CreateRequest("www.google.com", randomPayloadLength, "", randomAssignedIncrement)

	respBytes, reqSentTime, reqReceivedTime := ExecuteHTTPRequest(*req)
	equals(t, true, respBytes != nil)
	equals(t, true, reqReceivedTime.Sub(reqSentTime) > 0)
}