package tests

import (
	"context"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
	"math/big"
	"testing"
)


var currentExtractor *susy.SourceExtractor

var ctx context.Context


var amountTestCases []*AmountTestCase
type AmountTestCase struct {
	input int64
	expected *big.Int
}

func TestMain(t *testing.M) {
	amountTestCases = []*AmountTestCase {
		/**
			Waves: https://wavesexplorer.com/stagenet/tx/9MwvMvKDRBHZoZVaqvMWY38Qmsik7zeAjP1NJhxZ6sEg
			Ropsten: https://ropsten.etherscan.io/tx/0x85f3bbf31627f3881d374d934ea056ced906e7a96d361c2932bbcb35bebf6103
		 */
		&AmountTestCase{ input: 250000000, expected: big.NewInt(2500000000000000000) },
		/*
			Additional tests, base on the same checks
		*/
		&AmountTestCase{ input: 240000000, expected: big.NewInt(2400000000000000000) },
		&AmountTestCase{ input: 230005000, expected: big.NewInt(2300050000000000000) },
		&AmountTestCase{ input: 210000006, expected: big.NewInt(2100000060000000000) },
	}

	t.Run()
}

func TestExtractionWavesSourceLock(t *testing.T) {
	ctx = context.Background()

}

func errorHandler(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occured. %v", err)
		t.FailNow()
	}
}
