package bridge

import (
	"context"
	"runtime/debug"
	"testing"
)



func ValidateError(t *testing.T, err error) {
	if err != nil {
		t.Logf("Error: %v \n", err)
		debug.PrintStack()
		t.FailNow()
	}
} 

func TestEVMToSolanaExtractionRequest(t *testing.T) {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bridge := new(EthereumToSolanaExtractionBridge)
	err = bridge.Configure(ConfigureCommand{
		LUPortAddress: "0x10a785aa24d8540C583Ad99Bc82E5d7aF61b5806",
		IBPortAddress: "7quasTYVeRvbDQxWSbDRrQygPnQrWJBZkdY2EQ3hVF7h",
		SourceDecimals: 18,
		DestinationDecimals: 8,
		SourceNodeUrl: "https://data-seed-prebsc-1-s1.binance.org:8545",
		DestinationNodeUrl: "https://api.devnet.solana.com",
	})

	ValidateError(t, err)

	_, err = bridge.ExtractDirectTransferRequest(ctx)
	ValidateError(t, err)

	
}