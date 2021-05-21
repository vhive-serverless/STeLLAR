module main

go 1.16

replace (
	github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/functions/producer-consumer/common => ./common
	github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen => ./proto_gen
)

require (
	github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/functions/producer-consumer/common v0.0.0-00010101000000-000000000000
	github.com/ease-lab/vhive-bench/src/setup/deployment/raw-code/functions/producer-consumer/proto_gen v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.35.0
)
