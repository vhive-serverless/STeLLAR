
# Database Operation Scripts

This directory contains python scripts to perform basic CRUD operations on dynamoDB instance to store multiple experiment results from STeLLAR.
Further, these functions are to be deployed in AWS lambda.

* `read-results-from-db/lambda_function.py` : Queries experiment results from database and returns the populated response. 
* `write-results-to-db/lambda_function.py` : Writes experiment results to the database and returns the written JSON object.

A single database entry may include the following fields.
 * `experiment_type` : type of the experiment (warm/cold/data-transfers etc.)
 * `date` : the date which the results are belong in
 * `subtype`: a text field to store the subtype of an experiment (if applicable)
 * `min`: minimum latency (ms)
 * `max`: maximum latency (ms)
 * `median`: median latency (ms)
 * `tail_latency`: tail latency (ms) - 99th percentile 
 * `first_quartile`: 25th percentile latency
 * `third_quartile`: 75th percentile latency
 * `standard_deviation`: standard deviation for the latencies
 * `payload_size`: size of the payload
 * `burst_size`: number of requests sent in a single burst
 * `IAT_type`: inter-arrival time type (long/short)
 * `count`: number of experiment samples
 * `provider`: name of the cloud provider (aws/google/azure)
 * `createdAt`: a system timestamp for the record