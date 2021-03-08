package bridge

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"unsafe"

	_ "github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
)

type EthereumExtractionProvider struct{}

const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	EthereumRequestStatusSuccess // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 5 * 60 // 5 min
)

type EthereumToWavesExtractionBridge struct {
	config     ConfigureCommand
	configured bool

	//cache         map[RequestId]time.Time
	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper

	luPortContract *luport.LUPort
}

func (provider *EthereumToWavesExtractionBridge) Configure(config ConfigureCommand) error {
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
	provider.wavesClient, err = client.NewClient(client.Options{BaseUrl: config.DestinationNodeUrl})
	if err != nil {
		return err
	}
	provider.luPortContract, err = luport.NewLUPort(common.HexToAddress(config.LUPortAddress), provider.ethClient)
	if err != nil {
		return err
	}

	provider.wavesHelper = helpers.NewClientHelper(provider.wavesClient)

	provider.configured = true

	return nil
}

func byte32(s []byte) (a *[32]byte) {
	if len(a) <= len(s) {
		a = (*[len(a)]byte)(unsafe.Pointer(&s[0]))
	}
	return a
}

func (provider *EthereumToWavesExtractionBridge) pickRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (RequestId, *big.Int, error) {
	first := *byte32(firstRqId)

	if luState == nil || first == [32]byte{} {
		return "", nil, fmt.Errorf("invalid input")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ibStates, _, err := provider.wavesHelper.StateByAddress(provider.config.IBPortAddress, ctx)
	if err != nil {
		return "", nil, err
	}

	ibState := ParseState(ibStates)

	var rqIdInt *big.Int

	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {

		wavesRequestId := RequestId(base58.Encode(rqIdInt.Bytes()))
		//
		//// temp hardcode for testnet
		//if wavesRequestId == "2" {
		//	continue
		//}

		/**
		 * Due to a fact, that current gateway implementation
		 * on smart contracts (ports) does not have additional
		 * confirmation tx, we should check just for the existence of the swap with that id
		 */
		//if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) == CompletedStatus {
		if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) != CompletedStatus {
			continue
		}

		// validate waves target address
		luRequest, err := luState.Requests(nil, rqIdInt)
		if err != nil {
			continue
		}
		targetAddress := common.BytesToAddress(luRequest.ForeignAddress[:])
		if !ValidateWavesAddress(targetAddress.String(), 'W') {
			continue
		}

		break
	}

	if rqIdInt == nil {
		return "", nil, extractors.NotFoundErr
	}

	return RequestId(base58.Encode(rqIdInt.Bytes())), rqIdInt, nil
}

func (provider *EthereumToWavesExtractionBridge) rqBytesToBigInt(rqId [32]byte) *big.Int {
	id := big.NewInt(0)
	id.SetBytes(rqId[:])
	return id
}

func (provider *EthereumToWavesExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	rqId, rqIdInt, err := provider.pickRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqId == "" || rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luPortRequest, err := provider.luPortContract.Requests(nil, rqIdInt)
	if err != nil {
		return nil, err
	}

	amount := luPortRequest.Amount
	receiver := luPortRequest.ForeignAddress

	sourceDecimals := big.NewInt(10)
	sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	destinationDecimals := big.NewInt(10)
	destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	amount = amount.
		Mul(amount, destinationDecimals).
		Div(amount, sourceDecimals)

	var resultAction [8]byte
	// completed on waves side
	action := big.NewInt(int64(MintAction))
	result := action.FillBytes(resultAction[:])

	var bytesId [32]byte
	result = append(result, rqIdInt.FillBytes(bytesId[:])...)

	var bytesAmount [8]byte
	result = append(result, amount.FillBytes(bytesAmount[:])...)
	result = append(result, receiver[0:26]...)

	println(amount.String())
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (provider *EthereumToWavesExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	states, _, err := provider.wavesHelper.StateByAddress(provider.config.IBPortAddress, ctx)
	if err != nil {
		return nil, err
	}

	ibState := ParseState(states)

	requestIds, err := provider.luPortContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	var unlockRqId RequestId
	var burnRq *Request

	id := big.NewInt(0)
	id.SetBytes(requestIds.First[:])

	for burnRq = ibState.Request(ibState.FirstRq); burnRq != nil; burnRq = ibState.Request(burnRq.Next) {
		targetInt := big.NewInt(0)
		bRq, err := base58.Decode(string(burnRq.RequestID))
		if err != nil {
			return nil, err
		}

		targetInt.SetBytes(bRq)
		unlockRequest, err := provider.luPortContract.Requests(nil, targetInt)
		if err != nil {
			return nil, err
		}

		// if request exists and is processed, skip it
		// we pick only non-existing unlockRequests on LU
		if unlockRequest.Status != EthereumRequestStatusNone {
			continue
		}

		if burnRq.Receiver == "" {
			continue
		}
		if !ValidateEthereumBasedAddress(burnRq.Receiver) {
			continue
		}

		unlockRqId = burnRq.RequestID
		break
	}

	if unlockRqId == "" {
		return nil, extractors.NotFoundErr
	}

	if burnRq == nil {
		return nil, extractors.NotFoundErr
	}

	amount := big.NewInt(burnRq.Amount)
	receiver := burnRq.Receiver

	if receiver == "" {
		return nil, fmt.Errorf("receiver cannot be an empty string")
	}

	sourceDecimals := big.NewInt(10)
	sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	destinationDecimals := big.NewInt(10)
	destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	amount = amount.
		Mul(amount, sourceDecimals).
		Div(amount, destinationDecimals)

	rqId := burnRq.RequestID
	rqIdInt, err := rqId.ToBig()
	if err != nil {
		return nil, err
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	fmt.Printf("RQ ID: %v; AMOUNT: %v; RECEIVER: %v\n", burnRq.RequestID, amount.Int64(), receiver)

	result := []byte{'u'} // means 'unlock'
	result = append(result, rqIdInt.Bytes()[:]...)

	var bytesAmount [32]byte
	result = append(result, amount.FillBytes(bytesAmount[:])...)

	result = append(result, receiverBytes[0:20]...)
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}
