package tests

import (
	"math/big"
	"testing"
	"context"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy"
)


var currentExtractor *susy.SourceExtractor

var ctx context.Context


var amountTestCases []*AmountTestCase
type AmountTestCase struct {
	input int64
	expected *big.Int
}

func NewAmountTestCase(input int64) *AmountTestCase {
	expected := big.NewInt(input + 1)

	return &AmountTestCase{
		input: input,
		expected: expected,
	}
}

func TestMain(t *testing.M) {
	amountTestCases = []*AmountTestCase {
		NewAmountTestCase(2.5 * susy.WavesDecimals),
	}

	t.Run()
}

func TestExtractionWavesSourceLock(t *testing.T) {
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

	wavesProvider := &susy.WavesExtractionProvider{ ExtractorDelegate:currentExtractor }

	for testCaseIndex, testCase := range amountTestCases {
		testCaseNumber := testCaseIndex + 1

		mappedAmount := wavesProvider.MapWavesAmount(testCase.input)

		if mappedAmount != testCase.expected {
			t.Errorf(
				"#%v Amount map did not succeed. Input: %v; Output: %v; Expected: %v \n",
				testCaseNumber,
				testCase.input,
				mappedAmount,
				testCase.expected,
			)
			t.FailNow()
		}

		t.Logf("#%v Amount map succeed. Input: %v; Output: %v \n", testCaseNumber, testCase.input, mappedAmount)
	}
}

func errorHandler(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error occured. %v", err)
		t.FailNow()
	}
}
