package configuration

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

//Configuration is the schema for all experiment configurations.
type Configuration struct {
	Sequential     bool            `json:"Sequential"`
	SubExperiments []SubExperiment `json:"SubExperiments"`
}

//SubExperiment is the schema for sub-experiment configurations.
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
	ID                      int
}

//Extract will read the given JSON configuration file and load it as an array of sub-experiment configurations.
func Extract(configFile *os.File) Configuration {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfiguration Configuration
	if err := json.Unmarshal(configByteValue, &parsedConfiguration); err != nil {
		log.Fatalf("Could not extract configuration from file: %s", err.Error())
	}

	setDefaults(parsedConfiguration.SubExperiments)
	return parsedConfiguration
}

const defaultVisualization = "all-light"
const defaultIATType = "stochastic"
const defaultProvider = "aws"

func setDefaults(parsedSubExps []SubExperiment) {
	for index := range parsedSubExps {
		if parsedSubExps[index].Visualization == "" {
			parsedSubExps[index].Visualization = defaultVisualization
		}
		if parsedSubExps[index].IATType == "" {
			parsedSubExps[index].IATType = defaultIATType
		}
		if parsedSubExps[index].Provider == "" {
			parsedSubExps[index].Provider = defaultProvider
		}
	}
}
