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
	"encoding/json"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestGeneratePayload(t *testing.T) {
	emptyPayload := generateStringPayload("0")
	require.Equal(t, 0, len(emptyPayload))

	smallPayload := generateStringPayload("12")
	require.Equal(t, 12, len(smallPayload))

	mediumPayload := generateStringPayload("512")
	require.Equal(t, 512, len(mediumPayload))

	largePayload := generateStringPayload("1024")
	require.Equal(t, 1024, len(largePayload))
}

func TestExtractJSONTimestampChain(t *testing.T) {
	output, err := json.Marshal(ProducerConsumerResponse{
		RequestID:      "TestID",
		TimestampChain: []string{"1612371639523", "1612371639589"},
	})
	require.NoError(t, err)

	gatewayReply, err := json.Marshal(map[string]string{"body": string(output)})
	require.NoError(t, err)

	JSONTimestampChain := extractJSONTimestampChain(gatewayReply)
	require.Equal(t, 2, len(JSONTimestampChain))
	require.Equal(t, "1612371639523", JSONTimestampChain[0])
	require.Equal(t, "1612371639589", JSONTimestampChain[1])
}

func TestAppendTimestampToChain(t *testing.T) {
	updatedTimestampChain := AppendTimestampToChain([]string{})
	require.Equal(t, 1, len(updatedTimestampChain))
	timestamp1, err := strconv.Atoi(updatedTimestampChain[0])
	require.NoError(t, err)
	require.True(t, timestamp1 > 1612373289916) // 1612373289916 =  Wed Feb 03 2021 17:28:09

	time.Sleep(time.Millisecond)

	updatedTimestampChain = AppendTimestampToChain(updatedTimestampChain)
	require.Equal(t, 2, len(updatedTimestampChain))
	timestamp2, err := strconv.Atoi(updatedTimestampChain[1])
	require.NoError(t, err)
	require.Greater(t, timestamp2, timestamp1) // functions called later in the chain have a larger timestamp
}

func TestStringArrayToArrayOfString(t *testing.T) {
	require.Equal(t, StringArrayToArrayOfString("[]"), []string{""})
	require.Equal(t, StringArrayToArrayOfString("[35]"), []string{"35"})
	require.Equal(t, StringArrayToArrayOfString("[14 35 8]"), []string{"14", "35", "8"})
}
