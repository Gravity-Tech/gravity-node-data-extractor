package susy

import (
	"context"
	"fmt"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/contracts/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
	"time"
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
	MaxRqTimeout = 20

	WavesDecimals = 8
	EthDecimals   = 18
)

type ExtractImplementation ExtractionProvider
type ExtractImplementationType int

const (
	WavesSourceLock ExtractImplementationType = iota
	EthereumSourceBurn
)

type SourceExtractor struct {
	implementation ExtractImplementationType
	provider    ExtractionProvider

	cache       map[RequestId]time.Time
	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper
	luContract  string
	ibContract  *ibport.IBPort
}

func pickExtractionProvider(impl ExtractImplementationType, extractor *SourceExtractor) ExtractionProvider {
	switch impl {
	case WavesSourceLock:
		return &WavesExtractionProvider{ ExtractorDelegate:extractor }
	case EthereumSourceBurn:
		return &EthereumExtractionProvider{ ExtractorDelegate:extractor }
	}

	return nil
}

func New(sourceNodeUrl string, destinationNodeUrl string, luAddress string, ibAddress string, ctx context.Context, impl ExtractImplementationType) (*SourceExtractor, error) {
	ethClient, err := ethclient.DialContext(ctx, destinationNodeUrl)
	if err != nil {
		return nil, err
	}
	destinationContract, err := ibport.NewIBPort(common.HexToAddress(ibAddress), ethClient)
	if err != nil {
		return nil, err
	}
	wavesClient, err := client.NewClient(client.Options{BaseUrl: sourceNodeUrl})
	if err != nil {
		return nil, err
	}

	extractor := &SourceExtractor{
		implementation: impl,
		cache:          make(map[RequestId]time.Time),
		ethClient:      ethClient,
		wavesClient:    wavesClient,
		wavesHelper:    helpers.NewClientHelper(wavesClient),
		ibContract:     destinationContract,
		luContract:     luAddress,
	}

	extractor.provider = pickExtractionProvider(impl, extractor)

	return extractor, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	switch e.implementation {
	case WavesSourceLock:
		return &extractors.ExtractorInfo{
			Tag:         "source-waves",
			Description: "Source WAVES",
		}
	case EthereumSourceBurn:
		return &extractors.ExtractorInfo{
			Tag:         "source-eth",
			Description: "Source Ethereum",
		}
	}

	return nil
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {


	return nil, fmt.Errorf("No impl available")
}

func (e *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}
