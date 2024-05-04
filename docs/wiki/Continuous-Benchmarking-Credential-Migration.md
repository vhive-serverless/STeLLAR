# STeLLAR Offboarding / Credential Migration (Continuous Benchmarking)

This article contains a checklist of accounts/credentials that needs to be updated prior to the completion of FYPs to ensure the continued operation of daily experiments and CI/CD operations on STeLLAR.

## AWS

No action is needed for AWS as all credentails are under EASELab's account.

## Azure

### Azure Credentials on CI/CD Pipeline

The following values are needed for STeLLAR to deploy to Azure on the CI/CD pipeline:
- `AZURE_SUBSCRIPTION_ID`
- `AZURE_TENANT_ID`
- `AZURE_CLIENT_ID`
- `AZURE_CLIENT_SECRET`

You should already have an Azure account with a `AZURE_SUBSCRIPTION_ID`. The remaining values can be obtained via the Azure CLI or via the Azure console.

### Option 1 - Obtaining Azure credentials via the Azure CLI

[Create an Azure Service Principal](https://learn.microsoft.com/en-us/cli/azure/ad/sp?view=azure-cli-latest#az-ad-sp-create-for-rbac) with a secret via the Azure CLI:

```
az ad sp create-for-rbac --name "STeLLAR GitHub Actions" \
--role contributor \
--scopes /subscriptions/<your_azure_subscription_id> \
--sdk-auth
```

The command will output the credentials in JSON format.

### Option 2 - Obtaining Azure credentials via the Azure console

1. [Register an application on Microsoft Entra ID](https://learn.microsoft.com/en-sg/entra/identity-platform/howto-create-service-principal-portal#register-an-application-with-microsoft-entra-id-and-create-a-service-principal). The optional fields can be left empty. The `AZURE_CLIENT_ID` and `AZURE_TENANT_ID` will be displayed on the dashboard of the application.

<p align="center">
   <img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/9efc6e69-e63c-4596-b620-1f95dffa7da9"/>
</p>


2. [Assign a role to the application](https://learn.microsoft.com/en-sg/entra/identity-platform/howto-create-service-principal-portal#assign-a-role-to-the-application). The application should have a “Contributor” role.

3. [Create a new client secret for the application](https://learn.microsoft.com/en-sg/entra/identity-platform/howto-create-service-principal-portal#option-3-create-a-new-client-secret). The `AZURE_CLIENT_SECRET` will be displayed.

Finally, add the four values you have obtained as a secret on the STeLLAR repository.

### Azure self-hosted runners

Create an Azure VM with the following configuration:
- **Region:** (US) West US
   - **Note:** This is the current region all of our VMs are running at. Benchmarked Azure Functions are also deployed to this region.
- **Image:** Ubuntu Server 22.04 LTS - x64 Gen2
- **VM architecture:** x64
- **Size:** Standard_B1ms - 1 vcpu, 2GiB memory (US$18.10/month)
   - **Note:** The VM with larger 2GiB memory is required for image size experiments. Smaller VMs are known to run out of memory and crash when executing 100MB experiments.
- **OS disk type:** Standard HDD
   - **Note:** Standard HDD is cheaper and generally sufficient for our experiment needs.
- **Delete NIC when VM is deleted:** Checked
   - **Note:** Optional. Enabling this makes resource cleanup easier if you need to remove this self-hosted runner in the future.
- You may use the default options for those that were not specified above.

Execute the setup script to install the STeLLAR dependencies:

```
chmod +x ./scripts/setup.sh
./scripts/setup.sh
```

[Add the VM you created as a self-hosted runner for GitHub Actions.](https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners/adding-self-hosted-runners#adding-a-self-hosted-runner-to-a-repository)

The final `./run.sh` command in GitHub’s instructions to set up the self-hosted runner should be executed in a tmux terminal so that it can continue running after the ssh session ends.

## Google Cloud Run (GCR)

Google Cloud Run’s STeLLAR project (Project number: 567992090480, Project ID: stellar-benchmarking) is under EASELab’s organisation account (stellarbench@gmail.com), and the credentials are linked to the following service account (viewable under IAM dashboard in GCloud):

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/e79e8508-daed-4a09-b2b9-228ec2137c08"/>
</p>


### Changing Billing Account

However, a billing account needs to be added to the organisation in order to deploy GCR.
Search “Create a new billing account” in the GCloud console search bar:

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/10ba600c-6d56-4d1f-83c5-d22042122ea1"/>
</p>

Afterwards, the project needs to be shifted from the old billing account to the new billing account. Google has provided a comprehensive tutorial about the process [on this page](https://cloud.google.com/billing/docs/how-to/modify-project). (See “Change the Cloud Billing account linked to a project”)

### Google Cloud Self-hosted Runners

Self-hosted runners for GCR experiments use Google Cloud’s [Compute Engine](https://cloud.google.com/products/compute?hl=en) service. The list of all running Compute Engines can be seen on GCloud Console (search for “Compute Engine on the search bar):

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/0e940882-6ae0-491f-844b-340d22e74edd"/>
</p>

### Creating a new Compute Engine VM

To minimise setup and configuration of dependencies needed for STeLLAR, a machine image based on the existing configurations and setup has been set up which can be used to create new VMs. Choose “New VM instance from machine image” after selecting “CREATE INSTANCE”:

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/528e4532-5a85-4207-aa02-0c65e1b8fdca"/>
</p>

Choose the following configuration:
- **Region:** us-west1
- **Zone:** us-west1-a
- **Machine Configuration:** e2
- **Machine Type:** **e2-small** (e2-micro is unable to handle Java function builds and may fail Java deployments to GCR)

**Reference configuration:**
<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/41e99767-dda4-48d6-a02d-819847ddb553"/>
</p>

### Increasing VM Disk Size

Note that 16GB of disk space may not be enough and may risk running into disk space issues as experiments run over a long period of time. It is recommended to extend this to **64GB** of disk size.

Select the VM you wish to change from the Compute Engine dashboard and click on the boot disk in the Boot Disk section:
<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/03fd787b-23b6-42ff-bb19-92401322c315"/>
</p>

Select “Edit” and input your new Disk Size:
<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/6e59e236-5a58-41d1-aeb6-3e6f6111043b"/>
</p>

## Cloudflare

Deployment of Cloudflare workers requires a [Cloudflare Account](https://www.cloudflare.com/en-gb/) as well as a single API Token.

### New API Token

To create an API Token, go to “My Profile” > “API Tokens” > “Create” :

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/cc18ce7b-43ee-41a9-a946-ffbabf18a663"/>
</p>

Select the pre-defined template for editing Cloudflare Workers:

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/ef5e4497-fefa-418d-8fe4-099b8d693a3f"/>
</p>

Update the credential `CLOUDFLARE_API_TOKEN` under Github Actions secrets to deploy Cloudflare workers under the new account and API token:

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/62bf7b0d-16b9-4f6f-a0e2-e68cc7e86792"/>
</p>

### Cloudflare self-hosted Runners

As Cloudflare does not currently offer a service for Virtual machine computing services, AWS EC2 instances are used instead.

EC2 instances running Cloudflare experiments are located in **us-east-2 (Ohio)** as it was determined to be the closest to the deployed Cloudflare Workers. This may or may not change in the future, and a check on the geographical location of deployed Cloudflare Workers is recommended when deploying from the new account.

<p align="center">
<img src="https://github.com/vhive-serverless/STeLLAR/assets/76023265/8528322d-fe13-43e2-9b60-7181584d0497"/>
</p>
