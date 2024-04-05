## Introduction

Google Cloud Run allows for deployment of container applications which can be accessed via a generated endpoint. Note that only container-based deployment is supported (No ZIP deployment).

## Pre-requisites

1. Docker Hub Account with an access token (See [here](https://docs.docker.com/docker-hub/access-tokens/) on how to generate Docker Hub access tokens)
2. Google Cloud Account with [Google Cloud Run](https://cloud.google.com/run?hl=en) enabled

## Setup

To run experiments with STeLLAR for Google Cloud Run, the `gcloud` CLI tool must be set up.

1. Run the setup script to install dependencies required for running STeLLAR at `scripts/setup.sh`
2. Run the setup script for Google Cloud Run located at `scripts/gcr/setup.sh`  
(Note: You may need to change permissions of the script to allow execution)

To upload Docker images to Docker Hub for deployment with Google Cloud Run, the following environment values must also be set:
```
DOCKER_HUB_USERNAME=<your_docker_hub_username>
DOCKER_HUB_ACCESS_TOKEN=<your_docker_hub_access_token>
```

In addition, STeLLAR requires two core components to deploy and benchmark serverless functions: The function code and a JSON file specifying experiment parameters.

### Function code

Put your code directory at `src/setup/deployment/raw-code/serverless/gcr/<function_code_dir>`. Ensure that the Dockerfile is at the root of the `<function_code_dir>` directory. This function code must comply with [Google Cloud's Container runtime contract](https://cloud.google.com/run/docs/container-contract).

As Cloud Run is [based on Knative](https://cloud.google.com/blog/products/serverless/knative-based-cloud-run-services-are-ga), any Knative applications can be run out of the box on Cloud Run. See [here](https://knative.dev/docs/samples/serving/) for examples of code for various runtimes.

### Experiment JSON file

The JSON file specifying configurations can be placed at any path. The folder for our examples can be located under the `experiments/` folder. For examples on how to write the experiment JSON file, see [here](https://github.com/vhive-serverless/STeLLAR/tree/feature-serverless-framework-deployment/experiments/tests/gcr).

Summary of experiment JSON file values:

| Key | Type | Description |
| --- | --- | --- |
| Sequential | boolean | Specifies whether to run sub-experiments sequentially or in parallel |
| Provider | string | Specifies the Cloud Provider to deploy the function to. (Possible values: `aws`, `azure`, `gcr`, `cloudflare` |
| Runtime | string | Specifies the runtime of the function. (e.g. `go1.x`, `java11`, etc.) |
| SubExperiments | array | Array of objects specifying specific experiment parameters for each sub-experiment |
| Title | string | Name of the sub-experiment |
| Function | string | Name of the function. This must be the same value as `<function_code_dir>` |
| Handler | string | Entry point of the function. For GCR, the value must be `Dockerfile` |
| PackageType | string | Specifies the type of packaging for function upload. For GCR, the value must be `Container` |
| Bursts | number | Specifies the number of bursts to send to the deployed function(s) |
| BurstSizes | array | Specifies the size of each burst when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst. |
| IATSeconds | number | Specifies the interarrival time between each burst. |
| DesiredServiceTimes | array | Specifies the desired service execution time(s) when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst. These execution times are achieved by calculating the corresponding busy spin count for the desired time on the **host running the STeLLAR client**. |
| FunctionImageSizeMB | number | Specifies the target size of the function to upload. |
| Parallelism | number | Specifies the number of concurrent endpoints to deploy and benchmark. Useful for obtaining cold-start samples within a shorter period of time. |

## Benchmarking

1. Compile the STeLLAR binary at the `src` directory:

```
cd src
go build main.go
```

2. Run the compiled binary: `./main -o <output_folder_path> -c <experiment_json_file_path> -l <log_level>`  
(Note: You may need to add `sudo` at the front in order to access the Docker daemon)

Summary of the flags available:
| Flag | Description |
|------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-o` | Specifies the directory in which to store experiment results. Each experiment creates a folder named after the timestamp when it was created which is stored inside this directory. |
| `-c` | Specifies the path of the experiment JSON file which is used to run the experiments on. |
| `-l` | Specifies the log level to print to the console output. Default value is set to `info`. Possible values: `info`, `debug` |

### Obtaining Results

STeLLAR generates statistics as well as a CDF diagram after the experiments are completed. These results can be found under `<output_folder_path/<experiment_timestamp>`. It contains:

- `/<subexperiment_title>/latencies.csv`: All recorded request latencies for individual invocations
- `/<subexperiment_title>/statistics.csv`: Statistical information such as percentile, standard deviation, etc.
- `/<subexperiment_title>/empirical_CDF.png`: A CDF diagram of the sample latencies
