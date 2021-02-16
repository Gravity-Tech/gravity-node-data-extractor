package bridge

import (
	"context"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/mr-tron/base58"
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
