package bridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Gravity-Tech/gateway/abi/ethereum/luport"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/solana"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"

	portdelegate "github.com/Gravity-Tech/gravity-node-data-extractor/v2/abi/portdelegate"

	// solcommand "github.com/Gravity-Tech/solanoid/commands"

	solexecutor "github.com/Gravity-Tech/solanoid/commands/executor"
	"github.com/Gravity-Tech/solanoid/commands/ws"
	solclient "github.com/portto/solana-go-sdk/client"
	solcommon "github.com/portto/solana-go-sdk/common"

	"github.com/gorilla/websocket"
	gorilla "github.com/gorilla/websocket"
)

type SolanaExtractionProvider struct {}


type SolanaMintWatchCommand struct {
	TheMint, MintReceiver solcommon.PublicKey
	Amount                float64
}

func (smwc *SolanaMintWatchCommand) Equals(b *SolanaMintWatchCommand) bool {
	if b == nil {
		return false
	}

	return smwc.Amount == b.Amount && 
		bytes.Equal(smwc.TheMint.Bytes(), b.TheMint.Bytes()) &&
		bytes.Equal(smwc.MintReceiver.Bytes(), b.MintReceiver.Bytes())
}

type SolanaMintWatcher struct {
	delegateCtx         context.Context
	delegateTransactor *bind.TransactOpts
	delegateClient     *ethclient.Client

	solanaClient   *solclient.Client
	solanaCtx       context.Context

	wsEndpoint          string
	// processingRequests  []*SolanaMintWatchCommand
	currentRequest      *SolanaMintWatchCommand
	// currentExecutionContext   context.Context

	currentConnection *gorilla.Conn
	currentSubID      *int

	portDelegateProps *PortDelegateProps
}


func (smw *SolanaMintWatcher) request(ctx context.Context, method string, params []interface{}, response interface{}) error {
	j, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return err
	}

	// post request
	req, err := http.NewRequestWithContext(ctx, "POST", "/", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// http client and send request
	httpclient := &http.Client{}
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}

	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", res.StatusCode)
	}
	return nil
}

func (smw *SolanaMintWatcher) subscribe(swapID [32]byte, luRequest *struct {
	HomeAddress    common.Address
	Amount         *big.Int
	ForeignAddress [32]byte
	Status         uint8
}) {
	watchRequest := solana.BuildLogsSubscribeRequest(smw.currentRequest.MintReceiver.ToBase58())

	watchRequestBytes, err := json.Marshal(&watchRequest)
	if err != nil {
		fmt.Println("Error on marshal")
		debug.PrintStack()
	}

	err = smw.currentConnection.WriteMessage(websocket.TextMessage, watchRequestBytes)
	if err != nil {
		fmt.Println("Error on write message")
		debug.PrintStack()
	}

	go func() {
		swapID, luRequest := swapID, luRequest

		for {
			_, message, err := smw.currentConnection.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			if smw.currentSubID == nil {
				var responseSubscribe solana.SubscriptionResponse
				err = json.Unmarshal(message, &responseSubscribe)
				if err != nil {
					fmt.Printf("Error on solana.SubscriptionResponse unpack: %v \n", err)
				}

				smw.currentSubID = &responseSubscribe.Result

				continue
			}

			var responseUnpacked ws.LogsSubscribeNotification
			err = json.Unmarshal(message, &responseUnpacked)
			if err != nil {
				fmt.Printf("Error on ws.LogsSubscribeNotification unpack: %v \n", err)
				continue
			}

			txID := responseUnpacked.Params.Result.Value.Signature
			if txID == "" {
				continue
			}

			ctx := context.Background()

			var ibportMintToResponse solana.IBPortMintToTransactionResponse
			err = smw.request(
				ctx, "getConfirmedTransaction", []interface{}{txID, "json"}, &ibportMintToResponse,
			)
			if err != nil {
				fmt.Printf("Error on GetConfirmedTransaction unpack: %v \n", err)
				continue
			}

			fmt.Printf("RESPONSE: %+v \n", ibportMintToResponse)

			if !solana.IsMintToTx(&ibportMintToResponse.Result) {
				fmt.Printf("not mint to tx: %v \n", txID)
				continue
			}

			txMeta := ibportMintToResponse.Result.Meta
			if len(txMeta.PostTokenBalances) == 0 || len(txMeta.PreTokenBalances) == 0 {
				fmt.Println("post/pre token balances are empty")
				continue
			}

			preBalanceState, postBalanceState := txMeta.PreTokenBalances[0], txMeta.PostTokenBalances[0]

			balanceDiff := postBalanceState.UITokenAmount.UIAmount - preBalanceState.UITokenAmount.UIAmount

			if balanceDiff == smw.currentRequest.Amount {
				fmt.Printf("mint handled!")

				// portDelegateProps := provider.portDelegateProps()
				portDelegate, err := portdelegate.NewPortMemorizer(
					common.HexToAddress(smw.portDelegateProps.DelegateAddress),
					smw.delegateClient,
				)
				if err != nil {
					fmt.Printf("Error: portDelegate instance failed")
				}

				var amountBytes [32]byte
				luRequest.Amount.FillBytes(amountBytes[:])
				fmt.Printf("amount bytes: %v \n", amountBytes)

				fmt.Printf("LU Request: %+v \n", luRequest)

				persistByteArray := PortDelegateClient.PersistByteIX(
					swapID,
					amountBytes,
					luRequest.ForeignAddress,
					luRequest.Status,
				)
				tx, err := portDelegate.Persist(smw.delegateTransactor, persistByteArray)
				if err != nil {
					fmt.Println("persist failed")
				}
				fmt.Printf("Tx: Storage - Persist: %v \n", tx.Hash().Hex())

				return
			}
		}
	}()
}

func (smw *SolanaMintWatcher) unsubscribe(command *SolanaMintWatchCommand) {

}

func (smw *SolanaMintWatcher) HandleExtraction(swapID [32]byte, luRequest *struct {
	HomeAddress    common.Address
	Amount         *big.Int
	ForeignAddress [32]byte
	Status         uint8
}, command *SolanaMintWatchCommand) {
	if command == nil { return }
	if smw.currentRequest != nil && smw.currentRequest.Equals(command) { return }

	smw.currentRequest = command

	if smw.currentConnection == nil {
		var err error
		smw.currentConnection, _, err = websocket.DefaultDialer.Dial(smw.wsEndpoint, nil)
		if err != nil {
			fmt.Printf("Error on subscription open")
			debug.PrintStack()
		}
	}

	go smw.subscribe(swapID, luRequest)
}

func (smw *SolanaMintWatcher) Busy() bool {
	// return len(smw.processingRequests) > 0
	return smw.currentRequest != nil
}

type EthereumToSolanaExtractionBridge struct {
	config     ConfigureCommand
	configured   bool

	ethClient          *ethclient.Client

	luPortContract *luport.LUPort
	
	solanaClient   *solclient.Client
	solanaCtx       context.Context

	IBPortDataAccount string

	mintWatcher    *SolanaMintWatcher
}

func (provider *EthereumToSolanaExtractionBridge) Configure(config ConfigureCommand) error {
	var err error

	if provider.configured {
		return fmt.Errorf("bridge is configured already")
	}

	provider.config = config

	delegateCallerStrPK := os.Getenv("DELEGATE_PRIV_KEY")
	if delegateCallerStrPK == "" {
		return fmt.Errorf("delegate caller is not set")
	}
	privateKey, err := ethcrypto.HexToECDSA(delegateCallerStrPK)
	if err != nil {
		return err
	}

	provider.mintWatcher = &SolanaMintWatcher{}

	provider.mintWatcher.delegateCtx = context.Background()
	
	provider.mintWatcher.delegateTransactor = bind.NewKeyedTransactor(privateKey)
	provider.mintWatcher.delegateTransactor.GasLimit = 150000 * 5
	provider.mintWatcher.delegateTransactor.Context = provider.mintWatcher.delegateCtx

	portDelegateProps := provider.portDelegateProps()

	provider.mintWatcher.portDelegateProps = &portDelegateProps

	provider.mintWatcher.delegateClient, err = ethclient.DialContext(context.Background(), portDelegateProps.NodeUrl)
	if err != nil {
		return err
	}

	// Node clients instantiation
	provider.ethClient, err = ethclient.DialContext(context.Background(), config.SourceNodeUrl)
	if err != nil {
		return err
	}
	provider.luPortContract, err = luport.NewLUPort(common.HexToAddress(config.LUPortAddress), provider.ethClient)
	if err != nil {
		return err
	}

	provider.solanaClient = solclient.NewClient(config.DestinationNodeUrl)
	provider.solanaCtx = context.Background()

	provider.mintWatcher.solanaClient = solclient.NewClient(config.DestinationNodeUrl)
	provider.mintWatcher.solanaCtx = context.Background()

	provider.IBPortDataAccount = config.IBPortAddress

	return nil
}

func (provider *EthereumToSolanaExtractionBridge) rqBytesToBigInt(rqId [32]byte) *big.Int {
	id := big.NewInt(0)
	id.SetBytes(rqId[:])
	return id
}


func (provider *EthereumToSolanaExtractionBridge) pickRequestFromQueue(luState *luport.LUPort, firstRqId []byte) (SwapID, *big.Int, error) {
	first := *byte32(firstRqId)

	if luState == nil || first == [32]byte{} {
		return *new(SwapID), nil, fmt.Errorf("invalid input")
	}

	ibState, err := provider.IBPortState()
	if err != nil {
		return *new(SwapID), nil, err
	}

	portDelegateProps := provider.portDelegateProps()
	portDelegate, err := portdelegate.NewPortMemorizer(
		common.HexToAddress(portDelegateProps.DelegateAddress),
		provider.mintWatcher.delegateClient,
	)

	if err != nil {
		return *new(SwapID), nil, err
	}

	var request *struct {
		HomeAddress    common.Address
		Amount         *big.Int
		ForeignAddress [32]byte
		Status         uint8
	}
	var rqIdInt *big.Int
	for rqIdInt = provider.rqBytesToBigInt(first); rqIdInt != nil; rqIdInt, _ = luState.NextRq(nil, rqIdInt) {
		/**
		 * Due to a fact, that current gateway implementation
		 * on smart contracts (ports) does not have additional
		 * confirmation tx, we should check just for the existence of the swap with that id
		 */
		//if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) == CompletedStatus {
		// if ibRequest := ibState.Request(wavesRequestId); ibRequest != nil && Status(ibRequest.Status) != CompletedStatus {
		// 	continue
		// }
		var requestIdFixed SwapID
		copy(requestIdFixed[:], rqIdInt.Bytes()[0:16]) 
		
		fmt.Printf("ibState: %v \n", ibState)

		ibRequestStatus := ibState.SwapStatusDict[requestIdFixed]
		
		fmt.Printf("status: %v \n", ibRequestStatus)
		if ibRequestStatus != nil && *ibRequestStatus != EthereumRequestStatusSuccess {
			continue
		}

		// validate target address
		luRequest, err := luState.Requests(nil, rqIdInt)
		if err != nil {
			continue
		}

		fmt.Printf("Solana Address: %v \n", base58.Encode(luRequest.ForeignAddress[0:32]))
		if !ValidateSolanaAddress(base58.Encode(luRequest.ForeignAddress[0:32])) {
			continue
		}

		isTokenDataAccountPassed, tokenDataErr := ValidateSolanaTokenAccountOwnershipByTokenProgram(
			provider.solanaClient, 
			base58.Encode(luRequest.ForeignAddress[0:32]),
			provider.config.Meta,
		)
		if tokenDataErr != nil || !isTokenDataAccountPassed {
			continue
		}

		if luRequest.Amount.Uint64() == 0 {
			continue
		}

		portDelegateRequest, err := portDelegate.UnwrapRequests(nil, rqIdInt)
		if err != nil {
			continue
		}

		if !bytes.Equal(portDelegateRequest.ForeignAddress[:], make([]byte, 32)) {
			continue
		}

		fmt.Printf("rqID last: %v \n", rqIdInt.Bytes()[0:16])

		request = &luRequest
		break
	}

	fmt.Printf("rqID on input: %v \n", rqIdInt.Bytes()[0:16])

	if rqIdInt == nil {
		return *new(SwapID), nil, extractors.NotFoundErr
	}

	var swapID SwapID
	var evmSwapID [32]byte
	copy(evmSwapID[:], rqIdInt.Bytes()[0:16])
	copy(swapID[:], rqIdInt.Bytes()[0:16])

	if !provider.mintWatcher.Busy() {
		props := provider.portDelegateProps()
		amountFloat := float64(request.Amount.Int64()) / float64(math.Pow(10, float64(provider.config.DestinationDecimals)))

		go provider.mintWatcher.HandleExtraction(evmSwapID, request, &SolanaMintWatchCommand{
			TheMint: solcommon.PublicKeyFromString(props.TokenMint),
			MintReceiver: solcommon.PublicKeyFromBytes(request.ForeignAddress[:]),
			Amount: amountFloat,
		})
	}

	return swapID, rqIdInt, nil
}

func (provider *EthereumToSolanaExtractionBridge) IBPortState() (*IBPortContractState, error) {
	ibportStateResult, err := provider.solanaClient.GetAccountInfo(provider.solanaCtx, provider.IBPortDataAccount, solclient.GetAccountInfoConfig{
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

	return DecodeIBPortState(stateDecoded), nil
}

func (provider *EthereumToSolanaExtractionBridge) requestSerializer(array []byte) string {
	return base64.StdEncoding.EncodeToString(array)
}

type PortDelegateProps struct { DelegateAddress, NodeUrl, WSEndpoint, TokenMint string }

func (provider *EthereumToSolanaExtractionBridge) portDelegateProps() PortDelegateProps {
	return PortDelegateProps { 
		DelegateAddress: provider.config.Meta["fantom_port_delegate"],
		NodeUrl: provider.config.Meta["fantom_node_url"],
		WSEndpoint: provider.config.Meta["ws_endpoint"],
		TokenMint: provider.config.Meta["token_mint"],
	}
}

func (provider *EthereumToSolanaExtractionBridge) ExtractDirectTransferRequest(ctx context.Context) (*extractors.Data, error) {
	
	// // pick up unprocessed request
	// for swapId, requestStatus := range ibportContractState.SwapStatusDict {
	// 	if requestStatus == EthereumRequestStatusNew {
	// 		currentSwapId = swapId
	// 	}
	// }

	luRequestIds, err := provider.luPortContract.RequestsQueue(nil)
	
	if err != nil {
		return nil, err
	}

	if bytes.Equal(luRequestIds.First[:], make([]byte, 32)) {
		return nil, extractors.NotFoundErr
	}

	if err != nil {
		return nil, err
	}

	rqId, rqIdInt, err := provider.pickRequestFromQueue(provider.luPortContract, luRequestIds.First[:])
	if err != nil {
		return nil, err
	}
	if rqIdInt == nil {
		return nil, extractors.NotFoundErr
	}

	luRequest, err := provider.luPortContract.Requests(nil, rqIdInt)

	decimalsDiff := big.NewInt(provider.config.SourceDecimals - provider.config.DestinationDecimals)

	divideBy := big.NewInt(0).Exp(big.NewInt(10), decimalsDiff, nil)

	targetAmount := luRequest.Amount.Div(luRequest.Amount, divideBy).Uint64()

	solanaDecimals := big.NewInt(0).
		Exp(big.NewInt(10), big.NewInt(provider.config.DestinationDecimals), nil)

	targetAmountCasted := float64(targetAmount) / float64(solanaDecimals.Uint64())

	fmt.Printf("rqId[:]: %v \n", rqId[:])
	fmt.Printf("luRequest.ForeignAddress: %v \n", luRequest.ForeignAddress)
	fmt.Printf("luRequest.ForeignAddress: %v \n", base58.Encode(luRequest.ForeignAddress[:]))
	fmt.Printf("targetAmount: %v \n", targetAmount)
	fmt.Printf("targetAmountCasted: %v \n", targetAmountCasted)

	var resultByteVector [64]byte
	copy(resultByteVector[:], solexecutor.BuildCrossChainMintByteVector(rqId[:], luRequest.ForeignAddress, targetAmountCasted))

	fmt.Printf("result byte string: %v \n", provider.requestSerializer(resultByteVector[:]))
	fmt.Printf("result byte string(len): %v \n", len(resultByteVector))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(resultByteVector[:]),
	}, err
}

func (provider *EthereumToSolanaExtractionBridge) ExtractReverseTransferRequest(ctx context.Context) (*extractors.Data, error) {
	ibState, err := provider.IBPortState()
	if err != nil {
		return nil, err
	}

	var rqSwapID SwapID
	var reverseRequest *IBPortContractUnwrapRequest

	for swapID, burnRequest := range ibState.RequestsDict {
		status := ibState.SwapStatusDict[swapID]

		if status == nil {
			continue
		}

		luRequest, err := provider.luPortContract.Requests(nil, swapID.AsBigInt())
		if err != nil {
			return nil, err
		}

		if luRequest.Status == EthereumRequestStatusSuccess {
			continue
		}

		if luRequest.Status != EthereumRequestStatusNone {
			continue
		}

		if bytes.Equal(luRequest.HomeAddress[:], make([]byte, 20)) || bytes.Equal(luRequest.ForeignAddress[:], make([]byte, 32)) {
			continue
		}

		if !ValidateEthereumBasedAddress(hexutil.Encode(luRequest.ForeignAddress[0:20])) {
			continue
		}

		if burnRequest.Amount == 0 {
			continue
		}

		reverseRequest = burnRequest
		rqSwapID = swapID
		break
	}


	if reverseRequest == nil {
		return nil, extractors.NotFoundErr
	}

	amount := MapAmount(int64(reverseRequest.Amount), provider.config.DestinationDecimals, provider.config.SourceDecimals)

	fmt.Printf("amount mapped: %v \n", amount)

	var rqIDBytes [32]byte
	copy(rqIDBytes[:], rqSwapID[:])

	var amountBytes [32]byte
	amount.FillBytes(amountBytes[:])

	var evmReceiver [20]byte
	copy(evmReceiver[:], reverseRequest.ForeignAddress[0:20])

	result := BuildForEVMByteArray('u', rqIDBytes, amountBytes, evmReceiver)

	fmt.Printf("result byte string: %v \n", provider.requestSerializer(result[:]))
	fmt.Printf("result byte string(len): %v \n", len(result))

	return &extractors.Data{
		Type:  extractors.Base64,
		Value: provider.requestSerializer(result[:]),
	}, err
}
