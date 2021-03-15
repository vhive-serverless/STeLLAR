module vhive-bench/client

go 1.15

replace (
	vhive-bench/client/benchmarking => ./benchmarking
	vhive-bench/client/benchmarking/networking => ./benchmarking/networking
	vhive-bench/client/benchmarking/networking/benchgrpc/proto_gen => ./benchmarking/networking/benchgrpc/proto_gen
	vhive-bench/client/benchmarking/visualization => ./benchmarking/visualization
	vhive-bench/client/benchmarking/writers => ./benchmarking/writers
	vhive-bench/client/setup => ./setup
	vhive-bench/client/util => ./util
)

require (
	github.com/ajstarks/svgo v0.0.0-20200725142600-7a3c8b57fecb // indirect
	github.com/aws/aws-sdk-go v1.37.13
	github.com/go-gota/gota v0.10.1
	github.com/golang/protobuf v1.4.3
	github.com/sirupsen/logrus v1.8.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	gonum.org/v1/gonum v0.8.2
	gonum.org/v1/plot v0.9.0
	google.golang.org/genproto v0.0.0-20210212180131-e7f2df4ecc2d // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
