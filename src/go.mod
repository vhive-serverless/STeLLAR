module stellar

go 1.19

replace (
	stellar/benchmarking => ./benchmarking
	stellar/benchmarking/networking => ./benchmarking/networking
	stellar/benchmarking/networking/benchgrpc/proto_gen => ./benchmarking/networking/benchgrpc/proto_gen
	stellar/benchmarking/visualization => ./benchmarking/visualization
	stellar/benchmarking/writers => ./benchmarking/writers
	stellar/setup => ./setup
	stellar/util => ./util
)

require (
	github.com/ajstarks/svgo v0.0.0-20210923152817-c3b6e2f0c527 // indirect
	github.com/aws/aws-sdk-go v1.41.9
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/go-fonts/liberation v0.2.0 // indirect
	github.com/go-gota/gota v0.12.0
	github.com/go-latex/latex v0.0.0-20210823091927-c0d11ff05a81 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.3
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gonum.org/v1/gonum v0.9.3
	gonum.org/v1/plot v0.10.0
	google.golang.org/grpc v1.61.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-pdf/fpdf v0.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
)
