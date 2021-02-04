package bridge

import (
	"context"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
)

/**
 * Struct, that conforms to ChainExtractionBridge
 * must provide 2 methods as origin to destination chain bridge.
 * Every separate origin is provided in separate file.
 */
type ChainExtractionBridge interface {
	ExtractDirectTransferRequest(context.Context) (*extractors.Data, error)
	ExtractReverseTransferRequest(context.Context) (*extractors.Data, error)
}

type RequestId string
type Status int
type Action int

type Request struct {
	RequestID RequestId
	Next      RequestId
	Prev      RequestId
	Receiver  string
	Amount    int64
	Status    int
	Type      int
}

const (
	NewStatus       Status = 1
	CompletedStatus Status = 2

	Approve         Action = 1
	Unlock          Action = 1
)