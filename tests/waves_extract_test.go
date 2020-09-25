package tests

import (
	"testing"
	"context"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
)


var currentExtractor *susy.SourceExtractor

var ctx context.Context

func TestMain(t *testing.M) {

}

func errorHandler(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occured. %v", err)
		t.FailNow()
	}
}

func TestExtractionWavesEthereumLock(t *testing.T) {
	ctx = context.Background()

	extractor, err := susy.New(
		"https://nodes-stagenet.wavesnodes.com",
		"https://ropsten.infura.io/v3/663ad61d27254aac874ba7fc298e0956",
		"3MdQFC6chdxJ2WrxYV4ZidmutZdpzea1Kqp",
		"0x617832f23efE1896c7cAC6f67AF92cdcFFAE5F64",
		ctx,
		susy.WavesSourceLock,
	)

	errorHandler(t, err)

	currentExtractor = extractor

	extractedData, err := currentExtractor.Extract(ctx)

	errorHandler(t, err)

	extractedData.Value
}