package model

type Payment struct {
	ChainName string
	Amount string
	Decimals int
	Timestamp int
	SourceAddress, DestinationAddress string
}

type PaymentTransfer interface {
	LockFunds()
	UnlockFunds()

	BurnFunds()
	IssueFunds()
}

