package golang

import (
	"github.com/epicchainlabs/epicchain-go/pkg/smartcontract/binding"
	"github.com/epicchainlabs/epicchain-go/pkg/smartcontract/rpcbinding"
)

func goOffChainConfig() binding.Config {
	return rpcbinding.NewConfig()
}

func goOffChainGenerate() generateFunction {
	return rpcbinding.Generate
}
