package setup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os/exec"
	"regexp"
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
	Handler     string  `yaml:"handler"`
	Description string  `yaml:"description"`
	Name        string  `yaml:"name"`
	Events      []Event `yaml:"events"`
}

type Event struct {
	HttpApi HttpApi `yaml:"httpApi"`
}

type HttpApi struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

// CreateHeader sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeader(config Configuration) {
	s.Service = "STeLLAR" // or some other string
	s.FrameworkVersion = "3"
	s.Provider = Provider{
		Name:    config.Provider,
		Runtime: "go1.x",
		Region:  "us-west-2",
	}
	s.Functions = map[string]*Function{}
}

// AddPackagePattern adds a string pattern to Package.Pattern as long as such a pattern does not already exist in Package.Pattern
func (s *Serverless) AddPackagePattern(pattern string) {
	if !util.StringContains(s.Package.Patterns, pattern) {
		s.Package.Patterns = append(s.Package.Patterns, pattern)
	}
}

// CreateServerlessConfigFile dumps the contents of the Serverless struct into a yml file.
func (s *Serverless) CreateServerlessConfigFile() {
	data, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile("src/setup/deployment/raw-code/serverless/aws/serverless.yml", data, 0644)

	if err2 != nil {
		log.Fatal(err2)
	}
}

// AddFunctionConfig creates parallel serverless function configuration for serverless.com deployment
func (s *Serverless) AddFunctionConfig(subex *SubExperiment, index int) {
	// serverless.com functions require alphanumeric names
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

	for i := 0; i < subex.Parallelism; i++ {
		name := fmt.Sprintf("%s_%d_%d", nonAlphanumericRegex.ReplaceAllString(subex.Title, ""), index, i)

		events := []Event{{HttpApi{Path: "/" + name, Method: "GET"}}}

		var handler string
		switch subex.Function {
		case "hellopy":
			handler = "hellopy/lambda_function.lambda_handler"
			s.AddPackagePattern("hellopy/lambda_function.py")
			break
		case "producer-consumer":
			handler = "hellogo"
			s.AddPackagePattern("hellogo/*")
		default:
			log.Fatalf("DeployFunction could not recognize function image %s", subex.Function)
		}

		f := &Function{Handler: handler, Name: name, Events: events}
		s.Functions[name] = f
		subex.addRoute(name)
	}

}

// RemoveService removes the service defined in serverless.yml
func RemoveService() string {
	slsRemoveCmd := exec.Command("sls", "remove")
	slsRemoveCmd.Dir = "src/setup/deployment/raw-code/serverless/aws"
	slsRemoveMessage := util.RunCommandAndLog(slsRemoveCmd)
	return slsRemoveMessage
}

// Deploys the functions defined in the serverless.com file
func deployService() string {
	slsDeployCmd := exec.Command("sls", "deploy")
	slsDeployCmd.Dir = "src/setup/deployment/raw-code/serverless/aws"
	slsDeployMessage := util.RunCommandAndLog(slsDeployCmd)
	return slsDeployMessage
}
