package setup

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func extractConfiguration(configFile *os.File) Configuration {
	configByteValue, _ := ioutil.ReadAll(configFile)

	var parsedConfig Configuration
	if err := json.Unmarshal(configByteValue, &parsedConfig); err != nil {
		log.Fatalf("Could not extract experiment configuration from file: %s", err.Error())
	}

	if parsedConfig.Provider == "" {
		parsedConfig.Provider = defaultProvider
	}
	if parsedConfig.Runtime == "" {
		parsedConfig.Runtime = defaultRuntime
	}

	for index := range parsedConfig.SubExperiments {
		if parsedConfig.SubExperiments[index].Visualization == "" {
			parsedConfig.SubExperiments[index].Visualization = defaultVisualization
		}
		if parsedConfig.SubExperiments[index].IATType == "" {
			parsedConfig.SubExperiments[index].IATType = defaultIATType
		}
		if parsedConfig.SubExperiments[index].FunctionMemoryMB == 0 {
			parsedConfig.SubExperiments[index].FunctionMemoryMB = defaultFunctionMemoryMB
		}
		if parsedConfig.SubExperiments[index].GatewaysNumber == 0 {
			parsedConfig.SubExperiments[index].GatewaysNumber = defaultGatewaysNumber
		}
	}

	log.Debugf("Extracted %d sub-experiments from given configuration file.", len(parsedConfig.SubExperiments))
	return parsedConfig
}