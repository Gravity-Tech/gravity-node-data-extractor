package ethereum

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"
	"math/big"
	"time"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/mr-tron/base58"
)

type DestinationExtractor struct {
	options *susy.WavesEthereumBridgeOptions
	//cache       map[susy.RequestId]time.Time
	//ethClient   *ethclient.Client
	//wavesClient *client.Client
	//wavesHelper helpers.ClientHelper
	//luContract  string
	//ibContract  *ibport.IBPort
}

func New(options *susy.WavesEthereumBridgeOptions) *DestinationExtractor {
	extractor := &DestinationExtractor { options: options }

	return extractor
}

func (e *DestinationExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
	if err != nil {
		return nil, err
	}

	luState := susy.ParseState(states)

	requestIds, homeAddresses, foreignAddresses, amounts, statuses, err := e.ibContract.GetRequests(nil)

	if err != nil {
		return nil, err
	}

	queueLength := len(requestIds)

	if length := queueLength; (length != len(homeAddresses)) || (length != len(foreignAddresses)) || (length != len(amounts)) || (length != len(statuses)) {
		return nil, fmt.Errorf("invalid response")
	}

	var rq susy.RequestId
	var rqInt *big.Int
	var matchIndex int

	// All arrays have the same length
	for i := 0; i < queueLength; i++ {
		requestId := requestIds[i]
		ibRequestStatus := statuses[i]
		stringifiedRequestId := susy.RequestId(base58.Encode(requestId.Bytes()))

		luPortRequest := luState.Request(stringifiedRequestId)

		// Must be no such request on lu port
		if luPortRequest != nil {
			continue
		}

		if ibRequestStatus != susy.EthereumRequestStatusNew {
			continue
		}

		// Check cache
		if v, ok := e.cache[stringifiedRequestId]; ok {
			if time.Now().After(v) {
				delete(e.cache, stringifiedRequestId)
			} else {
				continue
			}
		}

		rq = stringifiedRequestId
		rqInt = requestId
		matchIndex = i
		break
	}

	if rq == "" || rqInt == nil {
		return nil, extractors.NotFoundErr
	}

	amount := amounts[matchIndex]
	receiver := foreignAddresses[matchIndex]

	bigIntAmount := amount

	wavesDecimals := big.NewInt(10)
	// 10^8 = 1e8
	wavesDecimals.Exp(wavesDecimals, big.NewInt(susy.WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	// 10^18 = 1e18
	ethDecimals.Exp(ethDecimals, big.NewInt(susy.EthDecimals), nil)

	// mappedX = x / 1e18 * 1e8
	//
	// more commonly:
	//
	// mappedX = x / sourceChainDecimals * destinationChainDecimals
	newAmount := bigIntAmount.Div(bigIntAmount, ethDecimals).Mul(bigIntAmount, wavesDecimals)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	//
	// 2 - Unlock action
	//
	var resultAction [8]byte
	action := big.NewInt(int64(2))
	result := action.FillBytes(resultAction[:])

	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiver[:]...)

	e.cache[rq] = time.Now().Add(susy.MaxRqTimeout * time.Second)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (e *DestinationExtractor) Info() *extractors.ExtractorInfo {
	return &extractors.ExtractorInfo{
		Tag:         "source-eth",
		Description: "Source Ethereum",
	}
}

func (e *DestinationExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}