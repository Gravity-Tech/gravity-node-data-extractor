package bridge

import "math/big"

type crossChainTransfer struct {
	Amount    *big.Int // [8]byte
	Receiver  [32]byte // [32]byte
	RequestId [32]byte
	Action    [8]byte
}

type transferRequestSerializer interface {
	Serialize() ([]byte, error)
	Deserialize() (*crossChainTransfer, error)
}
