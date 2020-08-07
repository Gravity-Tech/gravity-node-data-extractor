package model

type EthChainExtractor struct {}

func (e *EthChainExtractor) DataFeedTag() string {
	return "susy-ETH-chain-ext"
}

func (e *EthChainExtractor) Description() string {
	return "susy ETH chain ext"
}

func (e *EthChainExtractor) Data() (interface{}, interface{}) {
	return nil, nil
}

func (e *EthChainExtractor) Info() *ExtractorInfo {
	return &ExtractorInfo{
		Description: e.Description(),
		DataFeedTag: e.DataFeedTag(),
	}
}

func (e *EthChainExtractor) FetchTransferState() *Payment {
	return &Payment{
		ChainName:          "WAVES",
		Amount:             "1000",
		Decimals:           8,
		Timestamp:          10000,
		SourceAddress:      "3b123",
		DestinationAddress: "0x321",
	}
}

func (e *EthChainExtractor) InitTransfer() {
	//...
}

func (e *EthChainExtractor) extractData(params interface{}) []RawData {
	return []RawData {}
}

func (e *EthChainExtractor) mapData(extractedData []RawData) interface{} {
	return nil
}

