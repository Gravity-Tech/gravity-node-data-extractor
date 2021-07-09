package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"

	// solcommand "github.com/Gravity-Tech/solanoid/commands"
	solexecutor "github.com/Gravity-Tech/solanoid/commands/executor"
	solclient "github.com/portto/solana-go-sdk/client"
)

type SolanaExtractionProvider struct{}


const (
	EVMLUPortRequestStatusNone = iota
	EVMLUPortRequestStatusNew
	EVMLUPortRequestStatusCompleted
	EVMLUPortRequestStatusConfirmed
)

type EthereumToSolanaExtractionBridge struct {
	config     ConfigureCommand
	configured   bool


	ethClient      *ethclient.Client

	luPortContract *luport.LUPort
	
	solanaClient   *solclient.Client
	solanaCtx       context.Context
	// solanaExecutor *solexecutor.GenericExecutor

	IBPortDataAccount string
}

func (provider *EthereumToSolanaExtractionBridge) Configure(config ConfigureCommand) error {
	if provider.configured {
		return fmt.Errorf("bridge is configured already")
	}

	provider.config = config

	// Node clients instantiation
	var err error
	provider.ethClient, err = ethclient.DialContext(context.Background(), config.SourceNodeUrl)
	if err != nil {
		return err
	}
	provider.luPortContract, err = luport.NewLUPort(common.HexToAddress(config.LUPortAddress), provider.ethClient)
	if err != nil {
		return err
	}

	provider.solanaClient = solclient.NewClient(config.DestinationNodeUrl)
	provider.solanaCtx = context.Background()

	provider.IBPortDataAccount = config.IBPortAddress

	// provider.solanaExecutor, err = solcommand.InitGenericExecutor(

	// )

	return nil
}

func (provider *EthereumToSolanaExtractionBridge) rqBytesToBigInt(rqId [32]byte) *big.Int {
	id := big.NewInt(0)
	id.SetBytes(rqId[:])
	return id
}


func (provider *EthereumToSolanaExtractionBridge) pickConfirmRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (SwapID, *big.Int, error) {
	first := *byte32(firstRqId)

	if luState == nil || first == [32]byte{} {
		return *new(SwapID), nil, fmt.Errorf("invalid input")
	}

	ibState, err := provider.IBPortState()
	if err != nil {
		return *new(SwapID), nil, err
	}

	var rqIdInt *big.Int

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {
		var requestIdFixed SwapID
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 

		fmt.Printf("ibState: %v \n", ibState)

		ibRequestStatus := ibState.SwapStatusDict[requestIdFixed]

		// checks for ibport
		if *ibRequestStatus != 3 { // success status
			continue
		}
		
		// validate target address
		luRequest, err := luState.Requests(nil, rqIdInt)
		if err != nil {
			continue
		}
		// checks for luport
		if luRequest.Status != EVMLUPortRequestStatusCompleted {
			continue
		}

		fmt.Printf("Solana Address: %v \n", base58.Encode(luRequest.ForeignAddress[0:32]))
		if !ValidateSolanaAddress(base58.Encode(luRequest.ForeignAddress[0:32])) {
			continue
		}

		isTokenDataAccountPassed, tokenDataErr := ValidateSolanaTokenAccountOwnershipByTokenProgram(
			provider.solanaClient, 
			base58.Encode(luRequest.ForeignAddress[0:32]),
			provider.config.Meta,
		)
		if tokenDataErr != nil || !isTokenDataAccountPassed {
			continue
		}

		if luRequest.Amount.Uint64() == 0 {
			continue
		}

		fmt.Printf("rqID last: %v \n", rqIdInt.Bytes()[0:16])

		break
	}

	fmt.Printf("rqID on input: %v \n", rqIdInt.Bytes()[0:16])

	if rqIdInt == nil {
		return *new(SwapID), nil, extractors.NotFoundErr
	}

	var swapID SwapID
	copy(swapID[:], rqIdInt.Bytes()[0:16])

	return swapID, rqIdInt, nil
}

func (provider *EthereumToSolanaExtractionBridge) pickMintRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (SwapID, *big.Int, error) {
	first := *byte32(firstRqId)

	if luState == nil || first == [32]byte{} {
		return *new(SwapID), nil, fmt.Errorf("invalid input")
	}

	ibState, err := provider.IBPortState()
	if err != nil {
		return *new(SwapID), nil, err
	}

	var rqIdInt *big.Int

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {
		var requestIdFixed SwapID
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 

		fmt.Printf("ibState: %v \n", ibState)

		ibRequestStatus := ibState.SwapStatusDict[requestIdFixed]
		
		fmt.Printf("status: %v \n", ibRequestStatus)
		if ibRequestStatus != nil && *ibRequestStatus == EthereumRequestStatusSuccess {
			continue
		}

		// validate target address
		luRequest, err := luState.Requests(nil, rqIdInt)
		if err != nil {
			continue
		}

		fmt.Printf("Solana Address: %v \n", base58.Encode(luRequest.ForeignAddress[0:32]))
		if !ValidateSolanaAddress(base58.Encode(luRequest.ForeignAddress[0:32])) {
			continue
		}

		isTokenDataAccountPassed, tokenDataErr := ValidateSolanaTokenAccountOwnershipByTokenProgram(
			provider.solanaClient, 
			base58.Encode(luRequest.ForeignAddress[0:32]),
			provider.config.Meta,
		)

		if tokenDataErr != nil || !isTokenDataAccountPassed {
			continue
		}

		if luRequest.Amount.Uint64() == 0 {
			continue
		}

		fmt.Printf("rqID last: %v \n", rqIdInt.Bytes()[0:16])

		break
	}

	fmt.Printf("rqID on input: %v \n", rqIdInt.Bytes()[0:16])

	if rqIdInt == nil {
		return *new(SwapID), nil, extractors.NotFoundErr
	}

	var swapID SwapID
	copy(swapID[:], rqIdInt.Bytes()[0:16])

	return swapID, rqIdInt, nil
}

func (provider *EthereumToSolanaExtractionBridge) IBPortState() (*IBPortContractState, error) {
	ibportStateResult, err := provider.solanaClient.GetAccountInfo(provider.solanaCtx, provider.IBPortDataAccount, solclient.GetAccountInfoConfig{
		Encoding: "base64",
	})
	if err != nil {
		return nil, err
	}

	ibportState := ibportStateResult.Data.([]interface{})[0].(string)

	stateDecoded, err := base64.StdEncoding.DecodeString(ibportState)
	if err != nil {
		return nil, err
	}

	return DecodeIBPortState(stateDecoded), nil
}

func (provider *EthereumToSolanaExtractionBridge) requestSerializer(array []byte) string {
	return base64.StdEncoding.EncodeToString(array)
}

type directExtractionResponse struct {
	rqId            *SwapID
	rqIdInt         *big.Int
	request          struct{ HomeAddress common.Address; Amount *big.Int; ForeignAddress [32]byte; Status uint8 }
	isConfirmationRq  bool
	targetAmountCasted float64
}
// func (*luport.LUPortCaller).Requests(opts *bind.CallOpts, arg0 *big.Int) (struct{HomeAddress common.Address; Amount *big.Int; ForeignAddress [32]byte; Status uint8}, error)

func (provider *EthereumToSolanaExtractionBridge) extractDirectConfirmRequest(ctx context.Context) (*directExtractionResponse, error) {
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	
	if bytes.Equal(luRequestIds.First[:], make([]byte, 32)) {
		return nil, extractors.NotFoundErr
	}

	if err != nil {
		return nil, err
	}

	rqId, rqIdInt, err := provider.pickConfirmRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luRequest, err := provider.luPortContract.Requests(nil, rqIdInt)
	if err != nil {
		return nil, err
	}

	decimalsDiff := big.NewInt(provider.config.SourceDecimals - provider.config.DestinationDecimals)

	divideBy := big.NewInt(0).Exp(big.NewInt(10), decimalsDiff, nil)

	targetAmount := luRequest.Amount.Div(luRequest.Amount, divideBy).Uint64()

	solanaDecimals := big.NewInt(0).
		Exp(big.NewInt(10), big.NewInt(provider.config.DestinationDecimals), nil)

	targetAmountCasted := float64(targetAmount) / float64(solanaDecimals.Uint64())

	_ = rqId

	return &directExtractionResponse{
		rqId: &rqId,
		rqIdInt: rqIdInt,
		request: luRequest,
		isConfirmationRq: true,
		targetAmountCasted: targetAmountCasted,
	}, nil
}

func (provider *EthereumToSolanaExtractionBridge) extractDirectMintRequest(ctx context.Context) (*directExtractionResponse, error) {
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	
	if bytes.Equal(luRequestIds.First[:], make([]byte, 32)) {
		return nil, extractors.NotFoundErr
	}

	if err != nil {
		return nil, err
	}

	rqId, rqIdInt, err := provider.pickMintRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luRequest, err := provider.luPortContract.Requests(nil, rqIdInt)
	if err != nil {
		return nil, err
	}
	decimalsDiff := big.NewInt(provider.config.SourceDecimals - provider.config.DestinationDecimals)

	divideBy := big.NewInt(0).Exp(big.NewInt(10), decimalsDiff, nil)

	targetAmount := luRequest.Amount.Div(luRequest.Amount, divideBy).Uint64()

	solanaDecimals := big.NewInt(0).
		Exp(big.NewInt(10), big.NewInt(provider.config.DestinationDecimals), nil)

	targetAmountCasted := float64(targetAmount) / float64(solanaDecimals.Uint64())

	_ = rqId

	return &directExtractionResponse{
		rqId: &rqId,
		rqIdInt: rqIdInt,
		request: luRequest,
		isConfirmationRq: false,
		targetAmountCasted: targetAmountCasted,
	}, nil
}

func (provider *EthereumToSolanaExtractionBridge) extractDirectRequest(ctx context.Context) (*directExtractionResponse, error) {
	confirmationRequest, err := provider.extractDirectConfirmRequest(ctx)
	if err != nil {
		return nil, err
	}
	if confirmationRequest == nil {
		return provider.extractDirectMintRequest(ctx)
	}

	return confirmationRequest, nil
}

func (provider *EthereumToSolanaExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	
	requestResponse, err := provider.extractDirectRequest(ctx)
	if err != nil {
		return nil, err
	}
	if requestResponse == nil {
		return nil, extractors.NotFoundErr
	}
	
	// confirmation flow
	if requestResponse.isConfirmationRq {
		var resultByteVector [64]byte
		var pos int

		resultByteVector[0] = 'c'
		pos += 1
		
		swapID := requestResponse.rqId
		copy(resultByteVector[pos:pos+16], swapID[:]) 

		pos += 16
		
		newStatus := 3
		resultByteVector[pos] = uint8(newStatus)

		return &extractors.Data{
			Type:  extractors.Base64,
			Value: provider.requestSerializer(resultByteVector[:]),
		}, err
	}

	rqId, luRequest, targetAmountCasted := requestResponse.rqId, requestResponse.request, requestResponse.targetAmountCasted

	fmt.Printf("rqId[:]: %v \n", requestResponse.rqId)
	fmt.Printf("luRequest.ForeignAddress: %v \n", luRequest.ForeignAddress)
	fmt.Printf("luRequest.ForeignAddress: %v \n", base58.Encode(luRequest.ForeignAddress[:]))
	fmt.Printf("targetAmount: %v \n", luRequest.Amount)
	fmt.Printf("targetAmountCasted: %v \n", targetAmountCasted)

	var resultByteVector [64]byte
	copy(resultByteVector[:], solexecutor.BuildCrossChainMintByteVector(rqId[:], luRequest.ForeignAddress, targetAmountCasted))

	fmt.Printf("result byte string: %v \n", provider.requestSerializer(resultByteVector[:]))
	fmt.Printf("result byte string(len): %v \n", len(resultByteVector))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(resultByteVector[:]),
	}, err
}

func (provider *EthereumToSolanaExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	ibState, err := provider.IBPortState()
	if err != nil {
		return nil, err
	}

	var rqSwapID SwapID
	var reverseRequest *IBPortContractUnwrapRequest

	for swapID, burnRequest := range ibState.RequestsDict {
		status := ibState.SwapStatusDict[swapID]

		if status == nil {
			continue
		}

		luRequest, err := provider.luPortContract.Requests(nil, swapID.AsBigInt())
		if err != nil {
			return nil, err
		}

		// checks on ibport side
		if *status != EVMLUPortRequestStatusNew {
			continue
		}

		// checks on luport side
		if luRequest.Status != EVMLUPortRequestStatusNone {
			continue
		}

		if bytes.Equal(luRequest.HomeAddress[:], make([]byte, 20)) || bytes.Equal(luRequest.ForeignAddress[:], make([]byte, 32)) {
			continue
		}

		if !ValidateEthereumBasedAddress(hexutil.Encode(luRequest.ForeignAddress[0:20])) {
			continue
		}

		if burnRequest.Amount == 0 {
			continue
		}

		reverseRequest = burnRequest
		rqSwapID = swapID
		break
	}

	if reverseRequest == nil {
		return nil, extractors.NotFoundErr
	}

	amount := MapAmount(int64(reverseRequest.Amount), provider.config.DestinationDecimals, provider.config.SourceDecimals)

	fmt.Printf("amount mapped: %v \n", amount)

	var rqIDBytes [32]byte
	copy(rqIDBytes[:], rqSwapID[:])

	var amountBytes [32]byte
	amount.FillBytes(amountBytes[:])

	var evmReceiver [20]byte
	copy(evmReceiver[:], reverseRequest.ForeignAddress[0:20])

	result := BuildForEVMByteArray('u', rqIDBytes, amountBytes, evmReceiver)

	fmt.Printf("result byte string: %v \n", provider.requestSerializer(result[:]))
	fmt.Printf("result byte string(len): %v \n", len(result))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(result[:]),
	}, err
}
