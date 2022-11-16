module main

go 1.19

replace (
	github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/common => ./common
	github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen => ./proto_gen
)

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/common v0.0.0-00010101000000-000000000000
)
