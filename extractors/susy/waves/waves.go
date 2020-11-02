package waves

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/mr-tron/base58"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
)

type SourceExtractor struct {
	options *susy.WavesEthereumBridgeOptions
}

func New(options *susy.WavesEthereumBridgeOptions) *SourceExtractor {
	extractor := &SourceExtractor { options: options }

	return extractor
}

func (ext *SourceExtractor) Info() *extractors.ExtractorInfo {
	return &extractors.ExtractorInfo{
		Tag:         "source-waves",
		Description: "Source WAVES",
	}
}

func (ext *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	options := ext.options

	states, _, err := options.WavesHelper.StateByAddress(options.LUContract, ctx)
	if err != nil {
		return nil, err
	}

	luState := susy.ParseState(states)

	rq, rqInt, _ := ext.pickRequestFromQueue(luState)

	if rq == "" || rqInt == nil {
		return nil, extractors.NotFoundErr
	}

	amount := luState.Requests()[rq].Amount
	receiver := luState.Requests()[rq].Receiver

	if !common.IsHexAddress(receiver) {
		options.Cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	if empty := make([]byte, 20, 20); bytes.Equal(receiverBytes, empty[:]) {
		options.Cache[rq] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	newAmount := ext.MapWavesAmount(amount)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	result := []byte{'m'}
	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiverBytes...)
	options.Cache[rq] = time.Now().Add(susy.MaxRqTimeout * time.Second)
	println(base64.StdEncoding.EncodeToString(result))
	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (ext *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}


func (ext *SourceExtractor) pickRequestFromQueue(luState *susy.LUWavesState) (susy.RequestId, *big.Int, error) {
	options := ext.options

	var rq susy.RequestId
	var rqInt *big.Int

	for target := luState.FirstRq; true; target = luState.Requests()[target].Next {
		if target == "" {
			break
		}
		if v, ok := options.Cache[target]; ok {
			if time.Now().After(v) {
				delete(options.Cache, target)
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
		status, err := options.IBContract.SwapStatus(nil, targetInt)
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
func (ext *SourceExtractor) MapWavesAmount(amount int64) *big.Int {
	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(susy.WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(susy.EthDecimals), nil)

	newAmount := bigIntAmount.Mul(bigIntAmount, ethDecimals).Div(bigIntAmount, wavesDecimals)

	return newAmount
}