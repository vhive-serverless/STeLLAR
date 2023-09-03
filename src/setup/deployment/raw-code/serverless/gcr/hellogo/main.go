package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HelloGoRequest struct {
	IncrementLimit int `json:"IncrementLimit"`
}

type HelloGoResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var reqBody HelloGoRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Errorf("Error decoding request body: %s", err)
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	incrementLimit := 0
	if reqBody.IncrementLimit > 0 {
		incrementLimit = reqBody.IncrementLimit
	}

	simulateWork(incrementLimit)

	res := HelloGoResponse{
		RequestID: "google-does-not-specify",
		TimestampChain: []string{
			strconv.Itoa(int(time.Now().Nanosecond())),
		},
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Errorf("Error encoding response body: %s", err)
		http.Error(w, "Error encoding response body", http.StatusInternalServerError)
		return
	}
}

func simulateWork(incrementLimit int) {
	log.Infof("Running function up to increment limit (%d)...", incrementLimit)
	for i := 0; i < incrementLimit; i++ {
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
