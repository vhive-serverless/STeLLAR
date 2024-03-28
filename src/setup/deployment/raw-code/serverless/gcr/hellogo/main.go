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

type HelloGoResponse struct {
	RequestID      string   `json:"RequestID"`
	TimestampChain []string `json:"TimestampChain"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	incrementLimit, err := extractIncrementLimit(r)
	if err != nil {
		log.Errorf("Error extracting IncrementLimit: %s", err)
		http.Error(w, "Error extracting IncrementLimit", http.StatusBadRequest)
		return
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

func extractIncrementLimit(r *http.Request) (int, error) {
	incrementLimit := 0
	reqIncrementLimit := r.URL.Query().Get("IncrementLimit")
	if reqIncrementLimit != "" {
		var err error
		incrementLimit, err = strconv.Atoi(r.URL.Query().Get("IncrementLimit"))
		if err != nil {
			return 0, fmt.Errorf("Error parsing IncrementLimit: %s", err)
		}
	}
	return incrementLimit, nil
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

