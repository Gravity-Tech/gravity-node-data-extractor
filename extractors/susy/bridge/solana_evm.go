package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	solexecutor "github.com/Gravity-Tech/solanoid/commands/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
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
	provider.ibportContract, err = ibport.NewIBPort(common.HexToAddress(config.IBPortAddress), provider.ethClient)
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
	luState, err := provider.LUPortState()
	if err != nil {
		return nil, err
	}

	var rqSwapID SwapID
	var reverseRequest *PortContractUnwrapRequest

	for swapID, lockRequest := range luState.RequestsDict {
		status := luState.SwapStatusDict[swapID]

		if status == nil {
			continue
		}

		var evmSwapID [32]byte
		copy(evmSwapID[:], swapID[:])

		ibRequestStatus, err := provider.ibportContract.SwapStatus(nil, BytesToBigInt(evmSwapID[:]))

		if err != nil {
			return nil, err
		}

		if ibRequestStatus != EthereumRequestStatusNone {
			continue
		}

		if bytes.Equal(lockRequest.OriginAddress[:], make([]byte, 20)) || bytes.Equal(lockRequest.ForeignAddress[:], make([]byte, 32)) {
			continue
		}

		if !ValidateEthereumBasedAddress(hexutil.Encode(lockRequest.ForeignAddress[0:20])) {
			continue
		}

		if lockRequest.Amount == 0 {
			continue
		}

		reverseRequest = lockRequest
		rqSwapID = swapID
		break
	}


	if reverseRequest == nil {
		return nil, extractors.NotFoundErr
	}

	amount := MapAmount(int64(reverseRequest.Amount), provider.config.SourceDecimals, provider.config.DestinationDecimals)

	var rqIDBytes [32]byte
	copy(rqIDBytes[:], rqSwapID[:])

	var amountBytes [32]byte
	amount.FillBytes(amountBytes[:])

	var evmReceiver [20]byte
	copy(evmReceiver[:], reverseRequest.ForeignAddress[0:20])

	fmt.Printf("rq swap id: %v \n", rqIDBytes[:])
	fmt.Printf("rq swap id: %v \n", len(rqIDBytes[:]))
	fmt.Printf("rq amountBytes: %v \n", amountBytes[:])
	fmt.Printf("rq amount: %v \n", reverseRequest.Amount)
	fmt.Printf("rq amount (big): %v \n", amount.Int64())
	fmt.Printf("rq reverseRequest.ForeignAddress: %v \n", reverseRequest.ForeignAddress[0:20])

	result := BuildForEVMByteArray('m', rqIDBytes, amountBytes, evmReceiver)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(result[:]),
	}, err
}

func (provider *SolanaToEVMExtractionBridge) requestSerializer(array []byte) string {
	return base64.StdEncoding.EncodeToString(array)
}

func (provider *SolanaToEVMExtractionBridge) ExtractReverseTransferRequest(context.Context) (*extractors.Data, error) {
	ibRequestIds, err := provider.ibportContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(ibRequestIds.First[:], make([]byte, 32)) {
		return nil, extractors.NotFoundErr
	}

	rqId, rqIdInt, err := provider.pickRequestFromQueue(provider.ibportContract, ibRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luRequest, err := provider.ibportContract.UnwrapRequests(nil, rqIdInt)

	decimalsDiff := big.NewInt(provider.config.SourceDecimals - provider.config.DestinationDecimals)

	divideBy := big.NewInt(0).Exp(big.NewInt(10), decimalsDiff, nil)

	targetAmount := luRequest.Amount.Div(luRequest.Amount, divideBy).Uint64()

	solanaDecimals := big.NewInt(0).
		Exp(big.NewInt(10), big.NewInt(provider.config.DestinationDecimals), nil)

	targetAmountCasted := float64(targetAmount) / float64(solanaDecimals.Uint64())

	var resultByteVector [64]byte
	copy(resultByteVector[:], solexecutor.BuildCrossChainMintByteVector(rqId[:], luRequest.ForeignAddress, targetAmountCasted))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(resultByteVector[:]),
	}, err
}


func (provider *SolanaToEVMExtractionBridge) pickRequestFromQueue(ibState *ibport.IBPort, firstRqId []byte) (SwapID, *big.Int, error) {
	first := *byte32(firstRqId)

	if ibState == nil || first == [32]byte{} {
		return *new(SwapID), nil, fmt.Errorf("invalid input")
	}

	luState, err := provider.LUPortState()
	if err != nil {
		return *new(SwapID), nil, err
	}

	var resultRqIdInt *big.Int
	var rqIdInt *big.Int

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt.Int64() != 0; rqIdInt, _ = ibState.NextRq(nil, rqIdInt) {
		/**
		 * Due to a fact, that current gateway implementation
		 * on smart contracts (ports) does not have additional
		 * confirmation tx, we should check just for the existence of the swap with that id
		 */
		var requestIdFixed SwapID

		// validate target address
		ibRequest, err := ibState.UnwrapRequests(nil, rqIdInt)
		ibRequestStatus, err := ibState.SwapStatus(nil, rqIdInt)
		if err != nil {
			continue
		}

		// read requests on lock only (direct flow)
		if ibRequestStatus != EthereumRequestStatusNew {
			continue
		}

		// fmt.Printf("lu: %+v \n", luRequest)
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 

		luRequestStatus := luState.SwapStatusDict[requestIdFixed]

		if luRequestStatus != nil && *luRequestStatus == EthereumRequestStatusSuccess {
			continue
		}

		if !ValidateSolanaAddress(base58.Encode(ibRequest.ForeignAddress[0:32])) {
			continue
		}

		isTokenDataAccountPassed, tokenDataErr := ValidateSolanaTokenAccountOwnershipByTokenProgram(
			provider.solanaClient, 
			base58.Encode(ibRequest.ForeignAddress[0:32]),
			provider.config.Meta,
		)
		if tokenDataErr != nil || !isTokenDataAccountPassed {
			// fmt.Printf("swap_id: %v; tokenDataErr: %v \n", requestIdFixed, tokenDataErr)
			continue
		}

		if ibRequest.Amount.Uint64() == 0 {
			continue
		}

		resultRqIdInt = rqIdInt
		break
	}

	if resultRqIdInt == nil || resultRqIdInt.Uint64() == 0 {
		return *new(SwapID), nil, extractors.NotFoundErr
	}

	var swapID SwapID
	copy(swapID[:], resultRqIdInt.Bytes()[0:16])

	return swapID, resultRqIdInt, nil
}
