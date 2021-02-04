package susy

import (
	"context"
	"fmt"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/bridge"
	"math/big"
	"time"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
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

const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	SuccessEthereum // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 5 * 60 // 5 min
)

var (
	accuracy = big.NewInt(1).Exp(big.NewInt(10), big.NewInt(8), nil)
)

type SourceExtractor struct {
	delegate *bridge.ChainExtractionBridge
}

func New(sourceNodeUrl string, destinationNodeUrl string,
	luAddress string, ibAddress string,
	sourceDecimals int64, destinationDecimals int64,
	ctx context.Context, impl extractors.ExtractorType) (*SourceExtractor, error) {

	var delegate bridge.ChainExtractionBridge
	delegate = &bridge.WavesToEthereumExtractionBridge{}


	extractor := &SourceExtractor{
		delegate: delegate,
	}

	return extractor, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	switch e.kind {
	case WavesToEthDirect, WavesToEthReverse, EthToWavesDirect, EthToWavesReverse:
		return &extractors.ExtractorInfo{
			Description: string(e.kind),
			Tag:         string(e.kind),
		}
	}

	return nil
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	switch e.kind {
	case WavesToEthDirect, WavesToEthReverse, EthToWavesDirect, EthToWavesReverse:

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
