package setup

import (
	"fmt"
	//"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"stellar/setup/deployment"
	"stellar/util"
	"strings"

	//"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
	Name           string      `yaml:"name"`
	Runtime        string      `yaml:"runtime"`
	Region         string      `yaml:"region"`
	Credentials    string      `yaml:"credentials,omitempty"`
	SubscriptionId string      `yaml:"subscriptionId,omitempty"`
	TenantId       string      `yaml:"tenantId,omitempty"`
	AppId          string      `yaml:"appId,omitempty"`
	Password       string      `yaml:"password,omitempty"`
	FunctionApp    FunctionApp `yaml:"functionApp,omitempty"`
}

type FunctionApp struct {
	ExtensionVersion string `yaml:"extensionVersion"`
}

type Package struct {
	Individually bool `yaml:"individually"`
}

type Function struct {
	Handler   string          `yaml:"handler"`
	Runtime   string          `yaml:"runtime,omitempty"`
	Name      string          `yaml:"name,omitempty"`
	Events    []Event         `yaml:"events"`
	Package   FunctionPackage `yaml:"package,omitempty"`
	SnapStart bool            `yaml:"snapStart,omitempty"`
}

type FunctionPackage struct {
	Patterns []string `yaml:"patterns"`
	Artifact string   `yaml:"artifact,omitempty"`
}

type Event struct {
	AWSEvent     *AWSEvent
	Http         *AzureHttpEvent `yaml:"http,omitempty"`
	AlibabaEvent *AlibabaEvent
}

func (event Event) MarshalYAML() (interface{}, error) {
	if event.Http != nil {
		// For Azure HTTP events, wrap the event in a map with the key "http"
		return map[string]interface{}{
			"http": event.Http,
		}, nil
	} else if event.AWSEvent != nil {
		// For AWS events, return the AWSEvent as is
		return event.AWSEvent, nil
	} else if event.AlibabaEvent != nil {
		// For Alibaba events, return the AlibabaEvent as is
		return event.AlibabaEvent, nil
	} else {
		// If no events are set, return nil
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
	Method    string `yaml:"method"`
	Route     string `yaml:"route"`
	AuthLevel string `yaml:"authLevel,omitempty"`
}

type AlibabaEvent struct {
	AlibabaHttpEvent AlibabaHttpEvent `yaml:"http"`
}

type AlibabaHttpEvent struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

var nonAlphanumericRegex *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
var providerFunctionNames map[string]([]string) = make(map[string]([]string))

const (
	AWS_DEFAULT_REGION         = "us-west-1"
	AZURE_DEFAULT_REGION       = "West US"
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

	runtimeValue := config.Runtime
	if config.Runtime == "go1.x" {
		log.Warnf("`go1.x` runtime is deprecated. Runtime `provided.al2023` will be used instead...")
		runtimeValue = "provided.al2023"
	}

	switch config.Provider {
	case "azure":
		s.Provider = Provider{
			Name:           config.Provider,
			Runtime:        runtimeValue,
			Region:         region,
			SubscriptionId: "${env:b1dcfb0e-cca5-40f4-af7c-f5e0d8e55681}",
			TenantId:       "${env:15ce9348-be2a-462b-8fc0-e1765a9b204a}",
			AppId:          "${env:afc74a64-222c-4560-a633-ca77a7a00eb5}",
			//Password:       "${env:P.D8Q~ykEd3VJfRVWcCvYTl~w0UrtpNRd1xX4bA2}",
			FunctionApp: FunctionApp{ExtensionVersion: "~4"},
		}
	case "aliyun":
		s.Provider = Provider{
			Name:        config.Provider,
			Runtime:     runtimeValue,
			Region:      region,
			Credentials: "~/.aliyuncli/credentials",
		}
	default:
		s.Provider = Provider{
			Name:    config.Provider,
			Runtime: runtimeValue,
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
func (s *Serverless) AddFunctionConfigAWS(subex *SubExperiment, index int, randomTag string, artifactPath string) {
	log.Infof("Adding function config of Subexperiment %s, index %d", subex.Function, index)
	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}

	for i := 0; i < subex.Parallelism; i++ {
		handler := subex.Handler
		runtime := subex.Runtime
		if runtime == "go1.x" {
			log.Warnf("`go1.x` runtime is deprecated. Using `provided.al2023` runtime instead...")
			runtime = "provided.al2023"
		}
		name := fmt.Sprintf("%s-%s", randomTag, createName(subex, index, i))
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

// AddFunctionConfigAzure adds a function to the service for Azure provider.
func (s *Serverless) AddFunctionConfigAzure(subex *SubExperiment, index int, name string) {
	log.Infof("Adding function config for Subexperiment %s, index %d", subex.Function, index)

	if s.Functions == nil {
		s.Functions = make(map[string]*Function)
	}

	handler := subex.Handler

	// Define the HTTP event
	events := []Event{
		{
			Http: &AzureHttpEvent{
				Method:    "GET",       // Or the appropriate HTTP method
				Route:     name,        // The route for the function
				AuthLevel: "anonymous", // Authentication level
			},
		},
	}

	// Create the function configuration
	function := &Function{
		Handler: handler,
		Events:  events,
	}

	// Add package patterns if needed
	function.AddPackagePattern(subex.PackagePattern)

	// Add the function to the service
	s.Functions[name] = function
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
	slsRemoveCmdOutput := util.RunCommandAndLogWithRetries(slsRemoveCmd, 3)

	// Optionally, remove the serverless.yml file if necessary
	util.RunCommandAndLog(exec.Command("rm", "-f", filepath.Join(path, "serverless.yml")))

	return slsRemoveCmdOutput
}

// RemoveServerlessServiceForcefully forcefully removes a service that was deployed using the Serverless framework
func RemoveServerlessServiceForcefully(path string) string {
	log.Infof(fmt.Sprintf("Removing Serverless service at %s", path))
	slsRemoveCmd := exec.Command("sls", "remove", "--force")
	slsRemoveCmd.Dir = path
	slsRemoveCmdOutput := util.RunCommandAndLogWithRetries(slsRemoveCmd, 3)

	deleteSlsConfigFileCmd := exec.Command("rm", "serverless.yml")
	deleteSlsConfigFileCmd.Dir = path
	util.RunCommandAndLog(deleteSlsConfigFileCmd)

	return slsRemoveCmdOutput
}

// RemoveAzureAllServices removes all Azure services
func RemoveAzureAllServices(subExperiments []SubExperiment, path string) []string {
	var removeServiceMessages []string
	for subExperimentIndex := range subExperiments {
		deploymentDir := filepath.Join(path, fmt.Sprintf("sub-experiment-%d", subExperimentIndex))
		removeMessage := RemoveServerlessService(deploymentDir)
		removeServiceMessages = append(removeServiceMessages, removeMessage)
	}
	return removeServiceMessages
}

/* //remove for now
func removeSubExperimentParallelismInBatches(path string, subExperimentIndex int, subExperiment SubExperiment, removeServiceMessages []string, functionsPerBatch int) {
	numberOfBatches := int(math.Ceil(float64(subExperiment.Parallelism) / float64(functionsPerBatch)))

	for batchNumber := 0; batchNumber < numberOfBatches; batchNumber++ {
		mu := sync.Mutex{}
		wg := sync.WaitGroup{}

		for parallelism := batchNumber * functionsPerBatch; parallelism < (batchNumber+1)*functionsPerBatch && parallelism < subExperiment.Parallelism; parallelism++ {
			wg.Add(1)

			go func(parallelism int) {
				defer wg.Done()

				deploymentDir := fmt.Sprintf("%ssub-experiment-%d/parallelism-%d", path, subExperimentIndex, parallelism)
				slsRemoveCmdOutput := RemoveServerlessServiceForcefully(deploymentDir)
				mu.Lock()
				defer mu.Unlock()
				removeServiceMessages = append(removeServiceMessages, slsRemoveCmdOutput)
			}(parallelism)
		}
		wg.Wait()
	}
}
*/
// RemoveGCRAllServices removes all GCR services defined in the Subexperiment array
func RemoveGCRAllServices(subExperiments []SubExperiment) []string {
	var deleteServiceMessages []string
	for _, functionName := range providerFunctionNames["gcr"] {
		deleteMsg := RemoveGCRSingleService(functionName)
		deleteServiceMessages = append(deleteServiceMessages, deleteMsg)
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
	for _, functionName := range providerFunctionNames["cloudflare"] {
		removeMessage := RemoveCloudflareSingleWorker(functionName)
		removeServiceMessages = append(removeServiceMessages, removeMessage)
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
	log.Infof(fmt.Sprintf("Deploying service at %s", path))
	slsDeployCmd := exec.Command("bash", "-c", "sls deploy")
	slsDeployCmd.Dir = path
	slsDeployMessage := util.RunCommandAndLogWithRetries(slsDeployCmd, 3)
	return slsDeployMessage
}

// DeployGCRContainerService deploys a container service to cloud provider
func (s *Serverless) DeployGCRContainerService(subex *SubExperiment, index int, randomTag string, imageLink string, path string, region string) {
	log.Infof("Deploying container service(s) to GCR...")
	for i := 0; i < subex.Parallelism; i++ {
		name := fmt.Sprintf("%s-%s", randomTag, createName(subex, index, i))
		providerFunctionNames["gcr"] = append(providerFunctionNames["gcr"], name) // Used for function removal

		gcrDeployCommand := exec.Command("gcloud", "run", "deploy", name, "--image", imageLink, "--allow-unauthenticated", "--region", region)
		deployMessage := util.RunCommandAndLog(gcrDeployCommand)
		subex.Endpoints = append(subex.Endpoints, EndpointInfo{ID: GetGCREndpointID(deployMessage)})
		subex.AddRoute("")
	}
}

func DeployCloudflareWorkers(subex *SubExperiment, index int, randomTag string, path string) {
	log.Infof("Deploying Cloudflare Workers...")
	for i := 0; i < subex.Parallelism; i++ {
		name := fmt.Sprintf("%s-%s", randomTag, createName(subex, index, i))
		providerFunctionNames["cloudflare"] = append(providerFunctionNames["cloudflare"], name) // Used for function removal

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

// GetAzureEndpointID extracts the Azure function endpoint URL from the deployment output
func GetAzureEndpointID(message string) string {
	// Use a regex to find the line containing the endpoint URL
	regex := regexp.MustCompile(`(?m)Invoke URL:\s+(https://.+)`)
	matches := regex.FindStringSubmatch(message)
	if len(matches) > 1 {
		endpointURL := matches[1]
		return endpointURL
	} else {
		log.Fatalf("Could not find the endpoint URL in deployment output.")
		return ""
	}
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

// Container function that calls RunDeployment
func RunAzureDeployment() {
	deployment.RunDeployment()
}
