package setup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"regexp"
	"stellar/setup/deployment/connection/amazon"
	"stellar/util"
	"strings"
)

// Serverless describes the serverless.yml contents.
type Serverless struct {
	Service          string               `yaml:"service"`
	FrameworkVersion string               `yaml:"frameworkVersion"`
	Provider         Provider             `yaml:"provider"`
	Package          Package              `yaml:"package,omitempty"`
	Plugins          []string             `yaml:"plugins,omitempty"`
	Functions        map[string]*Function `yaml:"functions"`
}

type Provider struct {
	Name           string      `yaml:"name"`
	Runtime        string      `yaml:"runtime"`
	Region         string      `yaml:"region"`
	FunctionApp    FunctionApp `yaml:"functionApp,omitempty"`
	SubscriptionId string      `yaml:"subscriptionId,omitempty"`
}

type FunctionApp struct {
	Version string `yaml:"version"`
}

type Package struct {
	Individually bool `yaml:"individually,omitempty"`
}

type Function struct {
	Handler string          `yaml:"handler"`
	Runtime string          `yaml:"runtime"`
	Name    string          `yaml:"name"`
	Events  []Event         `yaml:"events"`
	Package FunctionPackage `yaml:"package,omitempty"`
}

type FunctionPackage struct {
	Patterns []string `yaml:"patterns,omitempty"`
	Artifact string   `yaml:"artifact,omitempty"`
}

type Event struct {
	HttpApiAWS     HttpApi  `yaml:"httpApi,omitempty"`
	HttpAzure      bool     `yaml:"http"`
	MethodsAzure   []string `yaml:"methods"`
	AuthLevelAzure string   `yaml:"authLevel"`
}

type HttpAzure struct {
	HttpAzure      bool     `yaml:"http"`
	MethodsAzure   []string `yaml:"methods"`
	AuthLevelAzure string   `yaml:"authLevel"`
}

type HttpApi struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// CreateHeaderConfig sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeaderConfig(config *Configuration) {

	var region string
	var functionApp FunctionApp
	switch config.Provider {
	case "aws":
		region = amazon.AWSRegion
		s.Package.Individually = true
	case "azure":
		region = "East US 2"
		//functionApp = FunctionApp{Version: "~4"}
		s.AddPlugin("serverless-azure-functions")
		// individual packaging not available to azure
	default:
		log.Errorf("Deployment to provider %s not supported yet.", config.Provider)
	}

	s.Service = "STeLLAR" // or some other string
	s.FrameworkVersion = "3"

	s.Provider = Provider{
		Name:           config.Provider,
		Runtime:        config.Runtime,
		Region:         region,
		FunctionApp:    functionApp,
		SubscriptionId: "d3b34116-7b03-412c-997c-ca77fa672d76",
	}

	if config.Provider == "azure" {
		s.Provider.FunctionApp = functionApp
	}
}

func (s *Serverless) AddPlugin(plugin string) {
	s.Plugins = append(s.Plugins, plugin)
}

// AddPackagePattern adds a string pattern to Package.Pattern as long as such a pattern does not already exist in Package.Pattern
func (f *Function) AddPackagePattern(pattern string) {
	if !util.StringContains(f.Package.Patterns, pattern) {
		f.Package.Patterns = append(f.Package.Patterns, pattern)
	}
}

// AddFunctionConfig Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfig(subex *SubExperiment, index int, artifactPath string, provider string) {
	log.Warnf("Adding function config of Subexperiment %s, index %d", subex.Function, index)
	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}
	for i := 0; i < subex.Parallelism; i++ {

		handler := subex.Handler
		runtime := subex.Runtime
		name := createName(subex, index, i)

		f := &Function{Handler: handler, Runtime: runtime, Name: name}
		var events []Event
		switch provider {
		case "aws":
			// add indiviual packaging pattern
			f.AddPackagePattern(subex.PackagePattern)
			if artifactPath != "" {
				f.Package.Artifact = artifactPath
			}
			events = []Event{{HttpApiAWS: HttpApi{Path: "/" + name, Method: "GET"}}}
			f.Events = events
		case "azure":
			// individual packaging not available to azure
			events = []Event{{HttpAzure: true, MethodsAzure: []string{"GET"}, AuthLevelAzure: "anonymous"}}
			f.Events = events
		}
		s.Functions[name] = f
		subex.AddRoute(name)
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

	err2 := os.WriteFile(path, data, 0644)

	if err2 != nil {
		log.Fatal(err2)
	}
}

// RemoveService removes the service defined in serverless.yml
func RemoveService(path string) string {
	slsRemoveCmd := exec.Command("sls", "remove")
	slsRemoveCmd.Dir = path
	slsRemoveMessage := util.RunCommandAndLog(slsRemoveCmd)
	// cleanup
	util.RunCommandAndLog(exec.Command("rm", fmt.Sprintf("%sserverless.yml", path)))
	return slsRemoveMessage
}

// DeployService deploys the functions defined in the serverless.com file
func DeployService(path string) string {
	slsDeployCmd := exec.Command("sls", "deploy")
	slsDeployCmd.Dir = path
	slsDeployMessage := util.RunCommandAndLog(slsDeployCmd)
	return slsDeployMessage
}

// GetEndpointID scrapes the serverless deploy message for the endpoint ID
func GetEndpointID(slsDeployMessage string) string {
	lines := strings.Split(slsDeployMessage, "\n")
	if lines[1] == "endpoints:" {
		line := lines[2]
		link := strings.Split(line, " ")[4]
		httpId := strings.Split(link, ".")[0]
		endpointId := strings.Split(httpId, "//")[1]
		return endpointId
	}
	line := lines[1]
	link := strings.Split(line, " ")[3]
	httpId := strings.Split(link, ".")[0]
	endpointId := strings.Split(httpId, "//")[1]
	return endpointId
}
