package solana

import "github.com/Gravity-Tech/solanoid/commands/ws"

type SubscriptionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  int    `json:"result"`
	ID      int    `json:"id"`
}



type IBPortMintToTransaction struct {
	BlockTime int `json:"blockTime"`
	Meta      struct {
		Err               interface{} `json:"err"`
		Fee               int         `json:"fee"`
		InnerInstructions []struct {
			Index        int `json:"index"`
			Instructions []struct {
				Accounts       []int  `json:"accounts"`
				Data           string `json:"data"`
				ProgramIDIndex int    `json:"programIdIndex"`
			} `json:"instructions"`
		} `json:"innerInstructions"`
		LogMessages       []string `json:"logMessages"`
		PostBalances      []int    `json:"postBalances"`
		PostTokenBalances []struct {
			AccountIndex  int    `json:"accountIndex"`
			Mint          string `json:"mint"`
			UITokenAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UIAmount       float64 `json:"uiAmount"`
				UIAmountString string  `json:"uiAmountString"`
			} `json:"uiTokenAmount"`
		} `json:"postTokenBalances"`
		PreBalances      []int `json:"preBalances"`
		PreTokenBalances []struct {
			AccountIndex  int    `json:"accountIndex"`
			Mint          string `json:"mint"`
			UITokenAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UIAmount       float64 `json:"uiAmount"`
				UIAmountString string  `json:"uiAmountString"`
			} `json:"uiTokenAmount"`
		} `json:"preTokenBalances"`
		Rewards []interface{} `json:"rewards"`
		Status  struct {
			Ok interface{} `json:"Ok"`
		} `json:"status"`
	} `json:"meta"`
	Slot        int `json:"slot"`
	Transaction struct {
		Message struct {
			AccountKeys []string `json:"accountKeys"`
			Header      struct {
				NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
				NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
				NumRequiredSignatures       int `json:"numRequiredSignatures"`
			} `json:"header"`
			Instructions []struct {
				Accounts       []int  `json:"accounts"`
				Data           string `json:"data"`
				ProgramIDIndex int    `json:"programIdIndex"`
			} `json:"instructions"`
			RecentBlockhash string `json:"recentBlockhash"`
		} `json:"message"`
		Signatures []string `json:"signatures"`
	} `json:"transaction"`
}

type IBPortMintToTransactionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result IBPortMintToTransaction `json:"result"`
	ID int `json:"id"`
}

func BuildAccountUnsubscribeRequest(subID int64) ws.RequestBody {
	return ws.RequestBody {
		Jsonrpc: "2.0",
		ID: 1,
		Method: "accountSubscribe",
		Params: []interface{} {
			subID,
		},
	}
}

func BuildLogsSubscribeRequest(watched string) ws.LogsSubscribeBody {
	return ws.LogsSubscribeBody {
		Jsonrpc: "2.0",
		ID: 1,
		Method: "logsSubscribe",
		Params: []ws.LogsSubscribeParam {
			{
				Mentions: []string {
					watched,
				},
			},
			{
				Commitment: "finalized",
			},
		},
	}
}

func IsMintToTx(response *IBPortMintToTransaction) bool {
	patternMint := "Program log: Instruction: MintTo"

	for _, logLine := range response.Meta.LogMessages {
		if logLine == patternMint {
			return true
		}
	}

	return false
}