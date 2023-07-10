package setup

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"strconv"
)

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

func (s *Serverless) CreateHeader(config Configuration) {
	s.Service = "STeLLAR" // or some other string
	s.FrameworkVersion = "3"
	s.Provider = Provider{
		Name:    config.Provider,
		Runtime: "python3.9",
		Region:  "us-east-1",
	}
	s.Functions = map[string]*Function{}
	log.Info(s.Service)
}

func (s *Serverless) AddPackagePattern(pattern string) {
	s.Package.Patterns = append(s.Package.Patterns, pattern)
}

func (s *Serverless) CreateServerlessConfigFile() {
	data, err := yaml.Marshal(&s)
	log.Info(s.Service)
	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile("src/setup/deployment/raw-code/serverless/aws/serverless.yml", data, 0644)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (s *Serverless) AddFunctionConfig(subex SubExperiment, index int) {
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	name := nonAlphanumericRegex.ReplaceAllString(subex.Title, "") + "_" + strconv.Itoa(index)
	events := []Event{Event{HttpApi{Path: "/" + name, Method: "GET"}}}
	f := &Function{Handler: subex.Function, Name: name, Events: events}
	s.Functions[name] = f
	log.Info(s.Functions)
}

func createServerlessYml() {

	provider := Provider{
		Name:    "hellopy",
		Runtime: "python3.9",
		Region:  "us-east-1",
	}

	packageStruct := Package{Patterns: []string{"!**", "hellopy/lambda_function.py"}}

	//function1 := Function{
	//	Handler:     "aws/hellopy/lambda_function.lambda_handler",
	//	Name:        "test1",
	//	Description: "Testing serverless.com deployment for stellar.",
	//}

	//functions := map[string]Function{"test1": function1}
	serverless := Serverless{
		Service:          "Stellar",
		FrameworkVersion: "3",
		Provider:         provider,
		Package:          packageStruct,
		//Functions:        functions,
	}

	data, err := yaml.Marshal(serverless)

	if err != nil {
		log.Fatal(err)
	}

	err2 := ioutil.WriteFile("aws/serverless.yml", data, 0644)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func main() {
	createServerlessYml()
}
