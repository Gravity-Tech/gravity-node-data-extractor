package susy

import (
	"context"
	"fmt"

	extcfg "github.com/Gravity-Tech/gravity-node-data-extractor/v2/config"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/bridge"
)

type ExtractionProvider interface {
	Extract(context.Context) (*extractors.Data, error)
}

type SourceExtractor struct {
	delegate bridge.ChainExtractionBridge
	gateway *DirectedGateway
}

func New(cfg *extcfg.MainConfig, gateway *DirectedGateway) (*SourceExtractor, error) {
	delegate := gateway.BuildDelegate(cfg.IntoBridge())

	if delegate == nil {
		return nil, fmt.Errorf("no impl available")
	}

	return &SourceExtractor{
		delegate,
		gateway,
	}, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	return &extractors.ExtractorInfo{
		Description: "cross-chain extractor",
		Tag:         "direct extractor",
	}
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	var result *extractors.Data
	var err error

	// switch e.gateway.Kind() {
	// case WavesToEthDirect, EthToWavesDirect, EVMToSolanaDirect:
	// 	result, err = e.delegate.ExtractDirectTransferRequest(ctx)
	// case WavesToEthReverse, EthToWavesReverse, EVMToSolanaReverse:
	// 	result, err = e.delegate.ExtractReverseTransferRequest(ctx)
	// }

	if err != nil {
		// debug.PrintStack()
		return nil, err
	}
	if result != nil {
		return result, nil
	}
	
	return nil, fmt.Errorf("no corresponding implementation available")
}


func (e *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}
