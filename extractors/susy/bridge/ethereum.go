package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	_ "github.com/Gravity-Tech/gateway/abi/ethereum/ibport"
	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
	"math/big"
	"time"
	"unsafe"
)

type EthereumExtractionProvider struct {}

const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	SuccessEthereum // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 5 * 60 // 5 min
)

type EthereumToWavesExtractionBridge struct {
	config ConfigureCommand
	configured bool

	cache         map[RequestId]time.Time
	ethClient     *ethclient.Client
	wavesClient   *client.Client
	wavesHelper   helpers.ClientHelper

	luPortContract *luport.LUPort
}

func (provider *EthereumToWavesExtractionBridge) Configure(config ConfigureCommand) error {
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
	provider.wavesClient, err = client.NewClient(client.Options{ BaseUrl: config.SourceNodeUrl })
	if err != nil {
		return err
	}
	provider.luPortContract, err = luport.NewLUPort(common.HexToAddress(config.LUPortAddress), provider.ethClient)
	if err != nil {
		return err
	}

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

	for rqIdInt = provider.rqBytesToBigInt(first);
		rqIdInt != nil;
		rqIdInt, _ = luState.NextRq(nil, rqIdInt) {

		wavesRequestId := RequestId(base58.Encode(rqIdInt.Bytes()))

		if v, ok := provider.cache[wavesRequestId]; ok {
			if time.Now().After(v) {
				delete(provider.cache, wavesRequestId)
			} else {
				continue
			}
		}

		status := Status(ibState.Request(wavesRequestId).Status)

		if status == CompletedStatus {
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

	luPortRequest, err := provider.luPortContract.Requests(nil, rqIdInt)
	if err != nil {
		return nil, err
	}

	amount := luPortRequest.Amount
	receiver := luPortRequest.ForeignAddress
	receiverBytes := receiver[:]

	strRqId := RequestId(base58.Encode(rqIdInt.Bytes()))

	if empty := make([]byte, 20, 20); bytes.Equal(receiverBytes, empty[:]) {
		provider.cache[strRqId] = time.Now().Add(24 * time.Hour)
		return nil, extractors.NotFoundErr
	}

	sourceDecimals := big.NewInt(10)
	sourceDecimals.Exp(sourceDecimals, big.NewInt(provider.config.SourceDecimals), nil)

	destinationDecimals := big.NewInt(10)
	destinationDecimals.Exp(destinationDecimals, big.NewInt(provider.config.DestinationDecimals), nil)

	amount = amount.Mul(amount, accuracy).
		Div(amount, destinationDecimals).
		Mul(amount, sourceDecimals).
		Div(amount, accuracy)

	var resultAction [8]byte
	// completed on waves side
	action := big.NewInt(int64(MintAction))
	result := action.FillBytes(resultAction[:])

	var bytesId [32]byte
	result = append(result, rqIdInt.FillBytes(bytesId[:])...)

	var bytesAmount [8]byte
	result = append(result, amount.FillBytes(bytesAmount[:])...)
	result = append(result, receiver[:]...)

	provider.cache[rqId] = time.Now().Add(MaxRqTimeout * time.Second)

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: base64.StdEncoding.EncodeToString(result),
	}, err
}

func (provider *EthereumToWavesExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	return nil, nil
}