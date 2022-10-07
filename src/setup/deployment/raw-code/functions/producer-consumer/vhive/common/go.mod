module common

go 1.16

replace github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen => ../proto_gen

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.37.6
	github.com/vhive-serverless/stellar/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen v0.0.0-00010101000000-000000000000
	github.com/minio/minio-go/v7 v7.0.9
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.35.0
)
