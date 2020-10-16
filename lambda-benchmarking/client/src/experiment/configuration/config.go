package configuration

type ExperimentConfig struct {
	Bursts               int
	BurstSizes           []string
	PayloadLengthBytes   int
	FrequencySeconds     float64
	LambdaIncrementLimit []string
	GatewayEndpoints     []string
	Id                   int
	IatType              string
}
