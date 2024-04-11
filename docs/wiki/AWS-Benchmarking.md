## Introduction

This guide provides instructions on how to deploy and benchmark AWS Lambda using STeLLAR.

## Pre-requisites

1. [AWS account](https://portal.aws.amazon.com/billing/signup#/start/email) with
   active [access keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey)
   and [AWSLambda_FullAccess](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AWSLambda_FullAccess.html)
   permissions.

## Setup

To deploy functions and run experiments on AWS Lambda through STeLLAR,
the [Serverless](https://www.serverless.com/framework/docs) framework must be
installed.

1. Run the setup script to install dependencies required for running STeLLAR (Note: You may need to change permissions
   of the script to allow execution):
   ```shell
   chmod +x ./scripts/setup.sh
   ./scripts/setup.sh
    ```

2. Setup the configuration and credential files for your AWS Command Line Interface by following the
   instructions [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-methods):
   ```shell
   aws configure
   ```

In addition, STeLLAR requires two core components to deploy and benchmark serverless functions: The function code and a
JSON file specifying experiment parameters.

### Function code

Put your code directory at `src/setup/deployment/raw-code/serverless/aws/<function_code_dir>`.

### Experiment JSON file

The JSON file specifying configurations can be placed at any path. The folder for our examples can be located under
the `experiments/` folder. For examples of experiment JSON files,
see [here](https://github.com/vhive-serverless/STeLLAR/blob/main/experiments/tests/aws).

Summary of experiment JSON file values:

| Key                 | Type    | Description                                                                                                                                                                                                                                                                                                |
|---------------------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Sequential          | boolean | Specifies whether to run sub-experiments sequentially or in parallel.                                                                                                                                                                                                                                      |
| Provider            | string  | Specifies the Cloud Provider to deploy the function to. (e.g. `aws`, `azure`, `gcr`, `cloudflare` etc.)                                                                                                                                                                                                    |
| Runtime             | string  | Specifies the runtime of the function. (e.g. `python3.9`, `nodejs18.x`, `java11` etc.)                                                                                                                                                                                                                     |
| SubExperiments      | array   | Array of objects specifying specific experiment parameters for each sub-experiment.                                                                                                                                                                                                                        |
| Title               | string  | Name of the sub-experiment.                                                                                                                                                                                                                                                                                |
| Function            | string  | Name of the function. This must be the same value as `<function_code_dir>`.                                                                                                                                                                                                                                |
| Handler             | string  | Entry point of the function. This should follow the format of `<your_file_name_without_extension>.<your_handler_function_name>`.                                                                                                                                                                           |
| PackageType         | string  | Specifies the type of packaging for function upload. Only `Zip` is currently accepted for STeLLAR benchmarking of AWS.                                                                                                                                                                                     |
| Bursts              | number  | Specifies the number of bursts to send to the deployed function(s).                                                                                                                                                                                                                                        |
| BurstSizes          | array   | Specifies the size of each burst when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst.                                                                                                                                                                     |
| DesiredServiceTimes | array   | Specifies the desired service execution time(s) when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst. These execution times are achieved by calculating the corresponding busy spin count for the desired time on the **host running the STeLLAR client**. |
| FunctionImageSizeMB | number  | Specifies the target size of the function to upload.                                                                                                                                                                                                                                                       |
| Parallelism         | number  | Specifies the number of concurrent endpoints to deploy and benchmark. Useful for obtaining cold-start samples within a shorter period of time.                                                                                                                                                             |

## Benchmarking

1. Compile the STeLLAR binary at the `src` directory:
   ```sh
   cd src
   go build main.go
   ```

2. Run the compiled binary:
   ```sh
   ./main -o <output_folder_path> -c <experiment_json_file_path> -l <log_level>
   ```

Summary of the flags available:

| Flag | Description                                                                                                                                                                         |
|------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-o` | Specifies the directory in which to store experiment results. Each experiment creates a folder named after the timestamp when it was created which is stored inside this directory. |
| `-c` | Specifies the path of the experiment JSON file which is used to run the experiments on.                                                                                             |
| `-l` | Specifies the log level to print to the console output. Default value is set to `info`. Possible values: `info`, `debug`                                                            |

### Obtaining Results

STeLLAR generates statistics as well as a CDF diagram after the experiments are completed. These results can be found
under `<output_folder_path/<experiment_timestamp>`. It contains:

- `/<subexperiment_title>/latencies.csv`: All recorded request latencies for individual invocations
- `/<subexperiment_title>/statistics.csv`: Statistical information such as percentile, standard deviation, etc.
- `/<subexperiment_title>/empirical_CDF.png`: A CDF diagram of the sample latencies
