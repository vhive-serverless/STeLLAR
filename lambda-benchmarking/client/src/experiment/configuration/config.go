package configuration

type ExperimentConfig struct {
	Bursts               int
	BurstSizes           []string
	PayloadLengthBytes   int
	FrequencySeconds     int
	LambdaIncrementLimit []string
	GatewayEndpoints     []string
	Id                   int
}