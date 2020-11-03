package configuration

type ExperimentConfig struct {
	Bursts                  int
	BurstSizes              []string
	PayloadLengthBytes      int
	FrequencySeconds        float64
	FunctionIncrementLimits []int64 // If more than one, service time is dynamic
	GatewayEndpoints        []string
	Id                      int
	IatType                 string
	Provider                string
}
