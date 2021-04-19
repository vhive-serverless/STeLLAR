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
	github.com/ajstarks/svgo v0.0.0-20200725142600-7a3c8b57fecb // indirect
	github.com/aws/aws-sdk-go v1.38.21
	github.com/go-gota/gota v0.10.1
	github.com/golang/protobuf v1.5.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20210319071255-635bc2c9138d // indirect
	gonum.org/v1/gonum v0.9.0
	gonum.org/v1/plot v0.9.0
	google.golang.org/genproto v0.0.0-20210319143718-93e7006c17a6 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)
