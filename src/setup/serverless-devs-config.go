package setup

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerlessDevs struct {
	edition  string
	name     string
	access   string
	services map[string]*Service
}

type Service struct {
	component string
	props     Props
}

type Props struct {
	region        string
	service       ServiceDetails
	function      ServerlessDevsFunction
	triggers      []Trigger
	customDomains []CustomDomain
}

type ServiceDetails struct {
	name        string
	description string
}

type ServerlessDevsFunction struct {
	name        string
	description string
	runtime     string
	codeUri     string
	handler     string
	memorySize  string
	timeout     string
}

type Trigger struct {
	name   string
	type_  string
	config TriggerConfig
}

type TriggerConfig struct {
	authType string
	methods  []string
}

type CustomDomain struct {
	domainName  string
	protocol    string
	routeConfig []RouteConfig
}

type RouteConfig struct {
	path    string
	methods []string
}

func (s *ServerlessDevs) CreateServerlessDevsConfigFile(path string) {
	data, yamlMarshalErr := yaml.Marshal(&s)
	if yamlMarshalErr != nil {
		log.Fatal(yamlMarshalErr)
	}

	writeFileErr := os.WriteFile(path, data, 0644)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
}
