package susy

import "github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/bridge"

type Gateway uint64

type DirectedGateway struct {
	gateway  Gateway
	isDirect bool
}

const (
	WavesEVM Gateway = iota
	EVMWaves
	SolanaEVM
	EVMSolana
)

// func (gateway Gateway) ExtractionProvider() bridge.ChainExtractionBridge {
	
// }

func MatchGateway(pattern string) *DirectedGateway {
	for _, gate := range []Gateway { WavesEVM, EVMWaves, SolanaEVM, EVMSolana } {
		if pattern == gate.Direct() {
			return &DirectedGateway{ gateway: gate, isDirect: true }
		}
		if pattern == gate.Reverse() {
			return &DirectedGateway{ gateway: gate, isDirect: false }
		}
	}
	return nil
}

func (gateway DirectedGateway) ExtractionProvider() bridge.ChainExtractionBridge {
	switch gateway.gateway {
	case WavesEVM:
		return &bridge.WavesToEthereumExtractionBridge{}
	case EVMWaves:
		return &bridge.EthereumToWavesExtractionBridge{}
	case SolanaEVM:
		return &bridge.SolanaToEVMExtractionBridge{}
	case EVMSolana:
		return &bridge.EthereumToSolanaExtractionBridge{}
	}

	panic("no extraction provider")
}

func (gateway DirectedGateway) BuildDelegate(bridgeConfig bridge.ConfigureCommand) bridge.ChainExtractionBridge {
	provider := gateway.ExtractionProvider()
	provider.Configure(bridgeConfig)
	return provider
}

func (gateway Gateway) Direct() string {
	return gateway.name() + "-direct"
}

func (gateway Gateway) Reverse() string {
	return gateway.name() + "-reverse"
}

func (gateway Gateway) name() string {
	switch gateway {
		case WavesEVM:
			return "waves-based-to-eth"
		case EVMWaves:
			return "eth-based-to-waves"
		case SolanaEVM:
			return "solana-based-to-evm"
		case EVMSolana:
			return "evm-solana"
	}
	panic("not implemented")
}