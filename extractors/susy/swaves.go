package susy

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/mr-tron/base58"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/contracts/ibport"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/wavesplatform/gowaves/pkg/client"
)

const (
	MaxRqTimeout = 20

	SuccessEthereum = 3
	WavesDecimals   = 8
	EthDecimals     = 18
)

type SourceExtractor struct {
	cache       map[RequestId]time.Time
	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper
	luContract  string
	ibContract  *ibport.IBPort
}

func New(sourceNodeUrl string, destinationNodeUrl string, luAddress string, ibAddress string, ctx context.Context) (*SourceExtractor, error) {
	ethClient, err := ethclient.DialContext(ctx, destinationNodeUrl)
	if err != nil {
		return nil, err
	}
	destinationContract, err := ibport.NewIBPort(common.HexToAddress(ibAddress), ethClient)
	if err != nil {
		return nil, err
	}
	wavesClient, err := client.NewClient(client.Options{BaseUrl: sourceNodeUrl})
	if err != nil {
		return nil, err
	}

	return &SourceExtractor{
		ethClient:   ethClient,
		wavesClient: wavesClient,
		wavesHelper: helpers.NewClientHelper(wavesClient),
		ibContract:  destinationContract,
		luContract:  luAddress,
	}, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	return &extractors.ExtractorInfo{
		Tag:         "source-waves",
		Description: "Source waves",
	}
}
func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

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

		var tagetInt *big.Int
		bRq, err := base58.Decode(string(target))
		if err != nil {
			return nil, err
		}

		tagetInt.SetBytes(bRq)
		status, err := e.ibContract.SwapStatus(nil, tagetInt)
		if err != nil {
			return nil, err
		}

		if status == SuccessEthereum {
			continue
		}

		rq = target
		rqInt = tagetInt
		break
	}

	if rq == "" || rqInt == nil {
		return nil, extractors.NotFoundErr
	}

	result := []byte{'m'}
	result = append(result, rqInt.Bytes()...)

	amount := luState.requests[rq].Amount
	receiver := luState.requests[rq].Receiver

	receiverBytes, err := hexutil.Decode(receiver)
	if err != nil {
		return nil, err
	}

	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(EthDecimals), nil)

	newAmount := bigIntAmount.Div(bigIntAmount, wavesDecimals).Mul(bigIntAmount, ethDecimals)
	newAmountBytes := newAmount.Bytes()
	if len(newAmountBytes) < 32 {
		empty := make([]byte, 32-len(newAmountBytes), 32-len(newAmountBytes))
		newAmountBytes = append(newAmountBytes, empty...)
	}

	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes...)
	result = append(result, receiverBytes...)

	e.cache[rq] = time.Now().Add(MaxRqTimeout * time.Second)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: "",
	}, err
}

func (e *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}
