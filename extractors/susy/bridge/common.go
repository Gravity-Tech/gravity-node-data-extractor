package bridge

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	wavescrypto "github.com/wavesplatform/go-lib-crypto"

	solclient "github.com/portto/solana-go-sdk/client"
	solcommon "github.com/portto/solana-go-sdk/common"
	soltoken "github.com/portto/solana-go-sdk/tokenprog"

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

	Meta map[string]string
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

func ValidateSolanaTokenAccountOwnershipByTokenProgram(client *solclient.Client, tokenAccount string, metaData map[string]string) (bool, error) {
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

	tokenState := stateResult.Data.([]interface{})[0].(string)
	tokenStateDecoded, err := base64.StdEncoding.DecodeString(tokenState)
	if err != nil {
		return false, err
	}

	tokenAccountState, err := soltoken.TokenAccountFromData(tokenStateDecoded)
	if err != nil {
		return false, err
	}

	tokenMint := metaData["token_mint"]
	fmt.Printf("tokenMint: %+v \n", tokenMint)

	if tokenAccountState.Mint.ToBase58() != tokenMint {
		return false, nil
	}
	
	fmt.Printf("tokenAccount: %+v \n", tokenAccount)

	return true, nil
}

func ValidateSolanaAddress(address string) bool {
	_, err := base58.Decode(address)
	if err != nil {
		return false
	}
	return true
}