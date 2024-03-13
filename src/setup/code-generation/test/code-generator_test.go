package code_generation

import (
	code_generation "stellar/setup/code-generation"
	"testing"
)

func TestAddFunctionConfig(t *testing.T) {
	code_generation.GenerateCode("hellopy", "aws")
}
