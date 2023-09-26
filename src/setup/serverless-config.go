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
	"time"
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
	Name        string `yaml:"name"`
	Runtime     string `yaml:"runtime"`
	Region      string `yaml:"region"`
	Credentials string `yaml:"credentials,omitempty"`
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
	AWSEvent     *AWSEvent
	AzureEvent   *AzureEvent
	AlibabaEvent *AlibabaEvent
}

func (event Event) MarshalYAML() (interface{}, error) {
	if event.AWSEvent != nil {
		return event.AWSEvent, nil
	} else if event.AzureEvent != nil {
		return event.AzureEvent, nil
	} else if event.AlibabaEvent != nil {
		return event.AlibabaEvent, nil
	} else {
		return nil, nil
	}
}

type AWSEvent struct {
	AWSHttpEvent AWSHttpEvent `yaml:"httpApi"`
}

type AWSHttpEvent struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type AzureEvent struct {
	AzureHttpEvent AzureHttpEvent `yaml:",omitempty,inline"`
}

type AzureHttpEvent struct {
	AzureHttp      bool     `yaml:"http"`
	AzureMethods   []string `yaml:"methods"`
	AzureAuthLevel string   `yaml:"authLevel"`
}

type AlibabaEvent struct {
	AlibabaHttpEvent AlibabaHttpEvent `yaml:"http"`
}

type AlibabaHttpEvent struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

const (
	AWS_DEFAULT_REGION         = "us-west-1"
	AZURE_DEFAULT_REGION       = "West US 2"
	GCR_DEFAULT_REGION         = "us-west1"
	ALIBABA_DEFAULT_REGION     = "us-west-1"
	ALIBABA_DEFAULT_ACCOUNT_ID = "5776795023355240"
)

// CreateHeaderConfig sets the fields Service, FrameworkVersion, and Provider
func (s *Serverless) CreateHeaderConfig(config *Configuration, serviceName string) {

	var region string
	switch config.Provider {
	case "aws":
		region = AWS_DEFAULT_REGION
	case "gcr":
		region = GCR_DEFAULT_REGION
	case "azure":
		region = AZURE_DEFAULT_REGION
	case "aliyun":
		region = ALIBABA_DEFAULT_REGION
	default:
		log.Errorf("Deployment to provider %s not supported yet.", config.Provider)
	}

	s.Service = serviceName
	s.FrameworkVersion = "3"

	if config.Provider == "aliyun" {
		s.Provider = Provider{
			Name:        config.Provider,
			Runtime:     config.Runtime,
			Region:      region,
			Credentials: "~/.aliyuncli/credentials",
		}
	} else {
		s.Provider = Provider{
			Name:    config.Provider,
			Runtime: config.Runtime,
			Region:  region,
		}
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
		events := []Event{
			{
				AWSEvent: &AWSEvent{
					AWSHttpEvent: AWSHttpEvent{
						Path:   "/" + name,
						Method: "GET",
					},
				},
			},
		}

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
func (s *Serverless) AddFunctionConfigAzure(subex *SubExperiment, index int, parallelism int) {
	log.Infof("Adding function config of Subexperiment %s, index %d", subex.Function, index)

	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}

	handler := subex.Handler
	runtime := subex.Runtime
	name := createName(subex, index, parallelism)
	events := []Event{
		{
			AzureEvent: &AzureEvent{
				AzureHttpEvent: AzureHttpEvent{
					AzureHttp:      true,
					AzureMethods:   []string{"GET"},
					AzureAuthLevel: "anonymous",
				},
			},
		},
	}

	function := &Function{Handler: handler, Runtime: runtime, Name: name, Events: events}
	function.AddPackagePattern(subex.PackagePattern)
	s.Functions[name] = function
	subex.AddRoute(name)
}

// AddFunctionConfigAlibaba Adds a function to the service. If parallelism = n, then it defines n functions. Also deploys all producer-consumer subfunctions.
func (s *Serverless) AddFunctionConfigAlibaba(subex *SubExperiment, index int, artifactPath string) {
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
				AlibabaEvent: &AlibabaEvent{
					AlibabaHttpEvent: AlibabaHttpEvent{
						Path:   "/" + name,
						Method: "GET",
					},
				},
			},
		}

		function := &Function{Handler: handler, Runtime: runtime, Name: name, Events: events}
		function.AddPackagePattern(subex.PackagePattern)
		s.Functions[name] = function
		subex.AddRoute(name)
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

// RemoveService removes the services created by experiments
func RemoveService(config *Configuration, path string) string {
	switch config.Provider {
	case "aws":
		return RemoveServerlessService(path)
	case "azure":
		RemoveAzureAllServices(config.SubExperiments, path)
		return "All Azure services removed."
	case "gcr":
		RemoveGCRAllServices(config.SubExperiments)
		return "All GCR services deleted."
	case "cloudflare":
		RemoveCloudflareAllWorkers(config.SubExperiments)
		return "All Cloudflare Workers deleted."
	case "aliyun":
		RemoveAlibabaAllServices(path, len(config.SubExperiments))
		return "All Alibaba Cloud services removed."
	default:
		log.Fatalf(fmt.Sprintf("Failed to remove service for unrecognised provider %s", config.Provider))
		return ""
	}
}

// RemoveServerlessService removes a service that was deployed using the Serverless framework
func RemoveServerlessService(path string) string {
	log.Infof(fmt.Sprintf("Removing Serverless service at %s", path))
	slsRemoveCmd := exec.Command("sls", "remove")
	slsRemoveCmd.Dir = path
	slsRemoveCmdOutput := util.RunCommandAndLog(slsRemoveCmd)

	util.RunCommandAndLog(exec.Command("rm", fmt.Sprintf("%sserverless.yml", path)))

	return slsRemoveCmdOutput
}

// RemoveServerlessServiceForcefully forcefully removes a service that was deployed using the Serverless framework
func RemoveServerlessServiceForcefully(path string) string {
	log.Infof(fmt.Sprintf("Removing Serverless service at %s", path))
	slsRemoveCmd := exec.Command("sls", "remove", "--force")
	slsRemoveCmd.Dir = path
	slsRemoveCmdOutput := util.RunCommandAndLog(slsRemoveCmd)

	deleteSlsConfigFileCmd := exec.Command("rm", "serverless.yml")
	deleteSlsConfigFileCmd.Dir = path
	util.RunCommandAndLog(deleteSlsConfigFileCmd)

	return slsRemoveCmdOutput
}

// RemoveAzureAllServices removes all Azure services
func RemoveAzureAllServices(subExperiments []SubExperiment, path string) []string {
	var removeServiceMessages []string
	for subExperimentIndex, subExperiment := range subExperiments {
		for parallelism := 0; parallelism < subExperiment.Parallelism; parallelism++ {
			deploymentDir := fmt.Sprintf("%ssub-experiment-%d/parallelism-%d", path, subExperimentIndex, parallelism)
			slsRemoveCmdOutput := RemoveServerlessServiceForcefully(deploymentDir)
			removeServiceMessages = append(removeServiceMessages, slsRemoveCmdOutput)
		}
	}
	return removeServiceMessages
}

// RemoveGCRAllServices removes all GCR services defined in the Subexperiment array
func RemoveGCRAllServices(subExperiments []SubExperiment) []string {
	var deleteServiceMessages []string
	for index, subex := range subExperiments {
		for i := 0; i < subex.Parallelism; i++ {
			service := createName(&subex, index, i)
			deleteMsg := RemoveGCRSingleService(service)
			deleteServiceMessages = append(deleteServiceMessages, deleteMsg)
		}
	}
	return deleteServiceMessages
}

// RemoveGCRSingleService removes a single GCR service
func RemoveGCRSingleService(service string) string {
	log.Infof("Deleting GCR service %s...", service)
	deleteServiceCommand := exec.Command("gcloud", "run", "services", "delete", "--quiet", "--region", GCR_DEFAULT_REGION, service)
	deleteMessage := util.RunCommandAndLog(deleteServiceCommand)
	return deleteMessage
}

// RemoveCloudflareAllWorkers removes all Cloudflare Workers
func RemoveCloudflareAllWorkers(subExperiments []SubExperiment) []string {
	log.Infof("Removing Cloudflare Workers...")
	var removeServiceMessages []string
	for index, subex := range subExperiments {
		for i := 0; i < subex.Parallelism; i++ {
			workerName := createName(&subex, index, i)
			removeMessage := RemoveCloudflareSingleWorker(workerName)
			removeServiceMessages = append(removeServiceMessages, removeMessage)
		}
	}
	return removeServiceMessages
}

// RemoveCloudflareSingleWorker removes a single Cloudflare Worker specified by name
func RemoveCloudflareSingleWorker(workerName string) string {
	log.Infof("Removing Cloudflare Worker %s...", workerName)
	removeWorkerCommand := exec.Command("wrangler", "delete", "--name", workerName, "--force")
	removeMessage := util.RunCommandAndLog(removeWorkerCommand)
	return removeMessage
}

// RemoveAlibabaAllServices removes all Alibaba Cloud services
func RemoveAlibabaAllServices(path string, numSubExperiments int) []string {
	alibabaCloudAccountId := os.Getenv("ALIYUN_ACCOUNT_ID")
	if alibabaCloudAccountId == "" {
		alibabaCloudAccountId = ALIBABA_DEFAULT_ACCOUNT_ID
	}
	nameOfBucketToDelete := fmt.Sprintf("oss://sls-%s-%s", alibabaCloudAccountId, ALIBABA_DEFAULT_REGION)
	util.RunCommandAndLog(exec.Command("aliyun", "oss", "rm", "--bucket", "--recursive", "--force", nameOfBucketToDelete))

	var removeServiceMessages []string
	for i := 0; i < numSubExperiments; i++ {
		subExPath := fmt.Sprintf("%ssub-experiment-%d/", path, i)
		slsRemoveCmdOutput := RemoveServerlessService(subExPath)
		removeServiceMessages = append(removeServiceMessages, slsRemoveCmdOutput)
	}
	return removeServiceMessages
}

// DeployService deploys the functions defined in the serverless.com file
func DeployService(path string) string {
	log.Infof("AWS_ACCESS_KEY_ID: %s", os.Getenv("AWS_ACCESS_KEY_ID"))
	log.Infof(fmt.Sprintf("Deploying service at %s", path))
	slsDeployCmd := exec.Command("bash", "-c", "sls deploy")
	slsDeployCmd.Dir = path
	slsDeployMessage := util.RunCommandAndLog(slsDeployCmd)
	return slsDeployMessage
}

// DeployGCRContainerService deploys a container service to cloud provider
func (s *Serverless) DeployGCRContainerService(subex *SubExperiment, index int, imageLink string, path string, region string) {
	log.Infof("Deploying container service(s) to GCR...")
	for i := 0; i < subex.Parallelism; i++ {
		name := createName(subex, index, i)

		gcrDeployCommand := exec.Command("gcloud", "run", "deploy", name, "--image", imageLink, "--allow-unauthenticated", "--region", region)
		deployMessage := util.RunCommandAndLog(gcrDeployCommand)
		subex.Endpoints = append(subex.Endpoints, EndpointInfo{ID: GetGCREndpointID(deployMessage)})
		subex.AddRoute("")
	}
}

func DeployCloudflareWorkers(subex *SubExperiment, index int, path string) {
	log.Infof("Deploying Cloudflare Workers...")
	for i := 0; i < subex.Parallelism; i++ {
		name := createName(subex, index, i)

		cloudFlareDeployCommand := exec.Command("wrangler", "deploy", fmt.Sprintf("%s/%s/%s", path, subex.Function, subex.Handler), "--name", name, "--compatibility-date", time.Now().Format("2006-01-02"))
		deployMessage := util.RunCommandAndLog(cloudFlareDeployCommand)
		subex.Endpoints = append(subex.Endpoints, EndpointInfo{ID: GetCloudflareEndpointID(deployMessage)})
		subex.AddRoute("")
	}
}

// GetAWSEndpointID scrapes the serverless deploy message for the endpoint ID
func GetAWSEndpointID(slsDeployMessage string) string {
	regex := regexp.MustCompile(`https://(.*)\.execute`)
	return regex.FindStringSubmatch(slsDeployMessage)[1]
}

// GetGCREndpointID scrapes the gcloud run deploy message for the endpoint ID
func GetGCREndpointID(deployMessage string) string {
	regex := regexp.MustCompile(`https://.*\.run\.app`)
	endpointID := strings.Split(regex.FindString(deployMessage), "//")[1]
	return endpointID
}

// GetAzureEndpointID finds the Azure endpoint ID from the deployment message
func GetAzureEndpointID(message string) string {
	methodAndEndpointRegex := regexp.MustCompile(`\[GET] .+\n`)
	methodAndEndpoint := methodAndEndpointRegex.FindString(message) // e.g. [GET] sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_0
	endpoint := strings.Split(methodAndEndpoint, " ")[1]            // e.g. sls-seasi-dev-stellar-sub-experiment-1.azurewebsites.net/api/subexperiment2_1_0
	endpointId := strings.Split(endpoint, ".")[0]                   // e.g. sls-seasi-dev-stellar-sub-experiment-1
	return endpointId
}

// GetAlibabaEndpointID finds the Alibaba Cloud endpoint ID from the deployment message
func GetAlibabaEndpointID(message string) string {
	// Example Alibaba endpoint
	// GET http://5cfeb440ed6d4ad69ae29d8408aa606e-ap-southeast-1.alicloudapi.com/foo -> my-service-dev.my-service-dev-hello
	re := regexp.MustCompile(`GET http://(?P<endpointId>[A-Za-z0-9]+)[-a-z0-9]+.alicloudapi.com`)
	matches := re.FindStringSubmatch(message)
	endpointIdSubexpIndex := re.SubexpIndex("endpointId")
	return matches[endpointIdSubexpIndex]
}

// GetCloudflareEndpointID finds the Cloudflare endpoint ID from the deployment message
func GetCloudflareEndpointID(message string) string {
	regex := regexp.MustCompile(`https://.*\.workers\.dev`)
	endpointID := strings.Split(regex.FindString(message), "//")[1]
	return endpointID
}
