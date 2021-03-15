module main

go 1.15

replace (
	github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/proto_gen => ./../vhive/proto_gen
	github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/common => ./../vhive/common
)

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/ease-lab/vhive-bench/client/src/setup/deployment/raw-code/producer-consumer/go1.x/vhive/common v0.0.0-00010101000000-000000000000
)
