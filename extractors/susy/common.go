package susy

import (
	"context"
	"fmt"
	"math/big"
	"time"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	//"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
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

type RequestChainKind int
const (
	WavesSource RequestChainKind = iota
	EthereumSource
)
type RequestDirection struct {
	Kind     RequestChainKind
	IsDirect bool
}

type ExtractImplementation ExtractionProvider
type ExtractImplementationType int

const (
	WavesSourceLock ExtractImplementationType = iota
	EthereumSourceBurn
)

type SourceExtractor struct {
	implementation ExtractImplementationType
	provider       ExtractionProvider

	cache       map[RequestId]time.Time
	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper
	luContract  string
	ibContract  string

	sourceDecimals      int64
	destinationDecimals int64

	requestConfig RequestDirection
}

func pickExtractionProvider(dir RequestDirection, extractor *SourceExtractor) ExtractionProvider {
	switch dir.Kind {
	case WavesSource:
		return &WavesExtractionProvider{ ExtractorDelegate: extractor, IsDirect: dir.IsDirect }
	case EthereumSource:
		return &EthereumExtractionProvider{ ExtractorDelegate: extractor, IsDirect: dir.IsDirect }
	}

	return nil
}

func New(sourceNodeUrl string, destinationNodeUrl string,
	luAddress string, ibAddress string,
	sourceDecimals int64, destinationDecimals int64,
	ctx context.Context, requestConfig RequestDirection) (*SourceExtractor, error) {
	ethClient, err := ethclient.DialContext(ctx, destinationNodeUrl)
	if err != nil {
		return nil, err
	}

	wavesClient, err := client.NewClient(client.Options{BaseUrl: sourceNodeUrl})
	if err != nil {
		return nil, err
	}

	extractor := &SourceExtractor{
		cache:               make(map[RequestId]time.Time),
		ethClient:           ethClient,
		wavesClient:         wavesClient,
		wavesHelper:         helpers.NewClientHelper(wavesClient),
		ibContract:          ibAddress,
		luContract:          luAddress,
		sourceDecimals:      sourceDecimals,
		destinationDecimals: destinationDecimals,
		requestConfig:       requestConfig,
	}

	extractor.provider = pickExtractionProvider(requestConfig, extractor)

	return extractor, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	switch e.requestConfig.Kind {
	case WavesSource:
		return &extractors.ExtractorInfo{
			Tag:         "source-waves",
			Description: "Source WAVES",
		}
	case EthereumSource:
		return &extractors.ExtractorInfo{
			Tag:         "source-eth",
			Description: "Source Ethereum",
		}
	}

	return nil
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	if e.provider != nil {
		return e.provider.Extract(ctx)
	}

	return nil, fmt.Errorf("No impl available")
}

func (e *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}
