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
	HttpApi HttpApi `yaml:"httpApi"`
}

type HttpApi struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

const (
	AWS_DEFAULT_REGION = "us-west-1"
	GCR_DEFAULT_REGION = "us-west1"
)

// CreateHeaderConfig sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeaderConfig(config *Configuration) {

	var region string
	switch config.Provider {
	case "aws":
		region = AWS_DEFAULT_REGION
	case "gcr":
		region = GCR_DEFAULT_REGION
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
	s.Package.Individually = true
}

// AddPackagePattern adds a string pattern to Package.Pattern as long as such a pattern does not already exist in Package.Pattern
func (f *Function) AddPackagePattern(pattern string) {
	if !util.StringContains(f.Package.Patterns, pattern) {
		f.Package.Patterns = append(f.Package.Patterns, pattern)
	}
}

// AddFunctionConfig Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfig(subex *SubExperiment, index int, artifactPath string) {
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

// createName removes no-alphanumeric characters as serverless.com functions requires alphanumeric names. It also adds alphanumeric indexes to ensure a unique function name.
func createName(subex *SubExperiment, index int, parallelism int) string {
	return fmt.Sprintf("%s-%d-%d", nonAlphanumericRegex.ReplaceAllString(subex.Title, ""), index, parallelism)
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
func RemoveService(provider string, path string) string {
	switch provider {
	case "gcr":
		cmd := fmt.Sprintf("gcloud run services list --region %s | awk '{print $2}' | awk NR\\>1", GCR_DEFAULT_REGION)
		getServicesCommand := exec.Command("bash", "-c", cmd)
		lines := util.RunCommandAndLog(getServicesCommand)
		services := strings.Split(lines, "\n")
		services = services[:len(services)-1] // Remove last empty line
		for _, service := range services {
			log.Infof("Deleting GCR service %s...", service)
			deleteServiceCommand := exec.Command("gcloud", "run", "services", "delete", "--quiet", "--region", GCR_DEFAULT_REGION, service)
			util.RunCommandAndLog(deleteServiceCommand)
			log.Infof("Deleted GCR service %s", service)
		}
		return ""
	default:
		slsRemoveCmd := exec.Command("sls", "remove")
		slsRemoveCmd.Dir = path
		slsRemoveMessage := util.RunCommandAndLog(slsRemoveCmd)
		// cleanup
		util.RunCommandAndLog(exec.Command("rm", fmt.Sprintf("%sserverless.yml", path)))
		log.Info(slsRemoveMessage)
		return slsRemoveMessage
	}
}

// DeployService deploys the functions defined in the serverless.com file
func DeployService(path string) string {
	slsDeployCmd := exec.Command("sls", "deploy")
	slsDeployCmd.Dir = path
	slsDeployMessage := util.RunCommandAndLog(slsDeployCmd)
	return slsDeployMessage
}

func (s *Serverless) DeployContainerService(subex *SubExperiment, index int, imageLink string, path string, region string) {
	switch s.Provider.Name {
	case "gcr":
		log.Infof("Deploying container service(s) to GCR...")
		for i := 0; i < subex.Parallelism; i++ {
			name := createName(subex, index, i)

			gcrDeployCommand := exec.Command("gcloud", "run", "deploy", name, "--image", imageLink, "--allow-unauthenticated", "--region", region)
			deployMessage := util.RunCommandAndLog(gcrDeployCommand)
			log.Info(deployMessage)
			subex.Endpoints = append(subex.Endpoints, EndpointInfo{ID: GetGCREndpointID(deployMessage)})
			subex.AddRoute("")
		}
	default:
		log.Fatalf("Container deployment not supported for provider %s", s.Provider.Name)
	}
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

func GetGCREndpointID(deployMessage string) string {
	regex := regexp.MustCompile(`https:\/\/.*\.run\.app`)
	endpointID := strings.Split(regex.FindString(deployMessage), "//")[1]
	return endpointID

}
