package susy

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/mr-tron/base58"
)

type EthereumExtractionProvider struct {
	ExtractorDelegate *SourceExtractor
}

func (provider *EthereumExtractionProvider) Extract(ctx context.Context) (*extractors.Data, error) {
	e := provider.ExtractorDelegate

	states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

	requestIds, err := e.ibContract.RequestsQueue(nil)
	if err != nil {
		return nil, err
	}

	var rqId RequestId
	var intRqId *big.Int

	id := big.NewInt(0)
	id.SetBytes(requestIds.First[:])

	for {
		wavesRequestId := RequestId(base58.Encode(id.Bytes()))

		luPortRequest := luState.Request(wavesRequestId)

		// Must be no such request on lu port
		if luPortRequest != nil {
			id, err = e.ibContract.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		status, err := e.ibContract.SwapStatus(nil, id)
		if err != nil {
			fmt.Printf("Error get status rq: %s \n", err.Error())
			id, err = e.ibContract.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		if status != EthereumRequestStatusNew {
			id, err = e.ibContract.NextRq(nil, id)
			if err != nil {
				return nil, err
			}
			continue
		}

		// Check cache
		if v, ok := e.cache[wavesRequestId]; ok {
			if time.Now().After(v) {
				delete(e.cache, wavesRequestId)
			} else {
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

	rq, err := e.ibContract.UnwrapRequests(nil, intRqId)
	if err != nil {
		return nil, err
	}

	amount := rq.Amount
	receiver := rq.ForeignAddress

	wavesDecimals := big.NewInt(10)
	// 10^8 = 1e8
	wavesDecimals.Exp(wavesDecimals, big.NewInt(WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	// 10^18 = 1e18
	ethDecimals.Exp(ethDecimals, big.NewInt(EthDecimals), nil)

	// mappedX = x / 1e18 * 1e8
	//
	// more commonly:
	//
	// mappedX = x / sourceChainDecimals * destinationChainDecimals
	newAmount := amount.Div(amount, ethDecimals).Mul(amount, wavesDecimals)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	//
	// 2 - Unlock action
	//
	var resultAction [8]byte
	action := big.NewInt(int64(2))
	result := action.FillBytes(resultAction[:])

	var bytesId [32]byte
	result = append(result, intRqId.FillBytes(bytesId[:])...)

	var bytesAmount [32]byte
	result = append(result, intRqId.FillBytes(bytesAmount[:])...)
	result = append(result, receiver[:]...)

	e.cache[rqId] = time.Now().Add(MaxRqTimeout * time.Second)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}
