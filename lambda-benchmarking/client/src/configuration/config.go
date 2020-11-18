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
	Title                   string   `json:"Title"`
	Bursts                  int      `json:"Bursts"`
	BurstSizes              []int    `json:"BurstSizes"`
	PayloadLengthBytes      int      `json:"PayloadLengthBytes"`
	CooldownSeconds         float64  `json:"CooldownSeconds"`
	FunctionIncrementLimits []int64  `json:"FunctionIncrementLimits"`
	DesiredServiceTimes     []string `json:"DesiredServiceTimes"`
	IATType                 string   `json:"IATType"`
	Provider                string   `json:"Provider"`
	GatewaysNumber          int      `json:"GatewaysNumber"`
	Visualization           string   `json:"Visualization"`
	FunctionMemoryMB        int64    `json:"FunctionMemoryMB"`
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

	standardIncrement := int64(1e10)
	standardDurationMs := timeSession(standardIncrement).Milliseconds()
	cachedServiceTimeIncrement = make(map[string]int64)
	for subExperimentIndex := range parsedConfiguration.SubExperiments {
		determineFunctionIncrementLimits(&parsedConfiguration.SubExperiments[subExperimentIndex],
			standardIncrement, standardDurationMs)
	}

	setDefaults(parsedConfiguration.SubExperiments)
	return parsedConfiguration
}

const defaultVisualization = "cdf"
const defaultIATType = "stochastic"
const defaultProvider = "aws"
const defaultFunctionMemoryMB = 128
const defaultGatewaysNumber = 1

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
		if parsedSubExps[index].FunctionMemoryMB == 0 {
			parsedSubExps[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedSubExps[index].GatewaysNumber == 0 {
			parsedSubExps[index].GatewaysNumber = defaultGatewaysNumber
		}
	}
}
