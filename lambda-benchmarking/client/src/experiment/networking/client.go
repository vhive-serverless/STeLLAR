package networking

import (
	"context"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	timeout = 15 * time.Minute
)

func MakeHTTPRequest(req http.Request) *http.Response {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		log.Errorf("Response from %s had status %s:\n %s", req.URL.Hostname(), resp.Status, string(bodyBytes))
	}

	return resp
}
