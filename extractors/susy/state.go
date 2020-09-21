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

type LUWavesState struct {
	requests      map[RequestId]*Request
	FirstRq       RequestId
	LastRq        RequestId
	NebulaAddress string
}

func ParseState(states []helpers.State) *LUWavesState {
	luState := &LUWavesState{}
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
			requestID := RequestId(partsOfKey[len(partsOfKey)-1])

			staticPart := strings.Join(partsOfKey[:len(partsOfKey)-1], "_")

			hashmapRecord, ok := luState.requests[requestID]

			if !ok {
				hashmapRecord = &Request{}
			}

			hashmapRecord.RequestID = requestID

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

func (state *LUWavesState) Request(id RequestId) *Request {
	return state.requests[id]
}
