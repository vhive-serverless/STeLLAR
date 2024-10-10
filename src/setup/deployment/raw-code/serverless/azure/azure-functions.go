// src/setup/deployment/raw-code/serverless/azure/azure-functions.go

package azure

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

// Config holds all the configuration variables loaded from the .env file
type Config struct {
	AzureSubscriptionID     string
	AzureLocation           string
	AzureResourceGroupName  string
	AzureStorageAccountName string
	AzureFunctionAppName    string
	FunctionName            string
	FunctionTemplate        string
	AuthLevel               string
	KeepResource            string
}

// Global variables for Azure SDK clients
var (
	resourcesClientFactory *armresources.ClientFactory
	storageClientFactory   *armstorage.ClientFactory
	resourceGroupClient    *armresources.ResourceGroupsClient
	accountsClient         *armstorage.AccountsClient
)

// functionProjectDir defines the directory for your Function App project
const functionProjectDir = `C:\Projects\STeLLAR\functionapp` // Ensure this path exists

// loadConfig retrieves environment variables and populates the Config struct
func LoadConfig() Config {
	return Config{
		AzureSubscriptionID:     os.Getenv("AZURE_SUBSCRIPTION_ID"),
		AzureLocation:           os.Getenv("AZURE_LOCATION"),
		AzureResourceGroupName:  os.Getenv("AZURE_RESOURCE_GROUP_NAME"),
		AzureStorageAccountName: os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AzureFunctionAppName:    os.Getenv("AZURE_FUNCTION_APP_NAME"),
		FunctionName:            os.Getenv("FUNCTION_NAME"),
		FunctionTemplate:        os.Getenv("FUNCTION_TEMPLATE"),
		AuthLevel:               os.Getenv("AUTH_LEVEL"),
		KeepResource:            os.Getenv("KEEP_RESOURCE"),
	}
}

// validateConfig checks that all required environment variables are set
func ValidateConfig(cfg Config) {
	missingVars := []string{}

	if cfg.AzureSubscriptionID == "" {
		missingVars = append(missingVars, "AZURE_SUBSCRIPTION_ID")
	}
	if cfg.AzureLocation == "" {
		missingVars = append(missingVars, "AZURE_LOCATION")
	}
	if cfg.AzureResourceGroupName == "" {
		missingVars = append(missingVars, "AZURE_RESOURCE_GROUP_NAME")
	}
	if cfg.AzureStorageAccountName == "" {
		missingVars = append(missingVars, "AZURE_STORAGE_ACCOUNT_NAME")
	}
	if cfg.AzureFunctionAppName == "" {
		missingVars = append(missingVars, "AZURE_FUNCTION_APP_NAME")
	}
	if cfg.FunctionName == "" {
		missingVars = append(missingVars, "FUNCTION_NAME")
	}
	if cfg.FunctionTemplate == "" {
		missingVars = append(missingVars, "FUNCTION_TEMPLATE")
	}
	if cfg.AuthLevel == "" {
		missingVars = append(missingVars, "AUTH_LEVEL")
	}

	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %v", missingVars)
	}

	log.Println("All required environment variables are set.")
}

// isCommandAvailable checks if a command is available in the system's PATH.
func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// shouldKeepResource determines whether to keep Azure resources based on KEEP_RESOURCE value
func ShouldKeepResource(keep string) bool {
	switch keep {
	case "1", "true", "True", "TRUE":
		return true
	default:
		return false
	}
}

// getAzureCredential initializes Azure SDK credentials
func GetAzureCredential() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// initializeClients initializes Azure SDK clients
func InitializeClients(ctx context.Context, cred *azidentity.DefaultAzureCredential, cfg Config) error {
	var err error
	resourcesClientFactory, err = armresources.NewClientFactory(cfg.AzureSubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("Failed to create resources client factory: %v", err)
	}
	resourceGroupClient = resourcesClientFactory.NewResourceGroupsClient()

	storageClientFactory, err = armstorage.NewClientFactory(cfg.AzureSubscriptionID, cred, nil)
	if err != nil {
		return fmt.Errorf("Failed to create storage client factory: %v", err)
	}
	accountsClient = storageClientFactory.NewAccountsClient()

	return nil
}

// createResourceGroup creates an Azure Resource Group
func CreateResourceGroup(ctx context.Context, cfg Config) (*armresources.ResourceGroup, error) {
	resourceGroupResp, err := resourceGroupClient.CreateOrUpdate(
		ctx,
		cfg.AzureResourceGroupName,
		armresources.ResourceGroup{
			Location: to.Ptr(cfg.AzureLocation),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &resourceGroupResp.ResourceGroup, nil
}

// checkNameAvailability checks if the storage account name is available
func CheckNameAvailability(ctx context.Context, cfg Config) (*armstorage.CheckNameAvailabilityResult, error) {
	result, err := accountsClient.CheckNameAvailability(
		ctx,
		armstorage.AccountCheckNameAvailabilityParameters{
			Name: to.Ptr(cfg.AzureStorageAccountName),
			Type: to.Ptr("Microsoft.Storage/storageAccounts"),
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &result.CheckNameAvailabilityResult, nil
}

// createStorageAccount creates an Azure Storage Account
func CreateStorageAccount(ctx context.Context, cfg Config) (*armstorage.Account, error) {
	pollerResp, err := accountsClient.BeginCreate(
		ctx,
		cfg.AzureResourceGroupName,
		cfg.AzureStorageAccountName,
		armstorage.AccountCreateParameters{
			Kind:     to.Ptr(armstorage.KindStorageV2),
			SKU:      &armstorage.SKU{Name: to.Ptr(armstorage.SKUNameStandardLRS)},
			Location: to.Ptr(cfg.AzureLocation),
			Properties: &armstorage.AccountPropertiesCreateParameters{
				AccessTier: to.Ptr(armstorage.AccessTierCool),
				Encryption: &armstorage.Encryption{
					Services: &armstorage.EncryptionServices{
						File:  &armstorage.EncryptionService{KeyType: to.Ptr(armstorage.KeyTypeAccount), Enabled: to.Ptr(true)},
						Blob:  &armstorage.EncryptionService{KeyType: to.Ptr(armstorage.KeyTypeAccount), Enabled: to.Ptr(true)},
						Queue: &armstorage.EncryptionService{KeyType: to.Ptr(armstorage.KeyTypeAccount), Enabled: to.Ptr(true)},
						Table: &armstorage.EncryptionService{KeyType: to.Ptr(armstorage.KeyTypeAccount), Enabled: to.Ptr(true)},
					},
					KeySource: to.Ptr(armstorage.KeySourceMicrosoftStorage),
				},
			},
		}, nil)
	if err != nil {
		return nil, err
	}
	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.Account, nil
}

// storageAccountProperties retrieves properties of the Storage Account
func StorageAccountProperties(ctx context.Context, cfg Config) (*armstorage.Account, error) {
	storageAccountResponse, err := accountsClient.GetProperties(
		ctx,
		cfg.AzureResourceGroupName,
		cfg.AzureStorageAccountName,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &storageAccountResponse.Account, nil
}

// initializeFunctionProject initializes a new Azure Functions project if not already initialized
func InitializeFunctionProject() error {
	// Check if the project directory exists
	if _, err := os.Stat(functionProjectDir); os.IsNotExist(err) {
		// Create the project directory
		err := os.MkdirAll(functionProjectDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}
	}

	// Change to the project directory
	err := os.Chdir(functionProjectDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to project directory: %v", err)
	}

	// Initialize a new Functions project with Node.js runtime
	// This step is optional if your project is already initialized
	cmd := exec.Command("func", "init", "--worker-runtime", "node")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("func init failed: %v\nOutput: %s", err, string(output))
	}

	log.Printf("func init output:\n%s\n", string(output))
	return nil
}

// createNewFunction creates a new Azure Function using `func new`
func CreateNewFunction(cfg Config) error {
	// Ensure we are in the project directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	log.Println("Current Directory:", currentDir)

	// Define the arguments for `func new`
	cmdArgs := []string{
		"new",
		"--name", cfg.FunctionName,
		"--template", cfg.FunctionTemplate,
		"--authlevel", cfg.AuthLevel,
	}

	cmd := exec.Command("func", cmdArgs...)

	// Set environment variables if needed
	cmd.Env = os.Environ()

	// Capture standard output and error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("func new failed: %v\nOutput: %s", err, string(output))
	}

	log.Printf("func new output:\n%s\n", string(output))
	return nil
}

// createFunctionApp creates an Azure Function App using `az functionapp create`
func CreateFunctionApp(cfg Config) error {
	cmdArgs := []string{
		"functionapp", "create",
		"--resource-group", cfg.AzureResourceGroupName,
		"--consumption-plan-location", cfg.AzureLocation,
		"--runtime", "node",
		"--runtime-version", "18",
		"--functions-version", "4",
		"--name", cfg.AzureFunctionAppName,
		"--storage-account", cfg.AzureStorageAccountName,
	}

	cmd := exec.Command("az", cmdArgs...)

	// Set environment variables if needed (e.g., AZURE_SUBSCRIPTION_ID)
	cmd.Env = os.Environ()

	// Capture standard output and error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("az functionapp create failed: %v\nOutput: %s", err, string(output))
	}

	log.Printf("az functionapp create output:\n%s\n", string(output))
	return nil
}

// publishFunctionApp publishes the Function App using `func azure functionapp publish`
func PublishFunctionApp(cfg Config) error {
	// Ensure you are in the Function App project directory
	err := os.Chdir(functionProjectDir)
	if err != nil {
		return fmt.Errorf("failed to change directory to project directory: %v", err)
	}

	cmdArgs := []string{
		"azure", "functionapp", "publish", cfg.AzureFunctionAppName,
	}

	cmd := exec.Command("func", cmdArgs...)

	// Set environment variables if needed
	cmd.Env = os.Environ()

	// Capture standard output and error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("func azure functionapp publish failed: %v\nOutput: %s", err, string(output))
	}

	log.Printf("func azure functionapp publish output:\n%s\n", string(output))
	return nil
}

// cleanup deletes the Resource Group to clean up resources
func Cleanup(ctx context.Context, cfg Config) error {
	pollerResp, err := resourceGroupClient.BeginDelete(ctx, cfg.AzureResourceGroupName, nil)
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
