package bridge

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	solclient "github.com/portto/solana-go-sdk/client"
)

// type EVMToSolanaExtractionBridge

type SolanaToEVMExtractionBridge struct {
	config     ConfigureCommand
	configured   bool

	ethClient      *ethclient.Client

	ibportContract *ibport.IBPort
	
	solanaClient   *solclient.Client
	solanaCtx       context.Context

	// IBPortDataAccount string
	LUPortDataAccount string
}

func (provider *SolanaToEVMExtractionBridge) Configure(config ConfigureCommand) error {
	if provider.configured {
		return fmt.Errorf("bridge is configured already")
	}

	provider.config = config

	// Node clients instantiation
	var err error

	provider.solanaClient = solclient.NewClient(config.SourceNodeUrl)
	provider.solanaCtx = context.Background()

	provider.ethClient, err = ethclient.DialContext(context.Background(), config.DestinationNodeUrl)
	if err != nil {
		return err
	}
	provider.ibportContract, err = ibport.NewIBPort(common.HexToAddress(config.LUPortAddress), provider.ethClient)
	if err != nil {
		return err
	}
	provider.LUPortDataAccount = config.LUPortAddress

	return nil
}

func (provider *SolanaToEVMExtractionBridge) LUPortState() (*PortContractState, error) {
	return GetSolanaPortContractState(
		provider.solanaClient,
		provider.solanaCtx,
		provider.config.LUPortAddress,
	)
}

func (provider *SolanaToEVMExtractionBridge) rqBytesToBigInt(rqId [32]byte) *big.Int {
	id := big.NewInt(0)
	id.SetBytes(rqId[:])
	return id
}


func (provider *SolanaToEVMExtractionBridge) ExtractDirectTransferRequest(context.Context) (*extractors.Data, error) {
	panic("not impl")
}

func (provider *SolanaToEVMExtractionBridge) ExtractReverseTransferRequest(context.Context) (*extractors.Data, error) {
	panic("not impl")
}

// ExtractDirectTransferRequest(context.Context) (*extractors.Data, error)
// ExtractReverseTransferRequest(context.Context) (*extractors.Data, error)