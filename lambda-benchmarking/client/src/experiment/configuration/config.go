package configuration

type ExperimentConfig struct {
	Bursts               int
	BurstSizes           []string
	PayloadLengthBytes   int
	FrequencySeconds     int
	LambdaIncrementLimit int
	GatewayEndpoints     []string
	Id                   int
}