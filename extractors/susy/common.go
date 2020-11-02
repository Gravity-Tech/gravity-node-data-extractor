package susy

import (
	"context"
	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
	"time"
)

const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	SuccessEthereum // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 5 * 60 // 5 min

	WavesDecimals = 8
	EthDecimals   = 18
)

type ExtractImplementationType int

const (
	WavesSourceLock ExtractImplementationType = iota
	EthereumSourceBurn
)

/**
 * Builder pattern in action.
 * Similar options object can be used in different implementations of extractor (different extractor builders).
 *
 * Available builders:
 * - New(*susy.WavesEthereumBridgeOptions) *SourceExtractor
 * - New(*susy.WavesEthereumBridgeOptions) *DestinationExtractor
 */
type WavesEthereumBridgeOptions struct {
	Cache       map[RequestId]time.Time
	EthClient   *ethclient.Client
	WavesClient *client.Client
	WavesHelper helpers.ClientHelper
	LUContract  string
	IBContract  *ibport.IBPort
}


func NewOptions(sourceNodeUrl string, destinationNodeUrl string, luAddress string, ibAddress string, ctx context.Context) (*WavesEthereumBridgeOptions, error) {
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

	return &WavesEthereumBridgeOptions{
		Cache:          make(map[RequestId]time.Time),
		EthClient:      ethClient,
		WavesClient:    wavesClient,
		WavesHelper:    helpers.NewClientHelper(wavesClient),
		IBContract:     destinationContract,
		LUContract:     luAddress,
	}, nil
}