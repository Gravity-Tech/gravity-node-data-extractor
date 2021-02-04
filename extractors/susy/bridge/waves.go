package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
	"math/big"
	"strings"
	"time"
)

var (
	accuracy = big.NewInt(1).Exp(big.NewInt(10), big.NewInt(8), nil)
)

const (
	FirstRqKey       = "first_rq"
	LastRqKey        = "last_rq"
	NebulaAddressKey = "nebula_address"
)

type LUWavesState struct {
	requests      map[RequestId]*Request
	FirstRq       RequestId
	LastRq        RequestId
	NebulaAddress string
}

func ParseState(states []helpers.State) *LUWavesState {
	luState := &LUWavesState{
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

func (state *LUWavesState) Request(id RequestId) *Request {
	return state.requests[id]
}


type WavesToEthereumExtractionBridge struct {
	kind        extractors.ExtractorType

	cache         map[RequestId]time.Time
	ethClient     *ethclient.Client
	wavesClient   *client.Client
	wavesHelper   helpers.ClientHelper
	luPortAddress string
	ibPortAddress *ibport.IBPort

	sourceDecimals      int64
	destinationDecimals int64
}

func (provider *WavesToEthereumExtractionBridge) Configure(sourceNodeUrl string, destinationNodeUrl string,
	luAddress string, ibAddress string,
	sourceDecimals int64, destinationDecimals int64,
	ctx context.Context, impl extractors.ExtractorType) error {
	// Node clients instantiation
	ethClient, err := ethclient.DialContext(ctx, destinationNodeUrl)
	if err != nil {
		return err
	}
	wavesClient, err := client.NewClient(client.Options{BaseUrl: sourceNodeUrl})
	if err != nil {
		return err
	}
	//destinationContract, err := ibport.NewIBPort(common.HexToAddress(ibAddress), ethClient)
	//if err != nil {
	//	return nil, err
	//}
	//
	// extractor := &SourceExtractor{
	//	kind:                impl,
	//	cache:               make(map[RequestId]time.Time),
	//	ethClient:           ethClient,
	//	wavesClient:         wavesClient,
	//	wavesHelper:         helpers.NewClientHelper(wavesClient),
	//	ibPortAddress:       ibAddress,
	//	luPortAddress:       luAddress,
	//	sourceDecimals:      sourceDecimals,
	//	destinationDecimals: destinationDecimals,
	// }



	// destinationContract, err := ibport.NewIBPort(common.HexToAddress(ibAddress), ethClient)
	// if err != nil {
	//	 return nil, err
	// }

	destinationContract, err := ibport.NewIBPort(common.HexToAddress(ibAddress), ethClient)
	if err != nil {
		return err
	}

	provider.ibPortAddress = destinationContract
	provider.cache = make(map[RequestId]time.Time)
	provider.ethClient = ethClient
	provider.wavesClient = wavesClient
	provider.wavesHelper = helpers.NewClientHelper(wavesClient)
	provider.luPortAddress = luAddress
	provider.sourceDecimals = sourceDecimals
	provider.destinationDecimals = destinationDecimals

	return nil
}

func (provider *WavesToEthereumExtractionBridge) pickRequestFromQueue(luState *LUWavesState) (RequestId, *big.Int, error) {

	var rq RequestId
	var rqInt *big.Int

	for target := luState.FirstRq; true; target = luState.requests[target].Next {
		if target == "" {
			break
		}
		if v, ok := provider.cache[target]; ok {
			if time.Now().After(v) {
				delete(provider.cache, target)
			} else {
				continue
			}
		}

		targetInt := big.NewInt(0)
		bRq, err := base58.Decode(string(target))
		if err != nil {
			return "", nil, err
		}

		targetInt.SetBytes(bRq)
		status, err := provider.ibPortAddress.SwapStatus(nil, targetInt)
		if err != nil {
			return "", nil, err
		}

		if status == susy.SuccessEthereum {
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
func (provider *WavesToEthereumExtractionBridge) MapWavesAmount(amount int64) *big.Int {
	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(provider.sourceDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(provider.destinationDecimals), nil)

	newAmount := bigIntAmount.Mul(bigIntAmount, accuracy).
		Div(bigIntAmount, wavesDecimals).
		Mul(bigIntAmount, ethDecimals).
		Div(bigIntAmount, accuracy)

	return newAmount
}

//
// Decoupling is aimed for tests management
// It allows testing distinct functions
//
func (provider *WavesToEthereumExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {

	states, _, err := provider.wavesHelper.StateByAddress(provider.luPortAddress, ctx)
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

	if !common.IsHexAddress(receiver) {
		provider.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	if empty := make([]byte, 20, 20); bytes.Equal(receiverBytes, empty[:]) {
		provider.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	newAmount := provider.MapWavesAmount(amount)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	result := []byte{'m'}
	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiverBytes...)
	provider.cache[rq] = time.Now().Add(susy.MaxRqTimeout * time.Second)
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}


func (provider *WavesToEthereumExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	states, _, err := provider.wavesHelper.StateByAddress(provider.luPortAddress, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

	requestIds, err := e.ibPortAddress.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	var rqId RequestId
	var intRqId *big.Int

	id := big.NewInt(0)
	id.SetBytes(requestIds.First[:])

	for {
		if id.Cmp(big.NewInt(0)) == 0 {
			return nil, extractors.NotFoundErr
		}

		wavesRequestId := RequestId(base58.Encode(id.Bytes()))

		luPortRequest := luState.Request(wavesRequestId)

		// Must be no such request on lu port
		if luPortRequest != nil {
			id, err = provider.ibPortAddress.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		status, err := provider.ibPortAddress.SwapStatus(nil, id)
		if err != nil {
			fmt.Printf("Error get status rq: %s \n", err.Error())
			id, err = provider.ibPortAddress.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		if status != susy.EthereumRequestStatusNew {
			id, err = provider.ibPortAddress.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Check cache
		if v, ok := provider.cache[wavesRequestId]; ok {
			if time.Now().After(v) {
				delete(provider.cache, wavesRequestId)
			} else {
				id, err = provider.ibPortAddress.NextRq(nil, id)
				if err != nil {
					return nil, err
				}
				continue
			}
		}

		rqId = wavesRequestId
		intRqId = id
		break
	}

	if rqId == "" {
		return nil, extractors.NotFoundErr
	}

	rq, err := provider.ibPortAddress.UnwrapRequests(nil, intRqId)
	if err != nil {
		return nil, err
	}

	amount := rq.Amount
	receiver := rq.ForeignAddress

	sourceDecimals := big.NewInt(10)
	sourceDecimals.Exp(sourceDecimals, big.NewInt(e.sourceDecimals), nil)

	destinationDecimals := big.NewInt(10)
	destinationDecimals.Exp(destinationDecimals, big.NewInt(e.destinationDecimals), nil)

	amount = amount.Mul(amount, susy.accuracy).
		Div(amount, destinationDecimals).
		Mul(amount, sourceDecimals).
		Div(amount, susy.accuracy)

	//
	// 2 - Unlock action
	//
	var resultAction [8]byte
	action := big.NewInt(int64(2))
	result := action.FillBytes(resultAction[:])

	var bytesId [32]byte
	result = append(result, intRqId.FillBytes(bytesId[:])...)

	var bytesAmount [8]byte
	result = append(result, amount.FillBytes(bytesAmount[:])...)
	result = append(result, receiver[0:26]...)

	e.cache[rqId] = time.Now().Add(susy.MaxRqTimeout * time.Second)

	println(amount.String())
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}
