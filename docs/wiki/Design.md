## Design

To begin with, we provide an overview of our benchmarking solution and define the main terms used throughout the rest of the codebase:

![design](../../design/diagram.png)

- The coordinator orchestrates the entire benchmarking procedure.
- The experiment configuration is an input JSON file used to specify and customize the experiments.
- An endpoint is a URL used for locating the function instance over the Internet. As seen in the diagram, this URL most often points to resources such as AWS API Gateway, Azure HTTP Triggers, vHive Kubernetes Load Balancer, or similar.
- The vendor endpoints input JSON file is only used for providers such as vHive that do not currently support automated function management (e.g., function listing, deployment, repurposing, or removal via SDKs or APIs).
- The inter-arrival time (IAT) is the time interval that the client waits for in-between sending two bursts to the same endpoint. To add some variability and simulate a more realistic scenario, we sample this from a shifted exponential distribution. For example, if we set the IAT to 10 minutes (modeling cold starts for most vendors), generated values can be, e.g., 10m12s, 10m27s, 11m.
- Multiple endpoints can be used simultaneously by the same experiment to speed up the benchmarking. The JSON configuration field parallelism defines this number: the higher it is, the more endpoints will be allocated, and the more bursts will be sent in short succession (speeding up the process for large IATs).
- The latencies CSV files are the main output of the evaluation framework. They are used in our custom Python plotting utility suite to produce insightful visualizations.
- The logs text file (Figure 4.2) is the final output of the benchmarking client. Log records are useful for optimizing code and debugging problematic behavior.

## Flow Chart
Finally, we look at the procedural steps adopted by the framework:

![flow chart](../../design/flow-chart.png)


1. The JSON configuration file is read and parsed, and any default field values are assigned. If the configuration file is missing, the program throws a fatal error.
2. Experiment service times (e.g., 10 seconds) are translated on the client machine into numbers representing busy-spin increment limits (e.g., 10,000,000). In turn, those are used by the measurement function on the server machine to keep the processor busy-spinning.
3. A connection with the serverless vendor is established. This is abstracted away behind a common interface having only four functions: ListAPIs, DeployFunction, RemoveFunction, and UpdateFunction. Used exclusively throughout the codebase, this interface offers seamless integration functionality with any provider.
4. In the provisioning phase, serverless.com framework is used to deploy the functions to the cloud and to establish HTTP endpoints. The functions are configured based on the experiment JSON file.
5. The last step runs all the experiments either sequentially or in parallel: bursts are successively sent to each available endpoint, followed by a sleep duration specified by the IAT. The process is repeated until all responses have been recorded to disk. Finally, statistics and visualizations are generated.

## Serverless.com Framework Deployment

Stellar uses serverless.com framework to deploy serverless functions to cloud.

![deployment](../../design/deployment.jpg)

1. The JSON experiment configuration file is parsed and serverless.yml service configuration file is written.
2. The function source code is compiled if needed. (e.g. Java and Go functions)
3. The filler file is created to increase the function deployment size. Filler files are added to the deployment package to benchmark performance of different functions with different sizes.
4. Based on the experiment deployment method:
    1. ZIP: serverless.com framework is able to zip the function with the filler file and no further steps are needed from STeLLAR. For better control over what gets zipped we allow the user to create "artifacts" - that is that STeLLAR zips the function and is provided to serverless.com already zipped. In such case, the serverless.com does not perform the zipping.
    2. Docker: Docker image is built - for this docker file needs to be provided by the user.
5. Serverless.com framework deploys the service defined in the service.yml. (`serverless deploy`)
6. Serverless.com return a list of endpoints and routes for every function defined.
7. Benchmarking is performed.
8. Once all the experiments are finished, the service is removed using the serverless.com framework. (`serverless remove`)

### Deployment boilerplate outline

Below are listed key packages and components that are most essential to serverless.com framework deployment.

- src/setup/
    - extract-configuration.go
        - Reads the experiment.json file and assigns values accordingly to Go Configuration struct.
    - assign-endpoints.go
        - Assigns endpoints and routes to each function deployed. Endpoint values are extracted from the Serverless
          framework deploy message and stored in the Configuration struct.
    - serverless-config.go
        - Creates the serverless.yml file based on the data in the experiment json file.

- src/setup/building/
    - This package takes builds Go binaries and Gradle projects for Java deployment.
- src/setup/code-generation/
    - This package generates Hello World provider-specific function source-code used for deployment.
- src/setup/deployment/packaging/
    - This package creates filler files, zips function packages and creates artifacts.
- src/setup/deployment/connection/
    - This package contains legacy code from the original deployment implementation, which used provider specific APIs.
      Most of this package will be eventually obsolete and deprecated, but we may get use of some of its
      functionalities.
- src/setup/deployment/raw-code/serverless/
    - This directory contains the logic of the functions that are to be deployed.
    - Each subdirectory contains source code resources for deploying to a specific provider (aws, azure, alibaba, gcr).
    - The serverless.yml file is stored inside each provider's subdirectory.
    - Function source code is stored in each provider’s subdirectories (e.g. aws/hellogo, azure/hellopy).

## Data Transfer Measurement
We integrate all necessary server-side functionality into a single function that we call a _measurement function_. This approach is similar to that taken in [40] and other serverless performance evaluation frameworks. A measurement function can perform up to three tasks, depending on the use case:

1. It always collects function instance runtime information.
2. If applicable, the function will simulate work by incrementing a variable in a busy-spin loop. This can be as simple as “for i := 0; i\<incrementLimit; i++{}”.
3. If applicable, the function records invocation timing. This is particularly useful for our data transfer studies where we complement client-measured round-trip time with internal function timestamps for validation purposes.

![transfer method](../../design/transfer-method.png)
