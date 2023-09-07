package setup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"regexp"
	"stellar/util"
	"strings"
)

// Serverless describes the serverless.yml contents.
type Serverless struct {
	Service          string               `yaml:"service"`
	FrameworkVersion string               `yaml:"frameworkVersion"`
	Provider         Provider             `yaml:"provider"`
	Plugins          []string             `yaml:"plugins,omitempty"`
	Package          Package              `yaml:"package"`
	Functions        map[string]*Function `yaml:"functions"`
}

type Provider struct {
	Name    string `yaml:"name"`
	Runtime string `yaml:"runtime"`
	Region  string `yaml:"region"`
}

type Package struct {
	Individually bool `yaml:"individually"`
}

type Function struct {
	Handler   string          `yaml:"handler"`
	Runtime   string          `yaml:"runtime"`
	Name      string          `yaml:"name"`
	Events    []Event         `yaml:"events"`
	Package   FunctionPackage `yaml:"package"`
	SnapStart bool            `yaml:"snapStart,omitempty"`
}

type FunctionPackage struct {
	Patterns []string `yaml:"patterns"`
	Artifact string   `yaml:"artifact,omitempty"`
}

type Event struct {
	AWSHttpEvent   AWSHttpEvent `yaml:"httpApi,omitempty"`
	AzureHttp      bool         `yaml:"http,omitempty"`
	AzureMethods   []string     `yaml:"methods,omitempty"`
	AzureAuthLevel string       `yaml:"authLevel,omitempty"`
}

type AWSHttpEvent struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

const (
	AWS_DEFAULT_REGION   = "us-west-1"
	AZURE_DEFAULT_REGION = "Southeast Asia"
)

// CreateHeaderConfig sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeaderConfig(config *Configuration, serviceName string) {

	var region string
	switch config.Provider {
	case "aws":
		region = AWS_DEFAULT_REGION
	case "azure":
		region = AZURE_DEFAULT_REGION
	default:
		log.Errorf("Deployment to provider %s not supported yet.", config.Provider)
	}

	s.Service = serviceName
	s.FrameworkVersion = "3"

	s.Provider = Provider{
		Name:    config.Provider,
		Runtime: config.Runtime,
		Region:  region,
	}
}

func (s *Serverless) addPlugin(pluginName string) {
	s.Plugins = append(s.Plugins, pluginName)
}

// packageIndividually enables individual packaging for providers like AWS, it is not supported by Azure
func (s *Serverless) packageIndividually() {
	s.Package.Individually = true
}

// AddPackagePattern adds a string pattern to Package.Pattern as long as such a pattern does not already exist in Package.Pattern
func (f *Function) AddPackagePattern(pattern string) {
	if !util.StringContains(f.Package.Patterns, pattern) {
		f.Package.Patterns = append(f.Package.Patterns, pattern)
	}
}

// AddFunctionConfigAWS Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfigAWS(subex *SubExperiment, index int, artifactPath string) {
	log.Infof("Adding function config of Subexperiment %s, index %d", subex.Function, index)
	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}
	for i := 0; i < subex.Parallelism; i++ {

		handler := subex.Handler
		runtime := subex.Runtime
		name := createName(subex, index, i)
		events := []Event{{AWSHttpEvent: AWSHttpEvent{Path: "/" + name, Method: "GET"}}}

		f := &Function{Handler: handler, Runtime: runtime, Name: name, Events: events}
		f.AddPackagePattern(subex.PackagePattern)
		if artifactPath != "" {
			f.Package.Artifact = artifactPath
		}
		if subex.SnapStartEnabled { // Add SnapStart field only if it is enabled
			f.SnapStart = true
		}
		s.Functions[name] = f
		subex.AddRoute(name)
		// TODO: producer-consumer sub-function definition
	}
}

// AddFunctionConfigAzure Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfigAzure(subex *SubExperiment, index int, artifactPath string) {
	log.Infof("Adding function config of Subexperiment %s, index %d", subex.Function, index)

	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}

	for i := 0; i < subex.Parallelism; i++ {
		handler := subex.Handler
		runtime := subex.Runtime
		name := createName(subex, index, i)
		events := []Event{
			{
				AzureHttp:      true,
				AzureMethods:   []string{"GET"},
				AzureAuthLevel: "anonymous",
			},
		}

		f := &Function{Handler: handler, Runtime: runtime, Name: name, Events: events}
		f.AddPackagePattern(subex.PackagePattern)
		s.Functions[name] = f
		subex.AddRoute(name)
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

// RemoveServiceAWS removes the service defined in serverless.yml
func RemoveServiceAWS(path string) string {
	slsRemoveCmd := exec.Command("sls", "remove")
	slsRemoveCmd.Dir = path
	slsRemoveMessage := util.RunCommandAndLog(slsRemoveCmd)
	// cleanup
	util.RunCommandAndLog(exec.Command("rm", fmt.Sprintf("%sserverless.yml", path)))
	return slsRemoveMessage
}

// RemoveServiceAzure removes the service defined in serverless.yml
func RemoveServiceAzure(path string) string {
	slsRemoveCmd := exec.Command("sls", "remove", "--force")
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

// GetEndpointIDFromAWSDeployment scrapes the serverless deploy message for the endpoint ID
func GetEndpointIDFromAWSDeployment(slsDeployMessage string) string {
	regex := regexp.MustCompile(`https:\/\/(.*)\.execute`)
	return regex.FindStringSubmatch(slsDeployMessage)[1]
}

func GetEndpointIDFromAzureDeployment(message string) string {
	methodAndEndpointRegex := regexp.MustCompile(`\[GET] .+\n`)
	methodAndEndpoint := methodAndEndpointRegex.FindString(message) // e.g. [GET] sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_0
	endpoint := strings.Split(methodAndEndpoint, " ")[1]            // e.g. sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_0
	endpointId := strings.Split(endpoint, ".")[0]                   // e.g. sls-seasi-dev-stellar-sub-experiment-1
	return endpointId
}
