package bridge

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
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

	var rqIdInt *big.Int

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {
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
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 
		
		fmt.Printf("ibState: %v \n", ibState)
		ibRequestStatus := ibState.SwapStatusDict[requestIdFixed]
		fmt.Printf("status: %v \n", ibRequestStatus)
		if ibRequestStatus != nil && *ibRequestStatus != EthereumRequestStatusSuccess {
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

		break
	}

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

func (provider *EthereumToSolanaExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	
	// // pick up unprocessed request
	// for swapId, requestStatus := range ibportContractState.SwapStatusDict {
	// 	if requestStatus == EthereumRequestStatusNew {
	// 		currentSwapId = swapId
	// 	}
	// }
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
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
	fmt.Printf("luRequest.ForeignAddress: %v \n", luRequest.ForeignAddress)
	fmt.Printf("luRequest.ForeignAddress: %v \n", base58.Encode(luRequest.ForeignAddress[:]))
	fmt.Printf("targetAmount: %v \n", targetAmount)
	fmt.Printf("targetAmountCasted: %v \n", targetAmountCasted)

	resultByteVector := solexecutor.BuildCrossChainMintByteVector(rqId[:], luRequest.ForeignAddress, targetAmountCasted)

	// println(amount.String())
	println(base64.StdEncoding.EncodeToString(resultByteVector))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(resultByteVector),
	}, err
}

func (provider *EthereumToSolanaExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	// states, _, err := provider.wavesHelper.StateByAddress(provider.config.IBPortAddress, ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// ibState := ParseState(states)

	// requestIds, err := provider.luPortContract.RequestsQueue(nil)
	// if err != nil {
	// 	return nil, err
	// }

	// var unlockRqId RequestId
	// var burnRq *Request

	// id := big.NewInt(0)
	// id.SetBytes(requestIds.First[:])

	// for burnRq = ibState.Request(ibState.FirstRq); burnRq != nil; burnRq = ibState.Request(burnRq.Next) {
	// 	targetInt := big.NewInt(0)
	// 	bRq, err := base58.Decode(string(burnRq.RequestID))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	targetInt.SetBytes(bRq)
	// 	unlockRequest, err := provider.luPortContract.Requests(nil, targetInt)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// if request exists and is processed, skip it
	// 	// we pick only non-existing unlockRequests on LU
	// 	if unlockRequest.Status != EthereumRequestStatusNone {
	// 		continue
	// 	}

	// 	if burnRq.Receiver == "" {
	// 		continue
	// 	}
	// 	if !ValidateEthereumBasedAddress(burnRq.Receiver) {
	// 		continue
	// 	}

	// 	unlockRqId = burnRq.RequestID
	// 	break
	// }

	// if unlockRqId == "" {
	// 	return nil, extractors.NotFoundErr
	// }

	// if burnRq == nil {
	// 	return nil, extractors.NotFoundErr
	// }

	// amount := big.NewInt(burnRq.Amount)
	// receiver := burnRq.Receiver

	// if receiver == "" {
	// 	return nil, fmt.Errorf("receiver cannot be an empty string")
	// }

	// sourceDecimals := big.NewInt(10)
	// sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	// destinationDecimals := big.NewInt(10)
	// destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	// amount = amount.
	// 	Mul(amount, sourceDecimals).
	// 	Div(amount, destinationDecimals)

	// rqId := burnRq.RequestID
	// rqIdInt, err := rqId.ToBig()
	// if err != nil {
	// 	return nil, err
	// }

	// receiverBytes, err := hexutil.Decode(receiver)
	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Printf("RQ ID: %v; AMOUNT: %v; RECEIVER: %v\n", burnRq.RequestID, amount.Int64(), receiver)

	// result := []byte{'u'} // means 'unlock'
	// result = append(result, rqIdInt.Bytes()[:]...)

	// var bytesAmount [32]byte
	// result = append(result, amount.FillBytes(bytesAmount[:])...)

	// result = append(result, receiverBytes[0:20]...)
	// println(base64.StdEncoding.EncodeToString(result))
	// return &extractors.Data {
	// 	Type:  extractors.Base64,
	// 	Value: base64.StdEncoding.EncodeToString(result),
	// }, err
	return nil, nil
}
