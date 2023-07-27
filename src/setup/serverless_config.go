package setup

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
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
	Name   string `yaml:"name"`
	Region string `yaml:"region"`
}

type Package struct {
	Patterns []string `yaml:"patterns"`
}

type Function struct {
	Handler     string  `yaml:"handler"`
	Description string  `yaml:"description"`
	Name        string  `yaml:"name"`
	Events      []Event `yaml:"events"`
	Runtime     string  `yaml:"runtime"`
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
		Name:   config.Provider,
		Region: "us-east-1",
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

	err2 := os.WriteFile("setup/deployment/raw-code/serverless/aws/serverless.yml", data, 0644)

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

		// TODO handler and language values are to be read from SubExperiment struct and from experiment JSON config file after GitHub issue 263 is completed
		var handler string
		var language string

		switch subex.Function {
		case "hellopy":
			handler = "hellopy/lambda_function.lambda_handler"
			language = "Python"
			s.AddPackagePattern("hellopy/lambda_function.py")
		case "hellogo":
			handler = "hellogo/hellogo"
			language = "Go"
			s.AddPackagePattern("hellogo/hellogo")
			buildGoBinary("hellogo")
		default:
			log.Fatalf("DeployFunction could not recognize function image %s", subex.Function)
		}

		f := &Function{Handler: handler, Name: name, Events: events, Runtime: GetAWSLambdaRuntimeIdentifierUsingLanguageName(language)}
		s.Functions[name] = f
		subex.addRoute(name)
	}

}

// RemoveService removes the service defined in serverless.yml
func RemoveService() string {
	slsRemoveCmd := exec.Command("sls", "remove")
	slsRemoveCmd.Dir = "setup/deployment/raw-code/serverless/aws"
	slsRemoveMessage := util.RunCommandAndLog(slsRemoveCmd)
	return slsRemoveMessage
}

// Deploys the functions defined in the serverless.com file
func deployService() string {
	slsDeployCmd := exec.Command("sls", "deploy")
	slsDeployCmd.Dir = "setup/deployment/raw-code/serverless/aws"
	slsDeployMessage := util.RunCommandAndLog(slsDeployCmd)
	return slsDeployMessage
}

func GetAWSLambdaRuntimeIdentifierUsingLanguageName(language string) string {
	supportedRuntimes := map[string]string{
		"Python": "python3.10",
		"Java":   "java17",
		"Go":     "go1.x",
	}
	runtime, ok := supportedRuntimes[language]
	if !ok {
		log.Fatalf("Unable to get AWS Lambda runtime identifier for the langauge %s", language)
	}
	return runtime
}

func buildGoBinary(functionImageName string) {
	fullPath := fmt.Sprintf("setup/deployment/raw-code/serverless/aws/%s", functionImageName)
	command := exec.Command("env", "GOOS=linux", "GOARCH=amd64", "go", "build", "-o", functionImageName, "-C", fullPath)
	util.RunCommandAndLog(command)
}
