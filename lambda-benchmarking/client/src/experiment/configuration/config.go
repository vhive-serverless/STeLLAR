package configuration

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type SubExperiment struct {
	Title                   string  `json:"Title"`
	Bursts                  int     `json:"Bursts"`
	BurstSizes              []int   `json:"BurstSizes"`
	PayloadLengthBytes      int     `json:"PayloadLengthBytes"`
	CooldownSeconds         float64 `json:"CooldownSeconds"`
	FunctionIncrementLimits []int64 `json:"FunctionIncrementLimits"`
	IATType                 string  `json:"IATType"`
	Provider                string  `json:"Provider"`
	GatewaysNumber          int     `json:"GatewaysNumber"`
	Visualization           string  `json:"Visualization"`
	GatewayEndpoints        []string
	Id                      int
}

func Extract(configFile *os.File) []SubExperiment {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfigs []SubExperiment
	if err := json.Unmarshal(configByteValue, &parsedConfigs); err != nil {
		log.Error(err)
	}

	return parsedConfigs
}
