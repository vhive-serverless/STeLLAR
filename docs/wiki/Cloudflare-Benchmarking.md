## Introduction

Cloudflare Workers allow for deployment of Javascript/WASM functions which can be accessed via a generated endpoint. 

## Pre-requisites

1. [Cloudflare Account](https://www.cloudflare.com) for deploying and deleting Cloudflare Workers.

## Setup

To run experiments with STeLLAR for Cloudflare, the `wrangler` CLI tool must be set up.

1. Run the setup script to install dependencies required for running STeLLAR at `scripts/setup.sh`
2. Run the setup script for Cloudflare located at `scripts/cloudflare/setup.sh`  
(Note: You may need to change permissions of the script to allow execution)

Alternatively, if you are using a headless machine setup with no GUI/Web browser for the `wrangler login` command used in `setup.sh`, set the `CLOUDFLARE_API_TOKEN` environment value by [creating an API token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token).

In addition, STeLLAR requires two core components to deploy and benchmark serverless functions: The function code and a JSON file specifying experiment parameters.

### Function code

Put your code directory at `src/setup/deployment/raw-code/serverless/cloudflare/<function_code_dir>`.

Cloudflare workers have native support for V8 Javascript and WebAssembly runtimes; other runtimes must be transpiled to Javascript first (See [here](https://blog.cloudflare.com/cloudflare-workers-announces-broad-language-support) for more details).

### Experiment JSON file

The JSON file specifying configurations can be placed at any path. The folder for our examples can be located under the `experiments/` folder. For examples on how to write the experiment JSON file, see [here](https://github.com/vhive-serverless/STeLLAR/tree/main/experiments/tests/cloudflare).

Summary of experiment JSON file values:

| Key | Type | Description |
| --- | --- | --- |
| Sequential | boolean | Specifies whether to run sub-experiments sequentially or in parallel |
| Provider | string | Specifies the Cloud Provider to deploy the function to. (Possible values: `aws`, `azure`, `gcr`, `cloudflare` |
| Runtime | string | Specifies the runtime of the function. (e.g. `go1.x`, `java11`, etc.) |
| SubExperiments | array | Array of objects specifying specific experiment parameters for each sub-experiment |
| Title | string | Name of the sub-experiment |
| Function | string | Name of the function. This must be the same value as `<function_code_dir>` |
| Handler | string | Entry point of the function. For Cloudflare, the value is the name of the script file (`index.js` for native Javascript and `dist/main.js` for transpiled runtimes) |
| PackageType | string | Specifies the type of packaging for function upload.  |
| Bursts | number | Specifies the number of bursts to send to the deployed function(s) |
| BurstSizes | array | Specifies the size of each burst when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst. |
| IATSeconds | number | Specifies the inter-arrival time between each burst. |
| DesiredServiceTimes | array | Specifies the desired service execution time(s) when invoking the deployed function(s). STeLLAR iterates and cycles through the array for each burst. These execution times are achieved by calculating the corresponding busy spin count for the desired time on the **host running the STeLLAR client**. |
| FunctionImageSizeMB | number | Specifies the target size of the function to upload. |
| Parallelism | number | Specifies the number of concurrent endpoints to deploy and benchmark. Useful for obtaining cold-start samples within a shorter period of time. |

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
