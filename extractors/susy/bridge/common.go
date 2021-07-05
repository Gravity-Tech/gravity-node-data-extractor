package bridge

import (
	"context"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	wavescrypto "github.com/wavesplatform/go-lib-crypto"

	solclient "github.com/portto/solana-go-sdk/client"
	solcommon "github.com/portto/solana-go-sdk/common"

	"math/big"
)

/**
 * Struct, that conforms to ChainExtractionBridge
 * must provide 2 methods as origin to destination chain bridge.
 * Every separate origin is provided in separate file.
 *
 *
 * Bridge represents an interface for bidirectional access between chains.
 */
type ChainExtractionBridge interface {
	Configure(ConfigureCommand) error
	ExtractDirectTransferRequest(context.Context) (*extractors.Data, error)
	ExtractReverseTransferRequest(context.Context) (*extractors.Data, error)
}

type ConfigureCommand struct {
	LUPortAddress, IBPortAddress        string
	SourceDecimals, DestinationDecimals int64

	SourceNodeUrl, DestinationNodeUrl   string

	Meta map[string]interface{}
}

type RequestId string

func (req RequestId) ToBig() (*big.Int, error) {
	targetInt := big.NewInt(0)
	bRq, err := base58.Decode(string(req))
	if err != nil {
		return nil, err
	}

	targetInt.SetBytes(bRq)
	return targetInt, nil
}

type Request struct {
	RequestID RequestId
	Next      RequestId
	Prev      RequestId
	Receiver  string
	Amount    int64
	Status    int
	Type      int
}

type Status int
type Action int
type RequestType int

const (
	NewStatus          Status = 1
	CompletedStatus    Status = 2

	ApproveAction      Action = 1
	UnlockAction       Action = 2
	MintAction         Action = 1
	ChangeStatusAction Action = 2

	IssueType          RequestType = 1
	BurnType           RequestType = 2
	LockType           RequestType = 1
	UnlockType         RequestType = 2
)


func ValidateEthereumBasedAddress(address string) bool {
	return common.IsHexAddress(address)
}

func ValidateWavesAddress(address string, chainId byte) bool {
	instance := wavescrypto.NewWavesCrypto()
	return instance.VerifyAddress(wavescrypto.Address(address), wavescrypto.WavesChainID(chainId))
}

func ValidateSolanaTokenAccountOwnershipByTokenProgram(client *solclient.Client, tokenAccount string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stateResult, err := client.GetAccountInfo(ctx, tokenAccount, solclient.GetAccountInfoConfig{
		Encoding: "base64",
	})
	if err != nil {
		return false, err
	}

	if stateResult.Owner != solcommon.TokenProgramID.ToBase58() {
		return false, nil
	}

	return true, nil
}

func ValidateSolanaAddress(address string) bool {
	_, err := base58.Decode(address)
	if err != nil {
		return false
	}
	return true
}