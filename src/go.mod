module vhive-bench

go 1.17

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
	github.com/ajstarks/svgo v0.0.0-20210923152817-c3b6e2f0c527 // indirect
	github.com/aws/aws-sdk-go v1.42.44
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-fonts/liberation v0.2.0 // indirect
	github.com/go-gota/gota v0.12.0
	github.com/go-latex/latex v0.0.0-20210823091927-c0d11ff05a81 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f // indirect
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3 // indirect
	golang.org/x/text v0.3.7 // indirect
	gonum.org/v1/gonum v0.9.3
	gonum.org/v1/plot v0.10.0
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-pdf/fpdf v0.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)
