package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
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

type EthereumExtractionProvider struct {}

// IB Port request state
const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	EthereumRequestStatusSuccess // is 3
	EthereumRequestStatusReturned
)


// LU port request state
const (
	EthereumRequestStatusCompleted = 2
)

const (
	MaxRqTimeout = 5 * 60 // 5 min
)

type EthereumToWavesExtractionBridge struct {
	config ConfigureCommand
	configured bool

	excludedRequests    []RequestID
	ethClient           *ethclient.Client
	wavesClient         *client.Client
	wavesHelper         helpers.ClientHelper

	luPortContract      *luport.LUPort
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
	provider.wavesClient, err = client.NewClient(client.Options{ BaseUrl: config.DestinationNodeUrl })
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

func (provider *EthereumToWavesExtractionBridge) pickRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (RequestID, *big.Int, *WavesRequestsState, error) {
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

	var rqIDInt *big.Int

	for rqIDInt = provider.rqBytesToBigInt(first);
		rqIDInt != nil;
		rqIDInt, _ = luState.NextRq(nil, rqIDInt) {


		wavesRequestID := RequestID(base58.Encode(rqIDInt.Bytes()))

		/**
		 * If lock request processed, but issue request didn't appear yet
		 */
		if ibRequest := ibState.Request(wavesRequestID); ibRequest != nil {
			// we ignore completed requests
			if ibRequest.Completed() {
				continue
			}

			break
		}

		break
	}

	if rqIDInt == nil {
		return "", nil, nil, extractors.NotFoundErr
	}

	return RequestID(base58.Encode(rqIDInt.Bytes())), rqIDInt, ibState, nil
}

func (provider *EthereumToWavesExtractionBridge) rqBytesToBigInt(rqID [32]byte) *big.Int {
	id := big.NewInt(0)
	id.SetBytes(rqID[:])
	return id
}

func (provider *EthereumToWavesExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	rqID, rqIDInt, ibState, err := provider.pickRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqID == "" || rqIDInt == nil {
		return nil, extractors.NotFoundErr
	}

	luPortRequest, err := provider.luPortContract.Requests(nil, rqIDInt)
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
	var action *big.Int
	var result []byte

	// no mint action passed, so we need process and instantiate one
	if ibState.Request(rqID) == nil {
		// if request is new/is not last unapproved we mint
		action = big.NewInt(int64(MintAction))

		result = action.FillBytes(resultAction[:])

		var bytesID [32]byte
		result = append(result, rqIDInt.FillBytes(bytesID[:])...)

		var bytesAmount [8]byte
		result = append(result, amount.FillBytes(bytesAmount[:])...)
		result = append(result, receiver[0:26]...)
	// if issue tx exists but is not completed we confirm it
	} else if !ibState.Request(rqID).Completed() {
		// if request is last unapproved & process we submit change status action
		action = big.NewInt(int64(ChangeStatusAction))

		result = action.FillBytes(resultAction[:])

		var bytesID [32]byte
		result = append(result, rqIDInt.FillBytes(bytesID[:])...)

		var bytesStatus [8]byte
		status := big.NewInt(int64(CompletedStatus))
		result = append(result, status.FillBytes(bytesStatus[:])...)
	}

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

	var unlockRqID RequestID
	var unlockRqStatus uint8
	var burnRq *Request

	id := big.NewInt(0)
	id.SetBytes(requestIds.First[:])

	for burnRq = ibState.Request(ibState.FirstRq); burnRq != nil; burnRq = ibState.Request(burnRq.Next) {
		targetInt := big.NewInt(0)
		bRq, err := base58.Decode(string(burnRq.RequestID))
		if err != nil {
			return nil, err
		}

		if burnRq.Receiver == "" {
			continue
		}

		targetInt.SetBytes(bRq)
		unlockRequest, err := provider.luPortContract.Requests(nil, targetInt)
		if err != nil {
			return nil, err
		}

		if unlockRequest.Status == EthereumRequestStatusCompleted {
			continue
		}

		if unlockRequest.Status == EthereumRequestStatusNew || unlockRequest.Status == EthereumRequestStatusNone {
			unlockRqID = burnRq.RequestID
			unlockRqStatus = unlockRequest.Status
			break
		}

		break
	}

	if unlockRqID == "" || burnRq == nil {
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

	rqID := burnRq.RequestID
	rqIDInt, err := rqID.ToBig()
	if err != nil {
		return nil, err
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	var result []byte
	if unlockRqStatus == uint8(EthereumRequestStatusNone) {
		result = []byte{'u'} // represents 'unlock' action
		result = append(result, rqIDInt.Bytes()[:]...)

		var bytesAmount [32]byte
		result = append(result, amount.FillBytes(bytesAmount[:])...)

		result = append(result, receiverBytes[0:20]...) // waves address is 20 bytes long

	} else if unlockRqStatus == uint8(EthereumRequestStatusNew) {
		result = []byte{'a'} // represents 'approve' action
		result = append(result, rqIDInt.Bytes()[:]...)
	} else {
		return nil, extractors.NotFoundErr
	}

	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}