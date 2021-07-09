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


func (provider *EthereumToSolanaExtractionBridge) pickRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (SwapID, *big.Int, error) {
	first := *byte32(firstRqId)

	if luState == nil || first == [32]byte{} {
		return *new(SwapID), nil, fmt.Errorf("invalid input")
	}

	ibState, err := provider.IBPortState()
	if err != nil {
		return *new(SwapID), nil, err
	}

	var resultRqIdInt *big.Int
	var rqIdInt *big.Int

	// fmt.Printf("ibState: %+v \n", ibState)

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt.Int64() != 0; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {
		/**
		 * Due to a fact, that current gateway implementation
		 * on smart contracts (ports) does not have additional
		 * confirmation tx, we should check just for the existence of the swap with that id
		 */
		//if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) == CompletedStatus {
		// if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) != CompletedStatus {
		// 	continue
		// }
		var requestIdFixed SwapID

		// validate target address
		luRequest, err := luState.Requests(nil, rqIdInt)
		if err != nil {
			// fmt.Printf("swap_id: %v; err on lu request fetch: %v \n", requestIdFixed, err)
			continue
		}

		// read requests on lock only (direct flow)
		if luRequest.Status != EthereumRequestStatusNew {
			continue
		}

		// fmt.Printf("lu: %+v \n", luRequest)
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 
		
		// fmt.Printf("requestIdFixed: %v \n", requestIdFixed)

		ibRequestStatus := ibState.SwapStatusDict[requestIdFixed]

		// fmt.Printf("RQ: %v; IB status: %v \n", requestIdFixed, ibRequestStatus)

		// if ibRequestStatus != nil {
		// 	fmt.Printf("IB status(unwrapped): %v \n", *ibRequestStatus)
		// } else {
		// 	fmt.Printf("LU request: %+v \n", luRequest)
		// 	fmt.Printf("IB status == nil: %v \n", requestIdFixed)
		// }

		if ibRequestStatus != nil && *ibRequestStatus == EthereumRequestStatusSuccess {
			continue
		}

		fmt.Printf("Solana Address: %v \n", base58.Encode(luRequest.ForeignAddress[0:32]))
		if !ValidateSolanaAddress(base58.Encode(luRequest.ForeignAddress[0:32])) {
			// fmt.Printf("swap_id: %v; solana address is invalid \n", requestIdFixed)
			continue
		}

		isTokenDataAccountPassed, tokenDataErr := ValidateSolanaTokenAccountOwnershipByTokenProgram(
			provider.solanaClient, 
			base58.Encode(luRequest.ForeignAddress[0:32]),
			provider.config.Meta,
		)
		if tokenDataErr != nil || !isTokenDataAccountPassed {
			// fmt.Printf("swap_id: %v; tokenDataErr: %v \n", requestIdFixed, tokenDataErr)
			continue
		}

		if luRequest.Amount.Uint64() == 0 {
			// fmt.Printf("swap_id: %v; amount is zero \n", requestIdFixed)
			continue
		}

		// fmt.Printf("rqID last: %v \n", rqIdInt.Bytes()[0:16])

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

func (provider *EthereumToSolanaExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	
	// // pick up unprocessed request
	// for swapId, requestStatus := range ibportContractState.SwapStatusDict {
	// 	if requestStatus == EthereumRequestStatusNew {
	// 		currentSwapId = swapId
	// 	}
	// }
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	
	if bytes.Equal(luRequestIds.First[:], make([]byte, 32)) {
		return nil, extractors.NotFoundErr
	}

	if err != nil {
		return nil, err
	}

	rqId, rqIdInt, err := provider.pickRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luRequest, err := provider.luPortContract.Requests(nil, rqIdInt)

	decimalsDiff := big.NewInt(provider.config.SourceDecimals - provider.config.DestinationDecimals)

	divideBy := big.NewInt(0).Exp(big.NewInt(10), decimalsDiff, nil)

	targetAmount := luRequest.Amount.Div(luRequest.Amount, divideBy).Uint64()

	solanaDecimals := big.NewInt(0).
		Exp(big.NewInt(10), big.NewInt(provider.config.DestinationDecimals), nil)

	targetAmountCasted := float64(targetAmount) / float64(solanaDecimals.Uint64())

	fmt.Printf("rqId[:]: %v \n", rqId[:])
	fmt.Printf("rqId(len): %v \n", len(rqId))

	fmt.Printf("luRequest.ForeignAddress: %v \n", luRequest.ForeignAddress)
	fmt.Printf("luRequest.ForeignAddress: %v \n", base58.Encode(luRequest.ForeignAddress[:]))

	fmt.Printf("targetAmount: %v \n", targetAmount)
	fmt.Printf("targetAmountCasted: %v \n", targetAmountCasted)

	var resultByteVector [64]byte
	copy(resultByteVector[:], solexecutor.BuildCrossChainMintByteVector(rqId[:], luRequest.ForeignAddress, targetAmountCasted))

	// fmt.Printf("result byte string: %v \n", provider.requestSerializer(resultByteVector[:]))
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
		
		// fmt.Printf("status: %v \n", status)
		// fmt.Printf("swapRequestsCount: %+v \n", burnRequest)

		var evmSwapID [32]byte
		copy(evmSwapID[:], swapID[:])

		luRequest, err := provider.luPortContract.Requests(nil, BytesToBigInt(evmSwapID[:]))

		// fmt.Printf("lu: %+v \n", luRequest)
		if err != nil {
			return nil, err
		}

		// fmt.Printf("lu.Status: %+v \n", luRequest.Status)

		if luRequest.Status != EthereumRequestStatusNone {
			// fmt.Printf("failed status check: %v \n", luRequest.Status
			continue
		}

		if bytes.Equal(burnRequest.OriginAddress[:], make([]byte, 20)) || bytes.Equal(burnRequest.ForeignAddress[:], make([]byte, 32)) {
			// fmt.Printf("failed comparison: %v; %v \n", burnRequest.OriginAddress[:], burnRequest.ForeignAddress[:])

			continue
		}

		if !ValidateEthereumBasedAddress(hexutil.Encode(luRequest.ForeignAddress[0:20])) {
			// fmt.Printf("invalid foreign address: %v \n", hexutil.Encode(luRequest.ForeignAddress[0:20]))
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
