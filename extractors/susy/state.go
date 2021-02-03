package susy

import (
	"strings"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/helpers"
)

const (
	FirstRqKey       = "first_rq"
	LastRqKey        = "last_rq"
	NebulaAddressKey = "nebula_address"

	NewStatus       Status = 1
	CompletedStatus Status = 2

	Approve Action = 1
	Unlock  Action = 1
)

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

type WavesRequestState struct {
	requests      map[RequestId]*Request
	FirstRq       RequestId
	LastRq        RequestId
	NebulaAddress string
}

func ParseState(states []helpers.State) *WavesRequestState {
	luState := &WavesRequestState{
		requests: make(map[RequestId]*Request),
	}
	for _, record := range states {
		switch record.Key {
		case FirstRqKey:
			luState.FirstRq = RequestId(record.Value.(string))
		case LastRqKey:
			luState.LastRq = RequestId(record.Value.(string))
		case NebulaAddressKey:
			luState.NebulaAddress = record.Value.(string)
		default:
			partsOfKey := strings.Split(record.Key, "_")
			if len(partsOfKey) != 3 {
				continue
			}
			requestID := RequestId(partsOfKey[2])
			if requestID == "" {
				continue
			}
			staticPart := partsOfKey[0] + "_" + partsOfKey[1]

			hashmapRecord, ok := luState.requests[requestID]
			if !ok {
				hashmapRecord = &Request{
					RequestID: requestID,
				}
			}

			switch staticPart {
			case "next_rq":
				hashmapRecord.Next = RequestId(record.Value.(string))
			case "prev_rq":
				hashmapRecord.Prev = RequestId(record.Value.(string))
			case "rq_receiver":
				hashmapRecord.Receiver = record.Value.(string)
			case "rq_amount":
				hashmapRecord.Amount = int64(record.Value.(float64))
			case "rq_status":
				hashmapRecord.Status = int(record.Value.(float64))
			case "rq_type":
				hashmapRecord.Status = int(record.Value.(float64))
			}

			luState.requests[requestID] = hashmapRecord
		}
	}
	return luState
}

func (state *WavesRequestState) Request(id RequestId) *Request {
	return state.requests[id]
}
