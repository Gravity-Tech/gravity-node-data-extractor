package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	//"fmt"
	"math/big"
	"strings"

	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
)

//
//var (
//	accuracy = big.NewInt(1).
//		Exp(big.NewInt(10), big.NewInt(18), nil)
//)

const (
	FirstRqKey       = "first_rq"
	LastRqKey        = "last_rq"
	NebulaAddressKey = "nebula_address"
)

// TODO Implement general queue iterator (Waves & ETH)
// bc too muchs queues
type WavesRequestsState struct {
	requests      map[RequestId]*Request
	FirstRq       RequestId
	LastRq        RequestId
	NebulaAddress string
}

func ParseState(states []helpers.State) *WavesRequestsState {
	luState := &WavesRequestsState{
		requests: make(map[RequestId]*Request),
	}
	for _, record := range states {
		switch record.Key {
		case FirstRqKey:
			luState.FirstRq = RequestId(record.Value.(string))
		case LastRqKey:
			luState.LastRq = RequestId(record.Value.(string))
		case NebulaAddressKey:
			luState.NebulaAddress = record.Value.(string)
		default:
			partsOfKey := strings.Split(record.Key, "_")
			if len(partsOfKey) != 3 {
				continue
			}
			requestID := RequestId(partsOfKey[2])
			if requestID == "" {
				continue
			}
			staticPart := partsOfKey[0] + "_" + partsOfKey[1]

			hashmapRecord, ok := luState.requests[requestID]
			if !ok {
				hashmapRecord = &Request{
					RequestID: requestID,
				}
			}

			switch staticPart {
			case "next_rq":
				hashmapRecord.Next = RequestId(record.Value.(string))
			case "prev_rq":
				hashmapRecord.Prev = RequestId(record.Value.(string))
			case "rq_receiver":
				hashmapRecord.Receiver = record.Value.(string)
			case "rq_amount":
				hashmapRecord.Amount = int64(record.Value.(float64))
			case "rq_status":
				hashmapRecord.Status = int(record.Value.(float64))
			case "rq_type":
				hashmapRecord.Status = int(record.Value.(float64))
			}

			luState.requests[requestID] = hashmapRecord
		}
	}
	return luState
}

func (state *WavesRequestsState) Request(id RequestId) *Request {
	return state.requests[id]
}

type WavesToEthereumExtractionBridge struct {
	config     ConfigureCommand
	configured bool

	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper

	ibPortContract *ibport.IBPort
}

func (provider *WavesToEthereumExtractionBridge) Configure(config ConfigureCommand) error {
	if provider.configured {
		return fmt.Errorf("bridge is configured already")
	}

	provider.config = config

	// Node clients instantiation
	var err error
	provider.ethClient, err = ethclient.DialContext(context.Background(), config.DestinationNodeUrl)
	if err != nil {
		return err
	}
	provider.wavesClient, err = client.NewClient(client.Options{BaseUrl: config.SourceNodeUrl})
	if err != nil {
		return err
	}
	provider.ibPortContract, err = ibport.NewIBPort(common.HexToAddress(config.IBPortAddress), provider.ethClient)
	if err != nil {
		return err
	}

	provider.wavesHelper = helpers.NewClientHelper(provider.wavesClient)

	provider.configured = true

	return nil
}

func (provider *WavesToEthereumExtractionBridge) pickRequestFromQueue(luState *WavesRequestsState) (RequestId, *big.Int, error) {
	var rq RequestId
	var rqInt *big.Int

	for target := luState.FirstRq; true; target = luState.requests[target].Next {
		if target == "" {
			break
		}

		targetInt := big.NewInt(0)
		bRq, err := base58.Decode(string(target))
		if err != nil {
			return "", nil, err
		}

		targetInt.SetBytes(bRq)
		status, err := provider.ibPortContract.SwapStatus(nil, targetInt)
		if err != nil {
			return "", nil, err
		}

		// if status exists:
		//  1. it has been immediately invoked, so we must skip this request
		if status == EthereumRequestStatusSuccess {
			continue
		}
		if !ValidateEthereumBasedAddress(luState.Request(target).Receiver) {
			continue
		}

		rq = target
		rqInt = targetInt
		break
	}

	return rq, rqInt, nil
}

//
// Map amount provided in waves attachment payment to ethereum
//
// Params:
// amount - asset value in wavelets
//
func MapAmount(amount int64, sourceDecimals, destinationDecimals int64) *big.Int {
	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(sourceDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(destinationDecimals), nil)

	newAmount := bigIntAmount.
		Mul(bigIntAmount, ethDecimals).
		Div(bigIntAmount, wavesDecimals)

	return newAmount
}

//
// Decoupling is aimed for tests management
// It allows testing distinct functions
//
func (provider *WavesToEthereumExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	states, _, err := provider.wavesHelper.StateByAddress(provider.config.LUPortAddress, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

	rq, rqInt, _ := provider.pickRequestFromQueue(luState)

	if rq == "" || rqInt == nil {
		return nil, extractors.NotFoundErr
	}

	amount := luState.requests[rq].Amount
	receiver := luState.requests[rq].Receiver

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	newAmount := MapAmount(
		amount,
		provider.config.SourceDecimals,
		provider.config.DestinationDecimals,
	)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	result := []byte{'m'}
	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiverBytes...)

	println(base64.StdEncoding.EncodeToString(result))
	
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (provider *WavesToEthereumExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	states, _, err := provider.wavesHelper.StateByAddress(provider.config.LUPortAddress, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

	requestIds, err := provider.ibPortContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	var rqId RequestId
	var intRqId *big.Int

	requestIDLength := 32

	id := big.NewInt(0)
	id.SetBytes(requestIds.First[:])

	for {
		if id.Cmp(big.NewInt(0)) == 0 {
			return nil, extractors.NotFoundErr
		}

		requestIDBuffer := bytes.NewBuffer(make([]byte, requestIDLength))
		requestIDBuffer.Write(id.Bytes()[:])

		requestID := requestIDBuffer.Bytes()[:]
		requestID = requestID[len(requestID) - requestIDLength:len(requestID)]

		wavesRequestId := RequestId(base58.Encode(requestID[:]))
		luPortRequest := luState.Request(wavesRequestId)

		onNext := func() error {
			id, err = provider.ibPortContract.NextRq(nil, id)
			return err
		}

		// Must be no such request on lu port
		if luPortRequest != nil {
			err = onNext()
			if err != nil {
				return nil, err
			}
			continue
		}

		status, err := provider.ibPortContract.SwapStatus(nil, id)
		if err != nil {
			fmt.Printf("Error get status rq: %s \n", err.Error())
			err = onNext()
                        if err != nil {
                                return nil, err
                        }
			continue
		}

		if status != EthereumRequestStatusNew {
			err = onNext()
                        if err != nil {
                                return nil, err
                        }
			continue
		}

		burnRequest, err := provider.ibPortContract.UnwrapRequests(nil, id)
		if err != nil {
			err = onNext()
                        if err != nil {
                                return nil, err
                        }
			continue
		}

		if !ValidateWavesAddress(base58.Encode(burnRequest.ForeignAddress[0:26]), 'W') {
			id, err = provider.ibPortContract.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		rqId = wavesRequestId
		intRqId = id
		break
	}

	if rqId == "" {
		return nil, extractors.NotFoundErr
	}

	rq, err := provider.ibPortContract.UnwrapRequests(nil, intRqId)
	if err != nil {
		return nil, err
	}

	amount := rq.Amount
	receiver := rq.ForeignAddress

	sourceDecimals := big.NewInt(10)
	sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	destinationDecimals := big.NewInt(10)
	destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	amount = amount.
		Mul(amount, sourceDecimals).
		Div(amount, destinationDecimals)

	//
	// 2 - Unlock action
	//
	var resultAction [8]byte
	// completed on waves side
	action := big.NewInt(int64(UnlockAction))
	result := action.FillBytes(resultAction[:])

	var bytesId [32]byte
	result = append(result, intRqId.FillBytes(bytesId[:])...)

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
