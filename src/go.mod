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
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/aws/aws-sdk-go v1.41.9
	github.com/go-fonts/liberation v0.3.1 // indirect
	github.com/go-gota/gota v0.12.0
	github.com/go-latex/latex v0.0.0-20230307184459-12ec69307ad9 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/image v0.11.0 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	gonum.org/v1/gonum v0.14.0
	gonum.org/v1/plot v0.14.0
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

require (
	git.sr.ht/~sbinet/gg v0.5.0 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-pdf/fpdf v0.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)
