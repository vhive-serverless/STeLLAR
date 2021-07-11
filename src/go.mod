module vhive-bench

go 1.16

replace (
	vhive-bench/benchmarking => ./benchmarking
	vhive-bench/benchmarking/networking => ./benchmarking/networking
	vhive-bench/benchmarking/networking/benchgrpc/proto_gen => ./benchmarking/networking/benchgrpc/proto_gen
	vhive-bench/benchmarking/visualization => ./benchmarking/visualization
	vhive-bench/benchmarking/writers => ./benchmarking/writers
	vhive-bench/setup => ./setup
	vhive-bench/util => ./util
)

require (
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75 // indirect
	github.com/aws/aws-sdk-go v1.39.4
	github.com/go-fonts/liberation v0.2.0 // indirect
	github.com/go-gota/gota v0.11.0
	github.com/golang/protobuf v1.5.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gonum.org/v1/gonum v0.9.3
	gonum.org/v1/plot v0.9.0
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
