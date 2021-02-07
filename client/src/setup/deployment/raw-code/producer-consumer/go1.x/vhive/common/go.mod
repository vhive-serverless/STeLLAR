module common

go 1.15

replace github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen => ../proto_gen

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.37.6
	github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.35.0
)
