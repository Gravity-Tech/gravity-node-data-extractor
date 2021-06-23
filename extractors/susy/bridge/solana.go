package bridge

import (
	"context"
	"fmt"

	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

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
	solanaExecutor *solexecutor.GenericExecutor
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

	// provider.solanaExecutor, err = solexecutor.NewGe

	return nil
}

// func (provider *EthereumToSolanaExtractionBridge) pickRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (RequestId, *big.Int, error) {
// 	first := *byte32(firstRqId)

// 	if luState == nil || first == [32]byte{} {
// 		return "", nil, fmt.Errorf("invalid input")
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	ibStates, _, err := provider.wavesHelper.StateByAddress(provider.config.IBPortAddress, ctx)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	ibState := ParseState(ibStates)

// 	var rqIdInt *big.Int

// 	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {

// 		wavesRequestId := RequestId(base58.Encode(rqIdInt.Bytes()))
// 		//
// 		//// temp hardcode for testnet
// 		//if wavesRequestId == "2" {
// 		//	continue
// 		//}

// 		/**
// 		 * Due to a fact, that current gateway implementation
// 		 * on smart contracts (ports) does not have additional
// 		 * confirmation tx, we should check just for the existence of the swap with that id
// 		 */
// 		//if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) == CompletedStatus {
// 		if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) != CompletedStatus {
// 			continue
// 		}

// 		// validate waves target address
// 		luRequest, err := luState.Requests(nil, rqIdInt)
// 		if err != nil {
// 			continue
// 		}
// 		if !ValidateWavesAddress(base58.Encode(luRequest.ForeignAddress[0:26]), 'W') {
// 			continue
// 		}

// 		break
// 	}

// 	if rqIdInt == nil {
// 		return "", nil, extractors.NotFoundErr
// 	}

// 	return RequestId(base58.Encode(rqIdInt.Bytes())), rqIdInt, nil
// }


func (provider *EthereumToSolanaExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	// luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	// if err != nil {
	// 	return nil, err
	// }

	// rqId, rqIdInt, err := provider.pickRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	// if err != nil {
	// 	return nil, err
	// }
	// if rqId == "" || rqIdInt == nil {
	// 	return nil, extractors.NotFoundErr
	// }

	// luPortRequest, err := provider.luPortContract.Requests(nil, rqIdInt)
	// if err != nil {
	// 	return nil, err
	// }

	// amount := luPortRequest.Amount
	// receiver := luPortRequest.ForeignAddress

	// sourceDecimals := big.NewInt(10)
	// sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	// destinationDecimals := big.NewInt(10)
	// destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	// amount = amount.
	// 	Mul(amount, destinationDecimals).
	// 	Div(amount, sourceDecimals)

	// var resultAction [8]byte
	// // completed on waves side
	// action := big.NewInt(int64(MintAction))
	// result := action.FillBytes(resultAction[:])

	// var bytesId [32]byte
	// result = append(result, rqIdInt.FillBytes(bytesId[:])...)

	// var bytesAmount [8]byte
	// result = append(result, amount.FillBytes(bytesAmount[:])...)
	// result = append(result, receiver[0:26]...)

	// println(amount.String())
	// println(base64.StdEncoding.EncodeToString(result))
	// return &extractors.Data{
	// 	Type:  extractors.Base64,
	// 	Value: base64.StdEncoding.EncodeToString(result),
	// }, err
	return nil, nil
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
	// return &extractors.Data{
	// 	Type:  extractors.Base64,
	// 	Value: base64.StdEncoding.EncodeToString(result),
	// }, err
	return nil, nil
}
