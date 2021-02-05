package susy

import (
	"context"
	"fmt"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/bridge"
)

const (
	WavesToEthDirect     extractors.ExtractorType = "waves-based-to-eth-direct"
	WavesToEthReverse    extractors.ExtractorType = "waves-based-to-eth-reverse"
	EthToWavesDirect     extractors.ExtractorType = "eth-based-to-waves-direct"
	EthToWavesReverse    extractors.ExtractorType = "eth-based-to-waves-reverse"
)

type ExtractionProvider interface {
	Extract(context.Context) (*extractors.Data, error)
}

type SourceExtractor struct {
	kind extractors.ExtractorType
	delegate bridge.ChainExtractionBridge
}

func New(sourceNodeUrl string, destinationNodeUrl string,
	luAddress string, ibAddress string,
	sourceDecimals int64, destinationDecimals int64, kind extractors.ExtractorType) (*SourceExtractor, error) {

	var delegate bridge.ChainExtractionBridge
	config := bridge.ConfigureCommand{
		SourceNodeUrl: sourceNodeUrl,
		DestinationNodeUrl: destinationNodeUrl,
		LUPortAddress: luAddress,
		IBPortAddress: ibAddress,
		SourceDecimals: sourceDecimals,
		DestinationDecimals: destinationDecimals,
	}

	switch kind {
		case WavesToEthDirect, WavesToEthReverse: {
			delegate = &bridge.WavesToEthereumExtractionBridge{}
		}
	}
	if delegate == nil {
		return nil, fmt.Errorf("no impl available")
	}

	err := delegate.Configure(config)
	if err != nil {
		return nil, err
	}

	extractor := &SourceExtractor{
		kind: kind,
		delegate: delegate,
	}

	return extractor, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	return nil
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	switch e.kind {
	case WavesToEthDirect, EthToWavesDirect:
		return e.delegate.ExtractDirectTransferRequest(ctx)
	case WavesToEthReverse, EthToWavesReverse:
		return e.delegate.ExtractReverseTransferRequest(ctx)
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
