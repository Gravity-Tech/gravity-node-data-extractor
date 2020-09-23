package susy

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

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
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	SuccessEthereum // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 20

	WavesDecimals = 8
	EthDecimals   = 18
)

type ExtractImplementation int

const (
	WavesSourceLock ExtractImplementation = iota
	EthereumSourceBurn
)

type SourceExtractor struct {
	implementation ExtractImplementation

	cache       map[RequestId]time.Time
	ethClient   *ethclient.Client
	wavesClient *client.Client
	wavesHelper helpers.ClientHelper
	luContract  string
	ibContract  *ibport.IBPort
}

func New(sourceNodeUrl string, destinationNodeUrl string, luAddress string, ibAddress string, ctx context.Context, impl ExtractImplementation) (*SourceExtractor, error) {
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
		implementation: impl,
		cache:          make(map[RequestId]time.Time),
		ethClient:      ethClient,
		wavesClient:    wavesClient,
		wavesHelper:    helpers.NewClientHelper(wavesClient),
		ibContract:     destinationContract,
		luContract:     luAddress,
	}, nil
}

func (e *SourceExtractor) Info() *extractors.ExtractorInfo {
	switch e.implementation {
	case WavesSourceLock:
		return &extractors.ExtractorInfo{
			Tag:         "source-waves",
			Description: "Source WAVES",
		}
	case EthereumSourceBurn:
		return &extractors.ExtractorInfo{
			Tag:         "source-eth",
			Description: "Source Ethereum",
		}
	}

	return nil
}

func (e *SourceExtractor) Extract(ctx context.Context) (*extractors.Data, error) {
	switch e.implementation {
	case WavesSourceLock:
		return e.wavesSourceLockExtract(ctx)
	case EthereumSourceBurn:
		return e.ethereumSourceBurnExtract(ctx)
	}

	return nil, fmt.Errorf("No impl available")
}

func (e *SourceExtractor) ethereumSourceBurnExtract(ctx context.Context) (*extractors.Data, error) {
	states, _, err := e.wavesHelper.StateByAddress(e.luContract, ctx)
	if err != nil {
		return nil, err
	}

	luState := ParseState(states)

	requestIds, homeAddresses, foreignAddresses, amounts, statuses, err := e.ibContract.GetRequests(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	})

	if err != nil {
		return nil, err
	}

	queueLength := len(requestIds)

	if length := queueLength; (length != len(homeAddresses)) || (length != len(foreignAddresses)) || (length != len(amounts)) || (length != len(statuses)) {
		return nil, fmt.Errorf("invalid response")
	}

	var rq RequestId
	var rqInt *big.Int
	var matchIndex int

	// All arrays have the same length
	for i := 0; i < queueLength; i++ {
		requestId := requestIds[i]
		ibRequestStatus := statuses[i]
		stringifiedRequestId := RequestId(base58.Encode(requestId.Bytes()))

		luPortRequest := luState.Request(stringifiedRequestId)

		// Must be no such request on lu port
		if luPortRequest != nil {
			continue
		}

		if ibRequestStatus != EthereumRequestStatusNew {
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
	wavesDecimals.Exp(wavesDecimals, big.NewInt(WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	// 10^18 = 1e18
	ethDecimals.Exp(ethDecimals, big.NewInt(EthDecimals), nil)

	// mappedX = x / 1e18 * 1e8
	//
	// more commonly:
	//
	// mappedX = x / sourceChainDecimals * destinationChainDecimals
	newAmount := bigIntAmount.Div(bigIntAmount, ethDecimals).Mul(bigIntAmount, wavesDecimals)

	var newAmountBytes [32]byte
	newAmount.FillBytes(newAmountBytes[:])

	// Won't change
	result := big.NewInt(int64(2)).Bytes()

	result = append(result, rqInt.Bytes()...)
	result = append(result, newAmountBytes[:]...)
	result = append(result, receiver[:]...)

	e.cache[rq] = time.Now().Add(MaxRqTimeout * time.Second)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (e *SourceExtractor) wavesSourceLockExtract(ctx context.Context) (*extractors.Data, error) {
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

		targetInt := big.NewInt(0)
		bRq, err := base58.Decode(string(target))
		if err != nil {
			return nil, err
		}

		targetInt.SetBytes(bRq)
		status, err := e.ibContract.SwapStatus(nil, targetInt)
		if err != nil {
			return nil, err
		}

		if status == SuccessEthereum {
			continue
		}

		rq = target
		rqInt = targetInt
		break
	}

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

	bigIntAmount := big.NewInt(amount)

	wavesDecimals := big.NewInt(10)
	wavesDecimals.Exp(wavesDecimals, big.NewInt(WavesDecimals), nil)

	ethDecimals := big.NewInt(10)
	ethDecimals.Exp(ethDecimals, big.NewInt(EthDecimals), nil)

	newAmount := bigIntAmount.Div(bigIntAmount, wavesDecimals).Mul(bigIntAmount, ethDecimals)
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

func (e *SourceExtractor) Aggregate(values []extractors.Data) (*extractors.Data, error) {
	//TODO most popular

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: values[0].Value,
	}, nil
}
