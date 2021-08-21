package bridge

import (
	"context"
	"encoding/base64"

	solclient "github.com/portto/solana-go-sdk/client"
)


func GetSolanaPortContractState(solanaClient *solclient.Client, ctx context.Context, dataAccount string) (*PortContractState, error) {
	ibportStateResult, err := solanaClient.GetAccountInfo(ctx, dataAccount, solclient.GetAccountInfoConfig{
		Encoding: "base64",
	})
	if err != nil {
		return nil, err
	}

	ibportState := ibportStateResult.Data.([]interface{})[0].(string)

	stateDecoded, err := base64.StdEncoding.DecodeString(ibportState)
	if err != nil {
		return nil, err
	}

	return DecodePortState(stateDecoded), nil
}