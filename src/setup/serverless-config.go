package setup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"stellar/setup/deployment/connection/amazon"
	"stellar/util"
)

// Serverless describes the serverless.yml contents.
type Serverless struct {
	Service          string               `yaml:"service"`
	FrameworkVersion string               `yaml:"frameworkVersion"`
	Provider         Provider             `yaml:"provider"`
	Package          Package              `yaml:"package"`
	Functions        map[string]*Function `yaml:"functions"`
}

type Provider struct {
	Name    string `yaml:"name"`
	Runtime string `yaml:"runtime"`
	Region  string `yaml:"region"`
}

type Package struct {
	Patterns []string `yaml:"patterns"`
}

type Function struct {
	Handler string  `yaml:"handler"`
	Runtime string  `yaml:"runtime"`
	Name    string  `yaml:"name"`
	Events  []Event `yaml:"events"`
}

type Event struct {
	HttpApi HttpApi `yaml:"httpApi"`
}

type HttpApi struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// CreateHeaderConfig sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeaderConfig(config *Configuration) {

	var region string
	switch config.Provider {
	case "aws":
		region = amazon.AWSRegion
	default:
		log.Errorf("Deployment to provider %s not supported yet.", config.Provider)
	}

	s.Service = "STeLLAR" // or some other string
	s.FrameworkVersion = "3"

	s.Provider = Provider{
		Name:    config.Provider,
		Runtime: config.Runtime,
		Region:  region,
	}
	s.addPackagePattern("!**")
}

// AddPackagePattern adds a string pattern to Package.Pattern as long as such a pattern does not already exist in Package.Pattern
func (s *Serverless) addPackagePattern(pattern string) {
	if !util.StringContains(s.Package.Patterns, pattern) {
		s.Package.Patterns = append(s.Package.Patterns, pattern)
	}
}

// AddFunctionConfig Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfig(subex *SubExperiment, index int) {
	log.Warnf("Adding function config of Subexperiment %s, index %d", subex.Function, index)
	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}
	for i := 0; i < subex.Parallelism; i++ {

		handler := subex.Handler
		runtime := subex.Runtime
		name := createName(subex, index, i)
		events := []Event{{HttpApi{Path: "/" + name, Method: "GET"}}}

		f := &Function{Handler: handler, Runtime: runtime, Name: name, Events: events}
		s.addPackagePattern(subex.PackagePattern)
		s.Functions[name] = f

		// TODO: producer-consumer sub-function definition
	}
}

// createName removes no-alphanumeric characters as serverless.com functions requires alphanumeric names. It also adds alphanumeric indexes to ensure a unique function name.
func createName(subex *SubExperiment, index int, parallelism int) string {
	return fmt.Sprintf("%s_%d_%d", nonAlphanumericRegex.ReplaceAllString(subex.Title, ""), index, parallelism)
}

// CreateServerlessConfigFile dumps the contents of the Serverless struct into a yml file.
func (s *Serverless) CreateServerlessConfigFile(path string) {
	data, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile(path, data, 0644)

	if err2 != nil {
		log.Fatal(err2)
	}
}

// RemoveService removes the service defined in serverless.yml
func RemoveService() string {
	// TODO: implement function
	return ""
}

// Deploys the functions defined in the serverless.com file

func DeployService() string {
	// TODO: implement function
	return ""
}
