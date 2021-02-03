package susy

import (
	"bytes"
	"context"
	"encoding/base64"
	"math/big"
	"time"
	"github.com/ethereum/go-ethereum/common"
	"github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/mr-tron/base58"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
)

type WavesExtractionProvider struct {
	ExtractorDelegate *SourceExtractor
	IsDirect           bool
}

func (provider *WavesExtractionProvider) pickRequestFromQueue(luState *WavesRequestState) (RequestId, *big.Int, error) {
	e := provider.ExtractorDelegate

	var rq RequestId
	var rqInt *big.Int


	for target := luState.FirstRq; true; target = luState.requests[target].Next {
		if target == "" {
			break
		}
		if v, ok := e.cache[target]; ok {
			if time.Now().After(v) {
				delete(e.cache, target)
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
		status, err := ibContract.SwapStatus(nil, targetInt)
		if err != nil {
			return "", nil, err
		}

		if status == SuccessEthereum {
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
func (provider *WavesExtractionProvider) MapWavesAmount(amount int64) *big.Int {
	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(provider.ExtractorDelegate.sourceDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(provider.ExtractorDelegate.destinationDecimals), nil)

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
func (provider *WavesExtractionProvider) Extract(ctx context.Context) (*extractors.Data, error) {
	if provider.IsDirect {
		return provider.extractDirectRequest(ctx)
	} else {
		return provider.extractReverseRequest(ctx)
	}
}

// encapsulation
func (provider *WavesExtractionProvider) extractDirectRequest(ctx context.Context) (*extractors.Data, error) {
	// direct request - minting
	e := provider.ExtractorDelegate

	states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
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
		e.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	if empty := make([]byte, 20, 20); bytes.Equal(receiverBytes, empty[:]) {
		e.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	newAmount := provider.MapWavesAmount(amount)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	result := []byte{'m'}
	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiverBytes...)
	e.cache[rq] = time.Now().Add(MaxRqTimeout * time.Second)
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (provider *WavesExtractionProvider) extractReverseRequest(ctx context.Context) (*extractors.Data, error) {
	// reverse request - unlocking
	e := provider.ExtractorDelegate

	//states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
	//if err != nil {
	//	return nil, err
	//}
	ibContract, err := ibport.NewIBPort(common.HexToAddress(e.ibContract), e.ethClient)

	if rq == "" || rqInt == nil {
		return nil, extractors.NotFoundErr
	}

	amount := luState.requests[rq].Amount
	receiver := luState.requests[rq].Receiver

	if !common.IsHexAddress(receiver) {
		e.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	if empty := make([]byte, 20, 20); bytes.Equal(receiverBytes, empty[:]) {
		e.cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	newAmount := provider.MapWavesAmount(amount)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	result := []byte{'u'}
	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiverBytes...)
	e.cache[rq] = time.Now().Add(MaxRqTimeout * time.Second)
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}
