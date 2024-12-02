package golang

import (
	"github.com/epicchainlabs/epicchain-go/pkg/smartcontract/binding"
)

func goOnChainConfig() binding.Config {
	return binding.NewConfig()
}

func goOnChainGenerate() generateFunction {
	return binding.Generate
}
