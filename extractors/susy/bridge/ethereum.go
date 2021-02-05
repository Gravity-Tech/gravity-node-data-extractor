package bridge

type EthereumExtractionProvider struct {
	//kind        extractors.ExtractorType
	//
	//cache         map[RequestId]time.Time
	//ethClient     *ethclient.Client
	//wavesClient   *client.Client
	//wavesHelper   helpers.ClientHelper
	//luPortAddress string
	//ibPortAddress *ibport.IBPort
	//
	//sourceDecimals      int64
	//destinationDecimals int64
}

const (
	EthereumRequestStatusNone = iota
	EthereumRequestStatusNew
	EthereumRequestStatusRejected
	SuccessEthereum // is 3
	EthereumRequestStatusReturned
)

const (
	MaxRqTimeout = 5 * 60 // 5 min
)