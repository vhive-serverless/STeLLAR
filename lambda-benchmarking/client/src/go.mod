module lambda-benchmarking/client

go 1.15

replace lambda-benchmarking/client/prompts => ./prompts

replace lambda-benchmarking/client/configuration => ./configuration

replace lambda-benchmarking/client/experiment => ./experiment

replace lambda-benchmarking/client/experiment/benchmarking => ./experiment/benchmarking

replace lambda-benchmarking/client/experiment/visualization => ./experiment/visualization

replace lambda-benchmarking/client/experiment/networking => ./experiment/networking

require (
	github.com/aws/aws-sdk-go v1.35.20
	github.com/go-gota/gota v0.10.1
	github.com/sirupsen/logrus v1.7.0
	gonum.org/v1/gonum v0.8.1
	gonum.org/v1/plot v0.8.0
)
