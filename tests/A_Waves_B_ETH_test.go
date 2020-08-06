package tests

import (
	waves "github.com/Gravity-Hub-Org/susy-data-extractor/v2/swagger-models/models"
	_ "go/types"
	"testing"
)

type TransactionKeyType = string
type TransactionEnumType = int


const (
	TransactionTypeKey TransactionKeyType = "type"
)

const (
	TransferTransactionType TransactionEnumType = 4
)

const (
	SwapRecipient = "william"
	SwapSender = "mary"
	SwapAmount int64 = 10 * 1e8
)

const (
	LastHeight = 4
	MaxBlocksPerRequest = 99
)

type ExtractionReadyState struct {
	ProcessedRequestIDList []string
	BlockForRequest int64
	BlockInterestRange []int64
	ReceivingAddress, SenderAddress string
}

func (state *ExtractionReadyState) HasBlockInterest (height int64) bool {
	begin, end := state.BlockInterestRange[0], state.BlockInterestRange[1]

	return height >= begin && height <= end
}

type SwapRequest struct {
	Sender, Recipient string
	Amount int64
	Currency string
}

func MockupBlock() []*waves.Block {
	heights := []int32 { 1, 2, 3 }

	return []*waves.Block {
		&waves.Block{
			Height: &heights[0],
			Transactions:     []waves.Transaction {
				map[string]interface{} {
					TransactionTypeKey: 1,
				},
				map[string]interface{} {
					TransactionTypeKey: 2,
				},
				map[string]interface{} {
					TransactionTypeKey: 3,
				},
				map[string]interface{} {
					TransactionTypeKey: TransferTransactionType,
					"recipient": "john",
				},
			},
		},
		&waves.Block{
			Height: &heights[1],
			Transactions:     []waves.Transaction {
				map[string]interface{} {
					TransactionTypeKey: 1,
				},
				map[string]interface{} {
					TransactionTypeKey: 2,
				},
				map[string]interface{} {
					TransactionTypeKey: 3,
				},
				map[string]interface{} {
					TransactionTypeKey: TransferTransactionType,
					"recipient": "john",
				},
			},
		},
		&waves.Block{
			Height: &heights[2],
			Transactions:     []waves.Transaction {
				map[string]interface{} {
					"type": 1,
				},
				map[string]interface{} {
					"type": 2,
				},
				map[string]interface{} {
					"type": 3,
				},
				map[string]interface{} {
					"type": TransferTransactionType,
					"recipient": SwapRecipient,
					"sender": SwapSender,
					"amount": SwapAmount,
				},
			},
		},
	}
}

func TestExtractionFromWavesToEth(t *testing.T) {

	swapDesiredState := &ExtractionReadyState{
		ProcessedRequestIDList: []string {},
		BlockForRequest:        0,
		BlockInterestRange:     []int64 { LastHeight - MaxBlocksPerRequest, LastHeight },
		ReceivingAddress:       SwapRecipient,
		SenderAddress:			SwapSender,
	}

	// Abstract for swap request
	// Абстракция над требованиями к свопу
	swapRequest := &SwapRequest{
		Sender:    swapDesiredState.SenderAddress,
		Recipient: swapDesiredState.ReceivingAddress,
		Amount:    SwapAmount,
		Currency:  "USDT",
	}

	// Latest blocks
	//
	// Последние 99 блоков

	t.Logf("Fetching latest %v blocks... \n", MaxBlocksPerRequest)
	latestBlocks := MockupBlock()

	var isSenderValid, isRecipientValid bool
	var matchedTx map[string]interface{}

	t.Logf("Starting to check the blocks")
	for _, currentBlock := range latestBlocks {

		// Check whether blocks conforms to interest
		// Проверяем на соот-е блока которое требует запрос на своп
		if !swapDesiredState.HasBlockInterest(int64(*currentBlock.Height)) { continue }

		// Then check TX list
		// Проверяем на транзакции

		t.Logf("Block #%v is in interest range. Checking for it transactions", *currentBlock.Height)


		for _, tx := range currentBlock.Transactions {

			// Check for Transfer Transaction
			// Проверяем на Transfer Transaction
			if tx[TransactionTypeKey] == nil ||
				tx[TransactionTypeKey] != TransferTransactionType ||
				tx["sender"] == nil || tx["recipient"] == nil{ continue }

			t.Logf("Block #%v. Handled transfer transaction." +
				"\n" +
				"Checking for swap request conformance...",
				*currentBlock.Height,
			)

			isSenderValid = tx["sender"] == swapRequest.Sender
			isRecipientValid = tx["recipient"] == swapRequest.Recipient

			if isSenderValid && isRecipientValid {
				matchedTx = tx
				break
			}
		}
		if isSenderValid && isRecipientValid { break }
	}

	if isSenderValid && isRecipientValid {
		t.Logf("Valid Sender: %v \n", matchedTx["sender"])
		t.Logf("Valid Recipient: %v \n", matchedTx["recipient"])
		t.Logf("Valid Amount: %v \n", matchedTx[TransactionTypeKey])
	} else {
		t.Errorf("Request TX not found")
	}

}